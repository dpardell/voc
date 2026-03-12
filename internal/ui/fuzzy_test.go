package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFuzzyInitialModel(t *testing.T) {
	m := InitialModel("Test Title", nil, nil, nil, nil, nil)
	if m.title != "Test Title" {
		t.Errorf("expected title 'Test Title', got %q", m.title)
	}
	if m.state != stateSearching {
		t.Error("expected initial state to be stateSearching")
	}
	if m.results == nil {
		t.Error("expected results to be initialized (not nil)")
	}
}

func TestFuzzyWindowResize(t *testing.T) {
	m := InitialModel("Test", nil, nil, nil, nil, nil)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	result, ok := updated.(model)
	if !ok {
		t.Fatal("Update did not return a model")
	}
	if result.width != 120 {
		t.Errorf("expected width 120, got %d", result.width)
	}
	if result.height != 40 {
		t.Errorf("expected height 40, got %d", result.height)
	}
}

func TestFuzzySearchInput(t *testing.T) {
	m := InitialModel("Test", nil, func(query string, limit int) ([]string, error) {
		return []string{"bonjour"}, nil
	}, nil, nil, nil)

	// Send a window size first so layout is initialized
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = updated.(model)

	// Type a character
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	result, ok := updated.(model)
	if !ok {
		t.Fatal("Update did not return a model")
	}
	if result.textInput.Value() != "b" {
		t.Errorf("expected text input value 'b', got %q", result.textInput.Value())
	}
}
