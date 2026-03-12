package ui

import (
	"fmt"
	"strings"

	"voc/internal/i18n"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

const (
	minWindowHeight       = 15
	availableHeightOffset = 10
	minAvailableHeight    = 5
	listWidthRatio        = 0.3
	minListWidth          = 20
	layoutPadding         = 4
	minContentWidth       = 10
	searchLimit           = 50
	wrapBuffer            = 2
	inputWidthOffset      = 10
	inputCharLimit        = 156
)

const (
	noMatchesMsg    = i18n.NoMatches
	selectItemMsg   = i18n.SelectItem
	errorPreviewMsg = i18n.ErrorPreview
	resizeWindowMsg = i18n.ResizeWindow
	hintsSearchMsg  = i18n.HintsSearch
	hintsDefMsg     = i18n.HintsDefinition
	savedMsg        = i18n.Saved
	notSavedMsg     = i18n.NotSaved
)

type SearchFunc func(query string, limit int) ([]string, error)
type DefFunc func(word string) (string, error)
type CheckVocabFunc func(word string) (bool, error)
type ToggleVocabFunc func(word string) (bool, error)

type uiState int

const (
	stateSearching uiState = iota
	stateViewingDefinition
)

type model struct {
	textInput       textinput.Model
	viewport        viewport.Model
	searchFunc      SearchFunc
	defFunc         DefFunc
	checkVocabFunc  CheckVocabFunc
	toggleVocabFunc ToggleVocabFunc
	earthSpinner    spinner.Model

	state         uiState
	title         string
	originalTitle string
	results       []string
	cursor        int
	selected      string
	quitted       bool
	inVocab       bool

	previewText           string
	width                 int
	height                int
	listWidth             int
	previewContainerWidth int
	availableHeight       int
}

type searchResultMsg struct {
	results []string
	err     error
}

type previewResultMsg struct {
	text string
	err  error
}

type vocabCheckMsg struct {
	inVocab bool
	err     error
}

type vocabToggleMsg struct {
	inVocab bool
	err     error
}

func InitialModel(title string, initialResults []string, search SearchFunc, def DefFunc, check CheckVocabFunc, toggle ToggleVocabFunc) model {
	InitStyles()

	s := GetSpinner()

	textInput := textinput.New()
	textInput.Placeholder = i18n.T(i18n.TypeToSearch)
	textInput.Focus()
	textInput.CharLimit = inputCharLimit

	viewportModel := viewport.New(0, 0)

	// If we have initial results, populate them
	results := initialResults
	if results == nil {
		results = []string{}
	}

	return model{
		title:           title,
		originalTitle:   title,
		textInput:       textInput,
		viewport:        viewportModel,
		searchFunc:      search,
		defFunc:         def,
		checkVocabFunc:  check,
		toggleVocabFunc: toggle,
		earthSpinner:    s,
		cursor:          0,
		state:           stateSearching,
		results:         results,
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, textinput.Blink, m.earthSpinner.Tick)
	// If we have results, we might want to trigger a preview update for the first item if desirable,
	// but the original logic waits for selection.
	// Actually, in stateSearching, the preview usually shows "Select an item".
	// Let's keep it simple for now.
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		inputWidth := m.width - inputWidthOffset
		if inputWidth < minContentWidth {
			inputWidth = minContentWidth
		}
		m.textInput.Width = inputWidth

		m.availableHeight = m.height - availableHeightOffset
		if m.availableHeight < minAvailableHeight {
			m.availableHeight = minAvailableHeight
		}

		m.listWidth = int(float64(m.width) * listWidthRatio)
		if m.listWidth < minListWidth {
			m.listWidth = minListWidth
		}

		m.previewContainerWidth = m.width - m.listWidth - layoutPadding

		previewContentWidth := m.previewContainerWidth - layoutPadding
		if previewContentWidth < minContentWidth {
			previewContentWidth = minContentWidth
		}

		m.viewport.Width = previewContentWidth
		m.viewport.Height = m.availableHeight

		m.updateViewportContent()

	case tea.KeyMsg:
		if m.state == stateViewingDefinition {
			switch msg.String() {
			case "esc", "backspace":
				m.state = stateSearching
				m.title = m.originalTitle
				m.viewport.Height = m.availableHeight
				m.updateViewportContent()
				return m, m.updatePreview()
			case "ctrl+c":
				m.quitted = true
				return m, tea.Quit
			case "j", "down":
				m.viewport.ScrollDown(1)
			case "k", "up":
				m.viewport.ScrollUp(1)
			case "ctrl+n":
				if m.cursor < len(m.results)-1 {
					m.cursor++
					m.title = strings.ToUpper(m.results[m.cursor])
					return m, tea.Batch(m.updateFullDefinition(), m.checkVocab())
				}
			case "ctrl+p":
				if m.cursor > 0 {
					m.cursor--
					m.title = strings.ToUpper(m.results[m.cursor])
					return m, tea.Batch(m.updateFullDefinition(), m.checkVocab())
				}
			case "ctrl+s":
				return m, m.toggleVocab()
			}
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitted = true
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.results) > 0 {
				m.state = stateViewingDefinition
				m.title = strings.ToUpper(m.results[m.cursor])
				m.viewport.Height = m.height - availableHeightOffset + 2 // Give more room in full view
				return m, tea.Batch(m.updateFullDefinition(), m.checkVocab())
			}
		case tea.KeyUp, tea.KeyCtrlP:
			if m.cursor > 0 {
				m.cursor--
				return m, m.updatePreview()
			}
		case tea.KeyDown, tea.KeyCtrlN:
			if m.cursor < len(m.results)-1 {
				m.cursor++
				return m, m.updatePreview()
			}
		}
	case searchResultMsg:
		if msg.err != nil {
			return m, nil
		}
		m.results = msg.results
		m.cursor = 0
		if len(m.results) == 0 {
			m.previewText = ""
			m.updateViewportContent()
			return m, nil
		}
		return m, m.updatePreview()

	case previewResultMsg:
		if msg.err != nil {
			m.previewText = i18n.T(errorPreviewMsg)
		} else {
			m.previewText = msg.text
		}
		m.updateViewportContent()
		m.viewport.GotoTop()
		m.updateViewportContent()
		m.viewport.GotoTop()
		return m, nil

	case vocabCheckMsg:
		m.inVocab = msg.inVocab
		return m, nil

	case vocabToggleMsg:
		if msg.err == nil {
			m.inVocab = msg.inVocab
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.earthSpinner, cmd = m.earthSpinner.Update(msg)
		return m, cmd
	}

	lastInputValue := m.textInput.Value()
	var inputCmd tea.Cmd
	m.textInput, inputCmd = m.textInput.Update(msg)

	if m.state == stateSearching && m.textInput.Value() != lastInputValue {
		if m.textInput.Value() != "" {
			cmd = tea.Batch(inputCmd, m.performSearch(m.textInput.Value()))
		} else {
			m.results = nil
			m.previewText = ""
			m.updateViewportContent()
			cmd = inputCmd
		}
	} else {
		cmd = inputCmd
	}

	if m.state == stateSearching {
		m.viewport, _ = m.viewport.Update(msg)
	}

	return m, cmd
}

