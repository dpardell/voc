package dictionary

import (
	"testing"
)

func TestGetDefaultKaikkiURL(t *testing.T) {
	tests := []struct {
		lang     string
		expected string
	}{
		{"fr", "https://kaikki.org/dictionary/French/kaikki.org-dictionary-French.jsonl.gz"},
		{"en", "https://kaikki.org/dictionary/English/kaikki.org-dictionary-English.jsonl.gz"},
		{"sk", "https://kaikki.org/dictionary/Slovak/kaikki.org-dictionary-Slovak.jsonl.gz"},
		{"xx", "https://kaikki.org/dictionary/xx/kaikki.org-dictionary-xx.jsonl.gz"},
	}

	for _, tt := range tests {
		got := GetDefaultKaikkiURL(tt.lang)
		if got != tt.expected {
			t.Errorf("GetDefaultKaikkiURL(%s) = %s; want %s", tt.lang, got, tt.expected)
		}
	}

	// Test override
	KaikkiURL = "https://custom.url"
	defer func() { KaikkiURL = "" }()
	
	got := GetDefaultKaikkiURL("fr")
	if got != "https://custom.url" {
		t.Errorf("Expected custom URL, got %s", got)
	}
}
