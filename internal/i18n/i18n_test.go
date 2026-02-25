package i18n

import (
	"testing"
)

func TestTranslations(t *testing.T) {
	langs := []string{"en", "fr", "cs", "sk", "es", "de"}
	
	for _, lang := range langs {
		locale, ok := locales[lang]
		if !ok {
			t.Errorf("Locale %s not found in locales map", lang)
			continue
		}

		// Check against 'en' as the reference
		for key := range en {
			if _, ok := locale[key]; !ok {
				t.Errorf("Key %s missing in %s translations", key, lang)
			}
		}

		// Check for extra keys in locale
		for key := range locale {
			if _, ok := en[key]; !ok {
				t.Errorf("Extra key %s found in %s translations", key, lang)
			}
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
