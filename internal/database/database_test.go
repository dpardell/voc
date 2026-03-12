package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "voc.db")
	os.Setenv("VOC_USER_DB_PATH", dbPath)
	defer os.Unsetenv("VOC_USER_DB_PATH")

	db, err := New()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Database file was not created at %s", dbPath)
	}
}

func TestAddGetDeleteWord(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voc-test-word-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "voc.db")
	os.Setenv("VOC_USER_DB_PATH", dbPath)
	defer os.Unsetenv("VOC_USER_DB_PATH")

	db, err := New()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	word := "testword"
	lang := "en"
	types := []WordType{
		{
			Type:        "noun",
			Definitions: []string{"a procedure intended to establish the quality, performance, or reliability of something"},
		},
	}

	// Test Add
	err = db.AddWord(word, lang, types, false)
	if err != nil {
		t.Errorf("Failed to add word: %v", err)
	}

	// Test Exists
	exists, err := db.WordExists(word, lang)
	if err != nil || !exists {
		t.Errorf("Word should exist: %v", err)
	}

	// Test Get
	got, err := db.GetWord(word, lang)
	if err != nil || got == nil {
		t.Errorf("Failed to get word: %v", err)
	}
	if got.Word != word {
		t.Errorf("Expected word %s, got %s", word, got.Word)
	}
	if got.Language != lang {
		t.Errorf("Expected language %s, got %s", lang, got.Language)
	}
	if len(got.Types) != 1 || got.Types[0].Type != "noun" {
		t.Errorf("Word types not retrieved correctly")
	}

	// Test Delete
	err = db.DeleteWord(word, lang)
	if err != nil {
		t.Errorf("Failed to delete word: %v", err)
	}

	exists, _ = db.WordExists(word, lang)
	if exists {
		t.Error("Word should no longer exist after deletion")
	}
}

func TestMultiLanguage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voc-test-multi-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "voc.db")
	os.Setenv("VOC_USER_DB_PATH", dbPath)
	defer os.Unsetenv("VOC_USER_DB_PATH")

	db, err := New()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	word := "chat"
	// 'chat' in English (informal talk) vs 'chat' in French (cat)
	typesEN := []WordType{{Type: "noun", Definitions: []string{"informal talk"}}}
	typesFR := []WordType{{Type: "noun", Definitions: []string{"cat"}}}

	_ = db.AddWord(word, "en", typesEN, false)
	_ = db.AddWord(word, "fr", typesFR, false)

	gotEN, _ := db.GetWord(word, "en")
	if gotEN == nil || gotEN.Types[0].Definitions[0] != "informal talk" {
		t.Errorf("Failed to get English 'chat'")
	}

	gotFR, _ := db.GetWord(word, "fr")
	if gotFR == nil || gotFR.Types[0].Definitions[0] != "cat" {
		t.Errorf("Failed to get French 'chat'")
	}

	wordsEN, _ := db.GetAllWords("en")
	if len(wordsEN) != 1 || wordsEN[0].Word != "chat" {
		t.Errorf("GetAllWords('en') should return 1 word")
	}

	wordsFR, _ := db.GetAllWords("fr")
	if len(wordsFR) != 1 || wordsFR[0].Word != "chat" {
		t.Errorf("GetAllWords('fr') should return 1 word")
	}
}

func newTestDB(t *testing.T) (*Database, string) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "voc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	dbPath := filepath.Join(tempDir, "voc.db")
	os.Setenv("VOC_USER_DB_PATH", dbPath)
	db, err := New()
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
		os.Unsetenv("VOC_USER_DB_PATH")
		os.RemoveAll(tempDir)
	})
	return db, tempDir
}

func addWord(t *testing.T, db *Database, word, lang string) {
	t.Helper()
	types := []WordType{{Type: "noun", Definitions: []string{"definition of " + word}}}
	if err := db.AddWord(word, lang, types, false); err != nil {
		t.Fatalf("AddWord(%q, %q) failed: %v", word, lang, err)
	}
}

