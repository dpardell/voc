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
	SetLanguage("fr")
	if GetLanguage() != "fr" {
		t.Errorf("Expected language fr, got %s", GetLanguage())
	}

	SetLanguage("invalid")
	if GetLanguage() != "fr" {
		t.Errorf("Should not switch to invalid language, got %s", GetLanguage())
	}

	SetLanguage("en")
	if GetLanguage() != "en" {
		t.Errorf("Expected language en, got %s", GetLanguage())
	}
}

func TestT(t *testing.T) {
	SetLanguage("en")
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
