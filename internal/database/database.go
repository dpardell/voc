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
	Language   string
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

func GetProgressPath(lang string) (string, error) {
	d := &Database{}
	return d.resolveDataPath("progress_" + lang + ".md")
}

func GetProgress(lang string) (string, error) {
	path, err := GetProgressPath(lang)
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

	if err := migrate(db); err != nil {
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
			word TEXT NOT NULL,
			language TEXT NOT NULL,
			incomplete INTEGER DEFAULT 0,
			created_at TEXT NOT NULL,
			UNIQUE(word, language)
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

func migrate(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT count(*) FROM pragma_table_info('words') WHERE name='language'").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Add language column. We'll default it to 'fr' for existing data as a best-effort guess.
	_, err = db.Exec("ALTER TABLE words ADD COLUMN language TEXT NOT NULL DEFAULT 'fr'")
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) WordExists(word string, language string) (bool, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM words WHERE word = ? AND language = ?", word, language).Scan(&count)
	return count > 0, err
}

func (d *Database) AddWord(word string, language string, types []WordType, incomplete bool) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO words (word, language, incomplete, created_at) VALUES (?, ?, ?, ?)",
		word, language, incomplete, time.Now().Format(time.RFC3339))
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

func (d *Database) GetWord(word string, language string) (*Word, error) {
	var w Word
	err := d.db.QueryRow("SELECT id, word, language, incomplete, created_at FROM words WHERE word = ? AND language = ?", word, language).
		Scan(&w.ID, &w.Word, &w.Language, &w.Incomplete, &w.CreatedAt)
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

func (d *Database) GetAllWords(language string) ([]Word, error) {
	query := `
		SELECT w.id, w.word, w.language, w.incomplete, w.created_at, wt.id, wt.type, d.definition
		FROM words w
		LEFT JOIN word_types wt ON w.id = wt.word_id
		LEFT JOIN definitions d ON wt.id = d.word_type_id
		WHERE w.language = ?
		ORDER BY w.word COLLATE UNICODE
	`
	rows, err := d.db.Query(query, language)
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
			lang       string
			incomplete bool
			createdAt  string
			wtID       sql.NullInt64
			wtType     sql.NullString
			def        sql.NullString
		)

		if err := rows.Scan(&id, &wordText, &lang, &incomplete, &createdAt, &wtID, &wtType, &def); err != nil {
			return nil, err
		}

		w, exists := wordMap[id]
		if !exists {
			w = &Word{
				ID:         id,
				Word:       wordText,
				Language:   lang,
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

func (d *Database) GetRandomWords(language string, count int) ([]Word, error) {
	idRows, err := d.db.Query("SELECT id FROM words WHERE language = ? ORDER BY RANDOM() LIMIT ?", language, count)
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

	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(`
		SELECT w.id, w.word, w.language, w.incomplete, w.created_at, wt.id, wt.type, d.definition
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

	for rows.Next() {
		var (
			wID        int
			wordText   string
			lang       string
			incomplete bool
			createdAt  string
			wtID       sql.NullInt64
			wtType     sql.NullString
			def        sql.NullString
		)

		if err := rows.Scan(&wID, &wordText, &lang, &incomplete, &createdAt, &wtID, &wtType, &def); err != nil {
			return nil, err
		}

		w, exists := wordMap[wID]
		if !exists {
			w = &Word{
				ID:         wID,
				Word:       wordText,
				Language:   lang,
				Incomplete: incomplete,
				CreatedAt:  createdAt,
			}
			wordMap[wID] = w
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

	result := make([]Word, 0, len(wordMap))
	for _, w := range wordMap {
		result = append(result, *w)
	}

	return result, nil
}

func (d *Database) DeleteWord(word string, language string) error {
	var wordID int
	err := d.db.QueryRow("SELECT id FROM words WHERE word = ? AND language = ?", word, language).Scan(&wordID)
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

func (d *Database) DeleteAllWords(language string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get all word IDs for this language
	rows, err := tx.Query("SELECT id FROM words WHERE language = ?", language)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []interface{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	placeholderStr := strings.Join(placeholders, ",")

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM definitions WHERE word_id IN (%s)", placeholderStr), ids...)
	if err != nil {
		return err
	}
	_, err = tx.Exec(fmt.Sprintf("DELETE FROM word_types WHERE word_id IN (%s)", placeholderStr), ids...)
	if err != nil {
		return err
	}
	_, err = tx.Exec(fmt.Sprintf("DELETE FROM words WHERE id IN (%s)", placeholderStr), ids...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
