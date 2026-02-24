package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSplashModel(t *testing.T) {
	m := InitialSplashModel("Test Saying")

	// Test window resize
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	sm := m2.(splashModel)
	if sm.width != 100 || sm.height != 40 {
		t.Errorf("Expected size 100x40, got %dx%d", sm.width, sm.height)
	}

	// Test navigation
	if sm.cursor != OptionSearch {
		t.Errorf("Expected initial cursor at OptionSearch, got %d", sm.cursor)
	}

	m3, _ := sm.Update(tea.KeyMsg{Type: tea.KeyDown})
	sm = m3.(splashModel)
	if sm.cursor != OptionQuiz {
		t.Errorf("Expected cursor at OptionQuiz after Down key, got %d", sm.cursor)
	}

	m4, _ := sm.Update(tea.KeyMsg{Type: tea.KeyUp})
	sm = m4.(splashModel)
	if sm.cursor != OptionSearch {
		t.Errorf("Expected cursor back at OptionSearch after Up key, got %d", sm.cursor)
	}
}
