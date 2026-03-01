package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadDefaults(t *testing.T) {
	// Stub deps to return nothing/error
	origGetenv := getenv
	origUserConfigDir := userConfigDir
	origReadFile := readFile
	defer func() {
		getenv = origGetenv
		userConfigDir = origUserConfigDir
		readFile = origReadFile
	}()

	getenv = func(key string) string { return "" }
	userConfigDir = func() (string, error) { return "", os.ErrNotExist }
	readFile = func(name string) ([]byte, error) { return nil, os.ErrNotExist }

	s, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if s.HostLang != "en" {
		t.Errorf("Expected default host lang en, got %s", s.HostLang)
	}
	if s.TargetLang != "fr" {
		t.Errorf("Expected default target lang fr, got %s", s.TargetLang)
	}
}

func TestConfigLoadEnv(t *testing.T) {
	origGetenv := getenv
	origUserConfigDir := userConfigDir
	origReadFile := readFile
	defer func() {
		getenv = origGetenv
		userConfigDir = origUserConfigDir
		readFile = origReadFile
	}()

	getenv = func(key string) string {
		switch key {
		case "VOC_HOST_LANG":
			return "fr"
		case "VOC_TARGET_LANG":
			return "es"
		}
		return ""
	}
	// No config file
	userConfigDir = func() (string, error) { return "", os.ErrNotExist }
	readFile = func(name string) ([]byte, error) { return nil, os.ErrNotExist }

	s, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if s.HostLang != "fr" {
		t.Errorf("Expected host lang fr from env, got %s", s.HostLang)
	}
	if s.TargetLang != "es" {
		t.Errorf("Expected target lang es from env, got %s", s.TargetLang)
	}
}

func TestConfigLoadFile(t *testing.T) {
	origGetenv := getenv
	origUserConfigDir := userConfigDir
	origReadFile := readFile
	defer func() {
		getenv = origGetenv
		userConfigDir = origUserConfigDir
		readFile = origReadFile
	}()

	getenv = func(key string) string { return "" }
	userConfigDir = func() (string, error) { return "/mock/config", nil }
	readFile = func(name string) ([]byte, error) {
		if filepath.Base(name) == "settings.yaml" {
			return []byte("host_lang: de\ntarget_lang: it"), nil
		}
		return nil, os.ErrNotExist
	}

	s, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if s.HostLang != "de" {
		t.Errorf("Expected host lang de from file, got %s", s.HostLang)
	}
	if s.TargetLang != "it" {
		t.Errorf("Expected target lang it from file, got %s", s.TargetLang)
	}
}

func TestConfigSaveLoad(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "voc-config-save-test-*")
	defer os.RemoveAll(tempDir)
	
	// os.UserConfigDir uses XDG_CONFIG_HOME on Linux
	origXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	defer os.Setenv("XDG_CONFIG_HOME", origXDG)

	s := &Settings{
		HostLang:   "de",
		TargetLang: "it",
		Dictionaries: map[string]string{
			"it": "https://example.com/it.jsonl",
		},
	}

	if err := s.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	configPath := filepath.Join(tempDir, "voc", "settings.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}

	// Load back
	s2, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if s2.HostLang != "de" {
		t.Errorf("Expected host lang de, got %s", s2.HostLang)
	}
	if s2.TargetLang != "it" {
		t.Errorf("Expected target lang it, got %s", s2.TargetLang)
	}
	if s2.Dictionaries["it"] != "https://example.com/it.jsonl" {
		t.Errorf("Expected dictionary URL, got %s", s2.Dictionaries["it"])
	}
}
