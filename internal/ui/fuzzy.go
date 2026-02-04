package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle.Copy()
	matchStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true) // Yellow bold for matches

	focusedButton = focusedStyle.Copy().Render
	blurredButton = fmt.Sprintf
)

type SearchFunc func(query string, limit int) ([]string, error)
type PreviewFunc func(word string) (string, error)

type model struct {
	textInput   textinput.Model
	searchFunc  SearchFunc
	previewFunc PreviewFunc

	results  []string
	cursor   int
	selected string // The final choice
	quitted  bool

	previewText string
	width       int
	height      int
	err         error
}

type searchResultMsg struct {
	results []string
	err     error
}

type previewResultMsg struct {
	text string
	err  error
}

func InitialModel(prompt string, search SearchFunc, preview PreviewFunc) model {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	ti.Prompt = prompt
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput:   ti,
		searchFunc:  search,
		previewFunc: preview,
		cursor:      0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitted = true
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.results) > 0 {
				m.selected = m.results[m.cursor]
				return m, tea.Quit
			}
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
				return m, m.updatePreview()
			}
		case tea.KeyDown:
			if m.cursor < len(m.results)-1 {
				m.cursor++
				return m, m.updatePreview()
			}
		}

	case searchResultMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.results = msg.results
		m.cursor = 0 // Reset cursor on new search
		return m, m.updatePreview()

	case previewResultMsg:
		if msg.err != nil {
			m.previewText = "Error loading preview"
			return m, nil
		}
		m.previewText = msg.text
	}

	var tiCmd tea.Cmd
	m.textInput, tiCmd = m.textInput.Update(msg)

	// If text changed, trigger search
	if m.textInput.Value() != "" {
		cmd = tea.Batch(tiCmd, m.performSearch(m.textInput.Value()))
	} else {
		m.results = nil
		m.previewText = ""
		cmd = tiCmd
	}

	return m, cmd
}

func (m model) performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		// Use a slight delay to debounce? (Maybe overkill for local DB)
		results, err := m.searchFunc(query, 50)
		return searchResultMsg{results: results, err: err}
	}
}

func (m model) updatePreview() tea.Cmd {
	if m.previewFunc == nil || len(m.results) == 0 {
		return nil
	}
	word := m.results[m.cursor]
	return func() tea.Msg {
		text, err := m.previewFunc(word)
		return previewResultMsg{text: text, err: err}
	}
}

func (m model) View() string {
	if m.quitted {
		return ""
	}

	s := fmt.Sprintf("\n%s\n\n", m.textInput.View())

	// Split view: Left results, Right preview

	// Prevent panic if window too small
	if m.height < 10 {
		return s + "(Resize window to view results)"
	}

	maxListHeight := m.height - 5 // Account for prompt and footer

	// Create list view
	var listItems []string
	start := 0
	end := len(m.results)
	if end > maxListHeight {
		start = m.cursor - (maxListHeight / 2)
		if start < 0 {
			start = 0
		}
		end = start + maxListHeight
		if end > len(m.results) {
			end = len(m.results)
			start = end - maxListHeight
			if start < 0 {
				start = 0
			}
		}
	}

	for i := start; i < end; i++ {
		res := m.results[i]
		if i == m.cursor {
			listItems = append(listItems, focusedStyle.Render("> "+res))
		} else {
			listItems = append(listItems, fmt.Sprintf("  %s", res))
		}
	}

	listView := strings.Join(listItems, "\n")

	// Layout
	// We want roughly 50% width for list, 50% for preview
	halfWidth := (m.width / 2) - 2

	listView = lipgloss.NewStyle().Width(halfWidth).Render(listView)
	previewView := lipgloss.NewStyle().
		Width(halfWidth).
		Border(lipgloss.NormalBorder(), false, false, false, true). // Left border
		PaddingLeft(1).
		Foreground(lipgloss.Color("245")).
		Render(m.previewText)

	// Combine horizontally
	content := lipgloss.JoinHorizontal(lipgloss.Top, listView, previewView)

	return s + content
}

func RunFuzzyFinder(prompt string, search SearchFunc, preview PreviewFunc) (string, error) {
	p := tea.NewProgram(InitialModel(prompt, search, preview), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	if finalModel, ok := m.(model); ok && !finalModel.quitted {
		return finalModel.selected, nil
	}

	return "", nil
}
