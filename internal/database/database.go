package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func init() {
	sql.Register("sqlite3_with_collation", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			cl := collate.New(language.English, collate.Loose)
			return conn.RegisterCollation("UNICODE", cl.CompareString)
		},
	})
}

type Database struct {
	db *sql.DB
}

type Word struct {
	ID         int
	Word       string
	Incomplete bool
	CreatedAt  string
	Types      []WordType
}

type WordType struct {
	ID          int
	Type        string
	Definitions []string
}

var DefaultUserDBPath string

func (d *Database) resolveDataPath(filename string) (string, error) {
	if envPath := os.Getenv("VOC_USER_DB_PATH"); envPath != "" {
		dbDir := filepath.Dir(envPath)
		return filepath.Join(dbDir, filename), nil
	}
	if DefaultUserDBPath != "" {
		dbDir := filepath.Dir(DefaultUserDBPath)
		return filepath.Join(dbDir, filename), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "voc", filename), nil
}

func GetProgressPath() (string, error) {
	d := &Database{}
	return d.resolveDataPath("progress.md")
}

func GetDailySayingPath() (string, error) {
	d := &Database{}
	return d.resolveDataPath("daily_saying.json")
}

func GetProgress() (string, error) {
	path, err := GetProgressPath()
	if err != nil {
		return "", err
	}
	
	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(bytes), nil
}

