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
			Question:           "Elle ___ (aller) au marché chaque semaine.",
			Options:            []string{"va", "allait", "ira", "est allée"},
			CorrectAnswerIndex: 0,
			CorrectAnswer:      "va",
			TargetWord:         "aller",
		},
		{
			Type:               "multiple_choice",
			Question:           "Which option correctly completes: \"Si j'avais le temps, je ___ voyager.\"",
			Options:            []string{"voudrais", "voulais", "veux", "voudrai"},
			CorrectAnswerIndex: 0,
			CorrectAnswer:      "voudrais",
			TargetWord:         "vouloir",
		},
	}

	// Manually transition to answering state (bypasses LLM)
	m := InitialQuizModel(nil, nil, "")
	m.questions = questions
	m.results = make([]bool, len(questions))
	m.state = stateAnswering

	// Test correct answer on first question
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

	// Test wrong answer on second question
	m4, _ := qm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	qm = m4.(quizModel)
	if qm.score != 1 {
		t.Errorf("Score should still be 1, got %d", qm.score)
	}
	if qm.state != stateFeedback {
		t.Errorf("State should be stateFeedback, got %v", qm.state)
	}

	// Move to next state (updating) after last question
	m5, _ := qm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	qm = m5.(quizModel)
	if qm.state != stateUpdating {
		t.Errorf("State should be stateUpdating, got %v", qm.state)
	}
}
