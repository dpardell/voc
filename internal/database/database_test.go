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

func TestGetProgressPath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voc-test-progress-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "voc.db")
	os.Setenv("VOC_USER_DB_PATH", dbPath)
	defer os.Unsetenv("VOC_USER_DB_PATH")

	path, err := GetProgressPath()
	if err != nil {
		t.Fatalf("Failed to get progress path: %v", err)
	}

	expected := filepath.Join(tempDir, "progress.md")
	if path != expected {
		t.Errorf("Expected progress path %s, got %s", expected, path)
	}
}