func New() (*Database, error) {
	d := &Database{}
	dbPath, err := d.resolveDataPath("voc.db")
	if err != nil {
		return nil, err
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3_with_collation", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	d.db = db
	return d, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word TEXT UNIQUE NOT NULL,
			incomplete INTEGER DEFAULT 0,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS word_types (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			FOREIGN KEY(word_id) REFERENCES words(id),
			UNIQUE(word_id, type)
		)`,
		`CREATE TABLE IF NOT EXISTS definitions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER NOT NULL,
			word_type_id INTEGER NOT NULL,
			definition TEXT NOT NULL,
			FOREIGN KEY(word_id) REFERENCES words(id),
			FOREIGN KEY(word_type_id) REFERENCES word_types(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) WordExists(word string) (bool, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM words WHERE word = ?", word).Scan(&count)
	return count > 0, err
}

func (d *Database) AddWord(word string, types []WordType, incomplete bool) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO words (word, incomplete, created_at) VALUES (?, ?, ?)",
		word, incomplete, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	wordID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, t := range types {
		res, err := tx.Exec("INSERT INTO word_types (word_id, type) VALUES (?, ?)", wordID, t.Type)
		if err != nil {
			return err
		}
		typeID, err := res.LastInsertId()
		if err != nil {
			return err
		}

		for _, def := range t.Definitions {
			_, err := tx.Exec("INSERT INTO definitions (word_id, word_type_id, definition) VALUES (?, ?, ?)",
				wordID, typeID, def)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (d *Database) GetWord(word string) (*Word, error) {
	var w Word
	err := d.db.QueryRow("SELECT id, word, incomplete, created_at FROM words WHERE word = ?", word).
		Scan(&w.ID, &w.Word, &w.Incomplete, &w.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := d.db.Query("SELECT id, type FROM word_types WHERE word_id = ?", w.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var wt WordType
		if err := rows.Scan(&wt.ID, &wt.Type); err != nil {
			return nil, err
		}

		defRows, err := d.db.Query("SELECT definition FROM definitions WHERE word_type_id = ?", wt.ID)
		if err != nil {
			return nil, err
		}
		defer defRows.Close()

		for defRows.Next() {
			var def string
			if err := defRows.Scan(&def); err != nil {
				return nil, err
			}
			wt.Definitions = append(wt.Definitions, def)
		}
		w.Types = append(w.Types, wt)
	}

	return &w, nil
}

func (d *Database) GetAllWords() ([]Word, error) {
	query := `
		SELECT w.id, w.word, w.incomplete, w.created_at, wt.id, wt.type, d.definition
		FROM words w
		LEFT JOIN word_types wt ON w.id = wt.word_id
		LEFT JOIN definitions d ON wt.id = d.word_type_id
		ORDER BY w.word COLLATE UNICODE
	`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wordMap := make(map[int]*Word)
	var orderedWords []*Word

	for rows.Next() {
		var (
			id         int
			wordText   string
			incomplete bool
			createdAt  string
			wtID       sql.NullInt64
			wtType     sql.NullString
			def        sql.NullString
		)

		if err := rows.Scan(&id, &wordText, &incomplete, &createdAt, &wtID, &wtType, &def); err != nil {
			return nil, err
		}

		w, exists := wordMap[id]
		if !exists {
			w = &Word{
				ID:         id,
				Word:       wordText,
				Incomplete: incomplete,
				CreatedAt:  createdAt,
			}
			wordMap[id] = w
			orderedWords = append(orderedWords, w)
		}

		if wtID.Valid {
			var targetType *WordType
			for i := range w.Types {
				if w.Types[i].ID == int(wtID.Int64) {
					targetType = &w.Types[i]
					break
				}
			}
			if targetType == nil {
				w.Types = append(w.Types, WordType{
					ID:   int(wtID.Int64),
					Type: wtType.String,
				})
				targetType = &w.Types[len(w.Types)-1]
			}
			if def.Valid {
				targetType.Definitions = append(targetType.Definitions, def.String)
			}
		}
	}

	words := make([]Word, 0, len(orderedWords))
	for _, w := range orderedWords {
		words = append(words, *w)
	}

	return words, nil
}

func (d *Database) GetRandomWords(count int) ([]Word, error) {
	// First, select random word IDs
	idRows, err := d.db.Query("SELECT id FROM words ORDER BY RANDOM() LIMIT ?", count)
	if err != nil {
		return nil, err
	}
	defer idRows.Close()

	var ids []interface{}
	for idRows.Next() {
		var id int
		if err := idRows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return []Word{}, nil
	}

	// Build query for details
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	
	query := fmt.Sprintf(`
		SELECT w.id, w.word, w.incomplete, w.created_at, wt.id, wt.type, d.definition
		FROM words w
		LEFT JOIN word_types wt ON w.id = wt.word_id
		LEFT JOIN definitions d ON wt.id = d.word_type_id
		WHERE w.id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := d.db.Query(query, ids...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wordMap := make(map[int]*Word)
	var result []Word

	for rows.Next() {
		var (
			wID        int
			wordText   string
			incomplete bool
			createdAt  string
			wtID       sql.NullInt64
			wtType     sql.NullString
			def        sql.NullString
		)

		if err := rows.Scan(&wID, &wordText, &incomplete, &createdAt, &wtID, &wtType, &def); err != nil {
			return nil, err
		}

		w, exists := wordMap[wID]
		if !exists {
			w = &Word{
				ID:         wID,
				Word:       wordText,
				Incomplete: incomplete,
				CreatedAt:  createdAt,
			}
			wordMap[wID] = w
			// We can't guarantee order with the map, but we'll collect values later
			// For random words, specific order in result doesn't matter much as long as it's the requested set
		}

		if wtID.Valid {
			// Find or create WordType
			// Since WordType is a struct in a slice, we need to track by ID to append definitions
			// A simple way is to re-scan the slice, but a map is O(1)
			
			// We need a unique key for the type within this word context or globally?
			// Locally is cleaner. But let's simplify:
			// Just iterate the existing types to find match. N is tiny (usually < 5 types per word)
			
			var targetType *WordType
			for i := range w.Types {
				if w.Types[i].ID == int(wtID.Int64) {
					targetType = &w.Types[i]
					break
				}
			}
			
			if targetType == nil {
				newType := WordType{
					ID:   int(wtID.Int64),
					Type: wtType.String,
				}
				w.Types = append(w.Types, newType)
				targetType = &w.Types[len(w.Types)-1]
			}
			
			if def.Valid {
				targetType.Definitions = append(targetType.Definitions, def.String)
			}
		}
	}

	for _, w := range wordMap {
		result = append(result, *w)
	}
	
	return result, nil
}

func (d *Database) DeleteWord(word string) error {
	var wordID int
	err := d.db.QueryRow("SELECT id FROM words WHERE word = ?", word).Scan(&wordID)
	if err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM definitions WHERE word_id = ?", wordID)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM word_types WHERE word_id = ?", wordID)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM words WHERE id = ?", wordID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *Database) DeleteAllWords() error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tables := []string{"definitions", "word_types", "words"}
	for _, table := range tables {
		if _, err := tx.Exec("DELETE FROM " + table); err != nil {
			return err
		}
	}
	return tx.Commit()
}
