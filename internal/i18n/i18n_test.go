package i18n

import (
	"testing"
)

func TestTranslations(t *testing.T) {
	// Verify that all keys in 'en' are also in 'fr'
	for key := range en {
		if _, ok := fr[key]; !ok {
			t.Errorf("Key %s missing in French translations", key)
		}
	}

	// Verify that all keys in 'fr' are also in 'en'
	for key := range fr {
		if _, ok := en[key]; !ok {
			t.Errorf("Key %s missing in English translations", key)
		}
	}
}

func TestSetLanguage(t *testing.T) {
	SetHostLanguage("fr")
	if GetHostLanguage() != "fr" {
		t.Errorf("Expected language fr, got %s", GetHostLanguage())
	}

	SetHostLanguage("invalid")
	if GetHostLanguage() != "fr" {
		t.Errorf("Should not switch to invalid language, got %s", GetHostLanguage())
	}

	SetHostLanguage("en")
	if GetHostLanguage() != "en" {
		t.Errorf("Expected language en, got %s", GetHostLanguage())
	}
}

func TestTargetLanguage(t *testing.T) {
	SetTargetLanguage("es")
	if GetTargetLanguage() != "es" {
		t.Errorf("Expected target language es, got %s", GetTargetLanguage())
	}
}

func TestGetLanguageName(t *testing.T) {
	tests := []struct {
		lang     string
		expected string
	}{
		{"en", "English"},
		{"fr", "Français"},
		{"sk", "Slovenčina"},
		{"xx", "xx"},
	}

	for _, tt := range tests {
		got := GetLanguageName(tt.lang)
		if got != tt.expected {
			t.Errorf("GetLanguageName(%s) = %s; want %s", tt.lang, got, tt.expected)
		}
	}
}

func TestT(t *testing.T) {
	SetHostLanguage("en")
	got := T(QuizTitle)
	expected := "VOCABULARY QUIZ"
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}

	// Test with arguments
	got = T(QuizQuestion, 1, 10)
	expected = "Question 1/10"
	if got != expected {
		t.Errorf("Expected %s, got %s", expected, got)
	}

	// Test fallback for missing key
	missingKey := StringID("MissingKey")
	got = T(missingKey)
	if got != string(missingKey) {
		t.Errorf("Expected fallback to key name %s, got %s", string(missingKey), got)
	}
}
