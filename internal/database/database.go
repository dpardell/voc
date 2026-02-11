package database

import (
	"database/sql"
	"os"
	"path/filepath"
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

func New() (*Database, error) {
	var dbPath string
	if envPath := os.Getenv("VOC_USER_DB_PATH"); envPath != "" {
		dbPath = envPath
	} else if DefaultUserDBPath != "" {
		dbPath = DefaultUserDBPath
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dbPath = filepath.Join(home, ".local", "share", "voc", "voc.db")
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

	return &Database{db: db}, nil
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
	rows, err := d.db.Query("SELECT id, word, incomplete, created_at FROM words ORDER BY word COLLATE UNICODE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Word, &w.Incomplete, &w.CreatedAt); err != nil {
			return nil, err
		}

		typeRows, err := d.db.Query("SELECT type FROM word_types WHERE word_id = ?", w.ID)
		if err != nil {
			return nil, err
		}

		for typeRows.Next() {
			var t string
			if err := typeRows.Scan(&t); err != nil {
				typeRows.Close()
				return nil, err
			}
			w.Types = append(w.Types, WordType{Type: t})
		}
		typeRows.Close()

		words = append(words, w)
	}
	return words, nil
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
