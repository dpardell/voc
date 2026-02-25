package app

import (
	"os"
	"testing"
)

func TestNewApp(t *testing.T) {
	// Mock environment to avoid side effects
	tempDir, _ := os.MkdirTemp("", "voc-app-test-*")
	defer os.RemoveAll(tempDir)
	
	os.Setenv("VOC_USER_DB_PATH", tempDir+"/voc.db")
	defer os.Unsetenv("VOC_USER_DB_PATH")

	app, err := NewApp("en", "fr")
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}
	defer app.Close()

	if app.HostLang != "en" {
		t.Errorf("Expected host lang en, got %s", app.HostLang)
	}
	if app.TargetLang != "fr" {
		t.Errorf("Expected target lang fr, got %s", app.TargetLang)
	}
	if app.DB == nil {
		t.Error("Expected DB to be initialized")
	}
}

func TestGetLLMClientErrors(t *testing.T) {
	app := &App{}

	// Clear env
	os.Unsetenv("VERTEX_API_KEY")
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("VERTEX_PROJECT_ID")

	_, err := app.GetLLMClient()
	if err == nil {
		t.Error("Expected error when API key and project ID are missing")
	}

	os.Setenv("VERTEX_API_KEY", "test-key")
	defer os.Unsetenv("VERTEX_API_KEY")
	
	_, err = app.GetLLMClient()
	if err == nil {
		t.Error("Expected error when project ID is missing")
	}
}
