package ui

import (
	"context"
	"fmt"
	"strings"

	"voc/internal/i18n"
	"voc/internal/llm"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type quizState int

const (
	stateGenerating quizState = iota
	stateAnswering
	stateFeedback
	stateFinished
	stateUpdating
)

type quizModel struct {
	client      llm.LLMClient
	targetWords []string
	progress    string
	ctx         context.Context

	questions []llm.Question
	results   []bool
	index     int
	score     int
	state     quizState
	selected  int
	quitted   bool
	width     int
	height    int
	textInput textinput.Model
	spinner   spinner.Model
	err       error

	// Results for the caller
	finalProgress string
	hostLang      string
	targetLang    string
}

func InitialQuizModel(client llm.LLMClient, targetWords []string, progress string) quizModel {
	InitStyles()
	ti := textinput.New()
	ti.Placeholder = i18n.T(i18n.QuizPlaceholder)
	ti.CharLimit = 156
	ti.Width = 30

	s := GetSpinner()

	return quizModel{
		client:      client,
		targetWords: targetWords,
		progress:    progress,
		ctx:         context.Background(),
		state:       stateGenerating,
		selected:    -1,
		textInput:   ti,
		spinner:     s,
		hostLang:    i18n.GetLanguageName(i18n.GetHostLanguage()),
		targetLang:  i18n.GetLanguageName(i18n.GetTargetLanguage()),
	}
}


type questionsMsg []llm.Question
type updateMsg string
type errMsg error

func (m quizModel) generateQuestions() tea.Cmd {
	return func() tea.Msg {
		questions, err := m.client.GenerateQuiz(m.ctx, m.targetWords, m.progress, m.hostLang, m.targetLang)
		if err != nil {
			return errMsg(err)
		}
		return questionsMsg(questions)
	}
}

func (m quizModel) updateProgress(sessionResults string, contextWords []string) tea.Cmd {
	return func() tea.Msg {
		newProgress, err := m.client.UpdateProgress(m.ctx, m.progress, sessionResults, contextWords)
		if err != nil {
			return errMsg(err)
		}
		return updateMsg(newProgress)
	}
}

func (m quizModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.generateQuestions())
}

func (m quizModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		spCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinner.TickMsg:
		m.spinner, spCmd = m.spinner.Update(msg)
		return m, spCmd

	case questionsMsg:
		m.questions = msg
		m.results = make([]bool, len(m.questions))
		m.state = stateAnswering
		if m.questions[0].Type == "fill_in_the_blank" {
			m.textInput.Focus()
		}
		return m, nil

	case updateMsg:
		m.finalProgress = string(msg)
		m.state = stateFinished
		return m, nil

	case errMsg:
		m.err = error(msg)
		return m, nil

	case tea.KeyMsg:
		if m.err != nil {
			return m, tea.Quit
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitted = true
			return m, tea.Quit

		case "q":
			if m.state == stateAnswering && m.questions[m.index].Type != "fill_in_the_blank" {
				m.quitted = true
				return m, tea.Quit
			}

		case "1", "2", "3", "4":
			if m.state == stateAnswering && m.questions[m.index].Type == "multiple_choice" {
				m.selected = int(msg.String()[0]-'1')
				if m.selected == m.questions[m.index].CorrectAnswerIndex {
					m.score++
					m.results[m.index] = true
				} else {
					m.results[m.index] = false
				}
				m.state = stateFeedback
				return m, nil
			}

		case "enter":
			if m.state == stateAnswering && m.questions[m.index].Type == "fill_in_the_blank" {
				answer := strings.TrimSpace(strings.ToLower(m.textInput.Value()))
				correct := strings.TrimSpace(strings.ToLower(m.questions[m.index].CorrectAnswer))
				if answer == correct {
					m.score++
					m.results[m.index] = true
				} else {
					m.results[m.index] = false
				}
				m.state = stateFeedback
				m.textInput.Blur()
				return m, nil
			}

			if m.state == stateFeedback {
				if m.index < len(m.questions)-1 {
					m.index++
					m.selected = -1
					m.state = stateAnswering
					m.textInput.Reset()
					if m.questions[m.index].Type == "fill_in_the_blank" {
						m.textInput.Focus()
					}
				} else {
					// Format results and start update
					var sessionResults strings.Builder
					sessionResults.WriteString(fmt.Sprintf("Score: %d/%d\n", m.score, len(m.questions)))
					for i, q := range m.questions {
						status := "Correct"
						if !m.results[i] {
							status = "Incorrect"
						}
						sessionResults.WriteString(fmt.Sprintf("- Word: %s, Question: %s, Result: %s\n", q.TargetWord, q.Question, status))
					}
					
					m.state = stateUpdating
					return m, m.updateProgress(sessionResults.String(), m.targetWords)
				}
				return m, nil
			} else if m.state == stateFinished {
				m.quitted = false
				return m, tea.Quit
			}
		}
	}

	if m.state == stateAnswering && m.questions[m.index].Type == "fill_in_the_blank" {
		m.textInput, tiCmd = m.textInput.Update(msg)
	}

	return m, tea.Batch(tiCmd, spCmd)
}

