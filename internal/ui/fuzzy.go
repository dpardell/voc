package ui

import (
	"fmt"
	"strings"

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
	noMatchesMsg    = "No matches found..."
	selectItemMsg   = "Select an item to see details..."
	errorPreviewMsg = "Error loading preview"
	resizeWindowMsg = "(Resize window to view results)"
)

var (
	subtleColor    = lipgloss.Color("241")
	highlightColor = lipgloss.Color("212")
	titleBgColor   = lipgloss.Color("62")

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(titleBgColor).
			Padding(0, 1).
			MarginLeft(1).
			Bold(true)

	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlightColor).
			Padding(0, 1).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(highlightColor).
				Bold(true)

	previewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor).
			Padding(0, 1).
			MarginLeft(1)
)

type SearchFunc func(query string, limit int) ([]string, error)
type PreviewFunc func(word string) (string, error)

type model struct {
	textInput   textinput.Model
	viewport    viewport.Model
	searchFunc  SearchFunc
	previewFunc PreviewFunc

	title    string
	results  []string
	cursor   int
	selected string
	quitted  bool

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

func InitialModel(title string, search SearchFunc, preview PreviewFunc) model {
	textInput := textinput.New()
	textInput.Placeholder = "Type to search..."
	textInput.Focus()
	textInput.Prompt = "> "
	textInput.CharLimit = inputCharLimit

	viewportModel := viewport.New(0, 0)

	return model{
		title:       title,
		textInput:   textInput,
		viewport:    viewportModel,
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
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitted = true
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.results) > 0 {
				m.selected = m.results[m.cursor]
				return m, tea.Quit
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
			m.previewText = errorPreviewMsg
		} else {
			m.previewText = msg.text
		}
		m.updateViewportContent()
		m.viewport.GotoTop()
		return m, nil
	}

	lastInputValue := m.textInput.Value()
	var inputCmd tea.Cmd
	m.textInput, inputCmd = m.textInput.Update(msg)

	if m.textInput.Value() != lastInputValue {
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

	m.viewport, _ = m.viewport.Update(msg)

	return m, cmd
}

func (m *model) updateViewportContent() {
	if m.viewport.Width > 0 {
		wrapWidth := m.viewport.Width - wrapBuffer
		if wrapWidth < minContentWidth {
			wrapWidth = minContentWidth
		}

		lines := strings.Split(m.previewText, "\n")
		var wrappedLines []string

		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				wrappedLines = append(wrappedLines, "")
				continue
			}

			wrappedBlock := wordwrap.String(line, wrapWidth)
			blockLines := strings.Split(wrappedBlock, "\n")

			wrappedLines = append(wrappedLines, blockLines[0])

			for i := 1; i < len(blockLines); i++ {
				arrow := lipgloss.NewStyle().Foreground(subtleColor).Render("↳ ")
				wrappedLines = append(wrappedLines, arrow+blockLines[i])
			}
		}

		m.viewport.SetContent(strings.Join(wrappedLines, "\n"))
	}
}

func (m model) performSearch(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := m.searchFunc(query, searchLimit)
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

	title := titleStyle.Render(m.title)
	searchBoxWidth := m.width - layoutPadding
	searchView := inputBoxStyle.Width(searchBoxWidth).Render(m.textInput.View())

	if m.height < minWindowHeight {
		return fmt.Sprintf("\n%s\n%s\n%s", title, searchView, resizeWindowMsg)
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
		listItems = append(listItems, itemStyle.Foreground(subtleColor).Render(noMatchesMsg))
	}

	listView := lipgloss.NewStyle().
		Width(m.listWidth).
		Render(strings.Join(listItems, "\n"))

	if m.previewText == "" {
		m.viewport.SetContent(selectItemMsg)
	}

	previewView := previewStyle.
		Width(m.previewContainerWidth).
		Height(m.availableHeight).
		Render(m.viewport.View())

	content := lipgloss.JoinHorizontal(lipgloss.Top, listView, previewView)

	return fmt.Sprintf("\n%s\n%s\n%s", title, searchView, content)
}

func RunFuzzyFinder(title string, search SearchFunc, preview PreviewFunc) (string, error) {
	program := tea.NewProgram(InitialModel(title, search, preview), tea.WithAltScreen())
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}

	if resultModel, ok := finalModel.(model); ok && !resultModel.quitted {
		return resultModel.selected, nil
	}

	return "", nil
}