func (m *model) updateViewportContent() {
	if m.viewport.Width > 0 {
		wrapWidth := m.viewport.Width - wrapBuffer
		if wrapWidth < minContentWidth {
			wrapWidth = minContentWidth
		}

		m.viewport.SetContent(wordwrap.String(m.previewText, wrapWidth))
	}
}

func (m model) performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := m.searchFunc(query, searchLimit)
		return searchResultMsg{results: results, err: err}
	}
}

func (m model) updatePreview() tea.Cmd {
	if m.defFunc == nil || len(m.results) == 0 {
		return nil
	}
	word := m.results[m.cursor]
	return func() tea.Msg {
		text, err := m.defFunc(word)
		return previewResultMsg{text: text, err: err}
	}
}

func (m model) updateFullDefinition() tea.Cmd {
	if m.defFunc == nil || len(m.results) == 0 {
		return nil
	}
	word := m.results[m.cursor]
	return func() tea.Msg {
		text, err := m.defFunc(word)
		return previewResultMsg{text: text, err: err}
	}
}

func (m model) checkVocab() tea.Cmd {
	if m.checkVocabFunc == nil || len(m.results) == 0 {
		return nil
	}
	word := m.results[m.cursor]
	return func() tea.Msg {
		inVocab, err := m.checkVocabFunc(word)
		return vocabCheckMsg{inVocab: inVocab, err: err}
	}
}