func (m quizModel) View() string {
	if m.quitted && m.state != stateFinished {
		return ""
	}

	title := titleStyle.Render(i18n.T(i18n.QuizTitle))

	if m.err != nil {
		return fmt.Sprintf("\n%s\n\n%s\n\n%s", 
			title, 
			lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("Error: "+m.err.Error()),
			hintStyle.Render("Press any key to exit"))
	}

	if m.state == stateGenerating {
		return fmt.Sprintf("\n%s\n\n%s %s", title, m.spinner.View(), i18n.T(i18n.QuizGenerating))
	}

	if m.state == stateUpdating {
		return fmt.Sprintf("\n%s\n\n%s %s", title, m.spinner.View(), i18n.T(i18n.QuizUpdating))
	}
	
	if m.state == stateFinished {
		scoreText := fmt.Sprintf("\n%s\n", i18n.T(i18n.QuizCompleted, m.score, len(m.questions)))
		return fmt.Sprintf("\n%s\n%s\n%s", title, scoreText, hintStyle.Render("Press Enter to exit"))
	}

	q := m.questions[m.index]
	header := subtleStyle.Render(i18n.T(i18n.QuizQuestion, m.index+1, len(m.questions)))
	score := subtleStyle.Render(i18n.T(i18n.QuizScore, m.score, len(m.questions)))
	
	questionText := lipgloss.NewStyle().
		Bold(true).
		Margin(1, 0, 1, 2).
		Width(m.width - 4).
		Render(q.Question)

	var optionsView string
	if q.Type == "multiple_choice" {
		var options []string
		for i, opt := range q.Options {
			prefix := fmt.Sprintf("%d) ", i+1)
			style := itemStyle
			
			if m.state == stateFeedback {
				if i == q.CorrectAnswerIndex {
					style = correctStyle
					prefix = "✅ "
				} else if i == m.selected {
					style = wrongStyle
					prefix = "❌ "
				} else {
					style = subtleStyle
				}
			} else if i == m.selected {
				style = selectedItemStyle
			}

			options = append(options, style.Render(prefix+opt))
		}
		optionsView = strings.Join(options, "\n")
	} else {
		// Fill in the blank
		optionsView = m.textInput.View()
	}

	feedback := ""
	if m.state == stateFeedback {
		if m.results[m.index] {
			feedback = "\n" + correctStyle.Render(i18n.T(i18n.QuizCorrect))
		} else {
			feedback = "\n" + wrongStyle.Render(i18n.T(i18n.QuizWrong, q.CorrectAnswer))
		}
	}

	hints := hintStyle.Render(i18n.T(i18n.HintsQuiz))
	if q.Type == "fill_in_the_blank" && m.state == stateAnswering {
		hints = hintStyle.Render(i18n.T(i18n.QuizHintType))
	}

	return fmt.Sprintf("\n%s %s\n%s\n%s\n\n%s\n%s\n\n%s", title, score, header, questionText, optionsView, feedback, hints)
}

type QuizResult struct {
	Score       int
	Questions   []llm.Question
	Answers     []bool // true if correct
	NewProgress string
}

func RunQuiz(client llm.LLMClient, targetWords []string, progress string) (*QuizResult, error) {
	m := InitialQuizModel(client, targetWords, progress)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	resModel := finalModel.(quizModel)
	
	return &QuizResult{
		Score:       resModel.score,
		Questions:   resModel.questions,
		Answers:     resModel.results,
		NewProgress: resModel.finalProgress,
	}, nil
}


