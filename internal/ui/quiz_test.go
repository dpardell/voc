package ui

import (
	"testing"
	"voc/internal/llm"

	tea "github.com/charmbracelet/bubbletea"
)

func TestQuizModel(t *testing.T) {
	questions := []llm.Question{
		{
			Type:               "multiple_choice",
			Question:           "What is apple?",
			Options:            []string{"fruit", "car", "phone", "building"},
			CorrectAnswerIndex: 0,
			CorrectAnswer:      "fruit",
		},
		{
			Type:          "fill_in_the_blank",
			Question:      "A ___ is a red fruit.",
			CorrectAnswer: "cherry",
		},
	}

	// Mocking LLM client is hard here, so we manually transition to answering state
	m := InitialQuizModel(nil, nil, "")
	m.questions = questions
	m.results = make([]bool, len(questions))
	m.state = stateAnswering
	m.textInput.Focus()

	// Test multiple choice - correct answer
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")})
	qm := m2.(quizModel)
	if qm.score != 1 {
		t.Errorf("Score should be 1, got %d", qm.score)
	}
	if qm.state != stateFeedback {
		t.Errorf("State should be stateFeedback, got %v", qm.state)
	}

	// Move to next question
	m3, _ := qm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	qm = m3.(quizModel)
	if qm.index != 1 {
		t.Errorf("Index should be 1, got %d", qm.index)
	}
	if qm.state != stateAnswering {
		t.Errorf("State should be stateAnswering, got %v", qm.state)
	}

	// Test fill in the blank - correct answer
	qm.textInput.SetValue("cherry")
	m4, _ := qm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	qm = m4.(quizModel)
	if qm.score != 2 {
		t.Errorf("Score should be 2, got %d", qm.score)
	}
	if qm.state != stateFeedback {
		t.Errorf("State should be stateFeedback, got %v", qm.state)
	}

	// Move to next state (updating)
	m5, _ := qm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	qm = m5.(quizModel)
	if qm.state != stateUpdating {
		t.Errorf("State should be stateUpdating, got %v", qm.state)
	}
}
