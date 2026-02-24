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
	types := []WordType{
		{
			Type:        "noun",
			Definitions: []string{"a procedure intended to establish the quality, performance, or reliability of something"},
		},
	}

	// Test Add
	err = db.AddWord(word, types, false)
	if err != nil {
		t.Errorf("Failed to add word: %v", err)
	}

	// Test Exists
	exists, err := db.WordExists(word)
	if err != nil || !exists {
		t.Errorf("Word should exist: %v", err)
	}

	// Test Get
	got, err := db.GetWord(word)
	if err != nil || got == nil {
		t.Errorf("Failed to get word: %v", err)
	}
	if got.Word != word {
		t.Errorf("Expected word %s, got %s", word, got.Word)
	}
	if len(got.Types) != 1 || got.Types[0].Type != "noun" {
		t.Errorf("Word types not retrieved correctly")
	}

	// Test Delete
	err = db.DeleteWord(word)
	if err != nil {
		t.Errorf("Failed to delete word: %v", err)
	}

	exists, _ = db.WordExists(word)
	if exists {
		t.Error("Word should no longer exist after deletion")
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
