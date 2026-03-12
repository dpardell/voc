package ui

import (
	"testing"
	"voc/internal/llm"

	tea "github.com/charmbracelet/bubbletea"
)

func TestConvoModel(t *testing.T) {
	m := InitialConvoModel(nil, "Novice level")

	// Test window resize (initializes viewports)
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	cm := m2.(convoModel)
	if cm.width != 80 || cm.height != 24 {
		t.Errorf("Expected size 80x24, got %dx%d", cm.width, cm.height)
	}

	// Test text input change
	cm.textInput.SetValue("Bonjour")
	if cm.textInput.Value() != "Bonjour" {
		t.Error("Text input value not updated")
	}

	// Test message update from LLM (mocked via chatMsg)
	res := &llm.ChatResponse{
		Response: "Bonjour ! Comment vas-tu ?",
		Corrections: []llm.Correction{
			{Incorrect: "Bonjur", Correct: "Bonjour"},
		},
	}
	history := []llm.Message{
		{Role: "user", Content: "Bonjur"},
		{Role: "model", Content: "Bonjour ! Comment vas-tu ?"},
	}

	m3, _ := cm.Update(chatMsg{response: res, history: history})
	cm = m3.(convoModel)

	if len(cm.history) != 2 {
		t.Errorf("Expected history length 2, got %d", len(cm.history))
	}
	if len(cm.lastCorr) != 1 || cm.lastCorr[0].Incorrect != "Bonjur" {
		t.Errorf("Expected corrections to be updated, got %+v", cm.lastCorr)
	}
}