func (m model) toggleVocab() tea.Cmd {
	if m.toggleVocabFunc == nil || len(m.results) == 0 {
		return nil
	}
	word := m.results[m.cursor]
	return func() tea.Msg {
		inVocab, err := m.toggleVocabFunc(word)
		return vocabToggleMsg{inVocab: inVocab, err: err}
	}
}

func (m model) View() string {
	if m.quitted {
		return ""
	}

	title := titleStyle.Render(m.title)

	prompt := m.earthSpinner.View() + " "
	m.textInput.Prompt = prompt

	searchBoxWidth := m.width - layoutPadding
	searchView := inputBoxStyle.Width(searchBoxWidth).Render(m.textInput.View())

	if m.state == stateViewingDefinition {
		m.viewport.Width = m.width - layoutPadding
		m.updateViewportContent()
		defView := previewStyle.
			Width(m.width - layoutPadding).
			Height(m.height - availableHeightOffset + 2).
			Render(m.viewport.View())

		if m.inVocab {
			title = lipgloss.JoinHorizontal(lipgloss.Left, title, savedStyle.Render(i18n.T(savedMsg)))
		} else {
			title = lipgloss.JoinHorizontal(lipgloss.Left, title, subtleStyle.Render(i18n.T(notSavedMsg)))
		}

		backMsg := hintStyle.Render(i18n.T(hintsDefMsg))
		return fmt.Sprintf("\n%s\n%s\n%s", title, defView, backMsg)
	}

	if m.height < minWindowHeight {
		return fmt.Sprintf("\n%s\n%s\n%s", title, searchView, i18n.T(resizeWindowMsg))
	}

	var listItems []string
	start := 0
	end := len(m.results)

	if end > m.availableHeight {
		start = m.cursor - (m.availableHeight / 2)
		if start < 0 {
			start = 0
		}
		end = start + m.availableHeight
		if end > len(m.results) {
			end = len(m.results)
			start = end - m.availableHeight
			if start < 0 {
				start = 0
			}
		}
	}

	for i := start; i < end; i++ {
		result := m.results[i]
		if i == m.cursor {
			listItems = append(listItems, selectedItemStyle.Render("✨ "+result))
		} else {
			listItems = append(listItems, itemStyle.Render(result))
		}
	}

	if len(m.results) == 0 && m.textInput.Value() != "" {
		listItems = append(listItems, itemStyle.Foreground(subtleColor).Render(i18n.T(noMatchesMsg)))
	}

	listView := lipgloss.NewStyle().
		Width(m.listWidth).
		Render(strings.Join(listItems, "\n"))

	if m.previewText == "" {
		m.viewport.SetContent(i18n.T(selectItemMsg))
	}

	previewView := previewStyle.
		Width(m.previewContainerWidth).
		Height(m.availableHeight).
		Render(m.viewport.View())

	content := lipgloss.JoinHorizontal(lipgloss.Top, listView, previewView)

	hints := hintStyle.Render(i18n.T(hintsSearchMsg))

	return fmt.Sprintf("\n%s\n%s\n%s\n%s", title, searchView, content, hints)
}

func RunFuzzyFinder(title string, initialResults []string, search SearchFunc, def DefFunc, check CheckVocabFunc, toggle ToggleVocabFunc) (string, error) {
	program := tea.NewProgram(InitialModel(title, initialResults, search, def, check, toggle), tea.WithAltScreen())
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}

	if resultModel, ok := finalModel.(model); ok && !resultModel.quitted {
		return resultModel.selected, nil
	}

	return "", nil
}