func TestGetAllWords(t *testing.T) {
	db, _ := newTestDB(t)

	addWord(t, db, "chat", "en")
	addWord(t, db, "apple", "en")
	addWord(t, db, "zebra", "en")
	addWord(t, db, "chien", "fr")

	enWords, err := db.GetAllWords("en")
	if err != nil {
		t.Fatalf("GetAllWords(en) failed: %v", err)
	}
	if len(enWords) != 3 {
		t.Fatalf("expected 3 en words, got %d", len(enWords))
	}
	// Verify alphabetical order
	if enWords[0].Word != "apple" || enWords[1].Word != "chat" || enWords[2].Word != "zebra" {
		t.Errorf("unexpected order: %v", []string{enWords[0].Word, enWords[1].Word, enWords[2].Word})
	}

	frWords, err := db.GetAllWords("fr")
	if err != nil {
		t.Fatalf("GetAllWords(fr) failed: %v", err)
	}
	if len(frWords) != 1 {
		t.Errorf("expected 1 fr word, got %d", len(frWords))
	}
}

func TestGetRandomWords(t *testing.T) {
	db, _ := newTestDB(t)

	for _, w := range []string{"un", "deux", "trois", "quatre", "cinq"} {
		addWord(t, db, w, "en")
	}

	got3, err := db.GetRandomWords("en", 3)
	if err != nil {
		t.Fatalf("GetRandomWords(en, 3) failed: %v", err)
	}
	if len(got3) != 3 {
		t.Errorf("expected 3 words, got %d", len(got3))
	}

	// Request more than available — should be capped at actual count
	gotAll, err := db.GetRandomWords("en", 100)
	if err != nil {
		t.Fatalf("GetRandomWords(en, 100) failed: %v", err)
	}
	if len(gotAll) != 5 {
		t.Errorf("expected 5 words (capped), got %d", len(gotAll))
	}

	// Request zero — should return empty
	gotZero, err := db.GetRandomWords("en", 0)
	if err != nil {
		t.Fatalf("GetRandomWords(en, 0) failed: %v", err)
	}
	if len(gotZero) != 0 {
		t.Errorf("expected 0 words, got %d", len(gotZero))
	}
}

func TestDeleteAllWords(t *testing.T) {
	db, _ := newTestDB(t)

	addWord(t, db, "apple", "en")
	addWord(t, db, "banana", "en")
	addWord(t, db, "chien", "fr")

	if err := db.DeleteAllWords("en"); err != nil {
		t.Fatalf("DeleteAllWords(en) failed: %v", err)
	}

	enWords, err := db.GetAllWords("en")
	if err != nil {
		t.Fatalf("GetAllWords(en) after delete failed: %v", err)
	}
	if len(enWords) != 0 {
		t.Errorf("expected 0 en words after DeleteAllWords, got %d", len(enWords))
	}

	frWords, err := db.GetAllWords("fr")
	if err != nil {
		t.Fatalf("GetAllWords(fr) failed: %v", err)
	}
	if len(frWords) != 1 {
		t.Errorf("expected fr words untouched (1), got %d", len(frWords))
	}
}

func TestGetProgress(t *testing.T) {
	db, tempDir := newTestDB(t)
	_ = db

	// No file yet — should return empty string, nil error
	content, err := GetProgress("fr")
	if err != nil {
		t.Fatalf("GetProgress with no file returned error: %v", err)
	}
	if content != "" {
		t.Errorf("expected empty content, got %q", content)
	}

	// Write a progress file manually
	progressPath := filepath.Join(tempDir, "progress_fr.md")
	expected := "## Summary\nTest progress."
	if writeErr := os.WriteFile(progressPath, []byte(expected), 0644); writeErr != nil {
		t.Fatalf("failed to write progress file: %v", writeErr)
	}

	content, err = GetProgress("fr")
	if err != nil {
		t.Fatalf("GetProgress returned error: %v", err)
	}
	if content != expected {
		t.Errorf("expected %q, got %q", expected, content)
	}
}

func TestGetProgressPath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voc-test-progress-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "voc.db")
	os.Setenv("VOC_USER_DB_PATH", dbPath)
	defer os.Unsetenv("VOC_USER_DB_PATH")

	path, err := GetProgressPath("fr")
	if err != nil {
		t.Fatalf("Failed to get progress path: %v", err)
	}

	expected := filepath.Join(tempDir, "progress_fr.md")
	if path != expected {
		t.Errorf("Expected progress path %s, got %s", expected, path)
	}
}
