package dictionary

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func createTestDict(t *testing.T, dbPath string) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to create test dict: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE dictionary (
		word TEXT,
		pos TEXT,
		gloss TEXT
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec(`INSERT INTO dictionary (word, pos, gloss) VALUES 
		('apple', 'noun', 'a round fruit with red or green skin'),
		('apple', 'noun', 'the tree which bears apples'),
		('run', 'verb', 'to move fast on foot'),
		('run', 'noun', 'an act of running')`)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}
}

func TestDictionaryLookup(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voc-dict-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "dictionary_en.db")
	createTestDict(t, dbPath)

	os.Setenv("VOC_DB_PATH", dbPath)
	defer os.Unsetenv("VOC_DB_PATH")

	dict := &Dictionary{}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	dict.db = db
	defer dict.Close()

	t.Run("Lookup word", func(t *testing.T) {
		data, err := dict.Lookup("apple")
		if err != nil {
			t.Errorf("Lookup failed: %v", err)
		}
		if data == nil {
			t.Fatal("Word 'apple' should be found")
		}
		if len(data.Types) != 1 || data.Types[0].Type != "noun" {
			t.Errorf("Expected 1 type 'noun', got %v", data.Types)
		}
		if len(data.Types[0].Definitions) != 2 {
			t.Errorf("Expected 2 definitions, got %d", len(data.Types[0].Definitions))
		}
	})

	t.Run("Lookup missing word", func(t *testing.T) {
		data, err := dict.Lookup("banana")
		if err != nil {
			t.Errorf("Lookup failed: %v", err)
		}
		if data != nil {
			t.Error("Word 'banana' should not be found")
		}
	})
}

func TestDictionarySearch(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voc-dict-search-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "dictionary_en.db")
	createTestDict(t, dbPath)

	os.Setenv("VOC_DB_PATH", dbPath)
	defer os.Unsetenv("VOC_DB_PATH")

	dict := &Dictionary{}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	dict.db = db
	defer dict.Close()

	results, err := dict.Search("ap", 5)
	if err != nil {
		t.Errorf("Search failed: %v", err)
	}
	if len(results) != 1 || results[0] != "apple" {
		t.Errorf("Expected ['apple'], got %v", results)
	}
}
