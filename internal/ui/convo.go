package ui

import (
	"context"
	"fmt"
	"strings"

	"voc/internal/i18n"
	"voc/internal/llm"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	userBubbleStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("2")).
			Padding(0, 1)

	coachBubbleStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("4")).
				Padding(0, 1)

	corrBubbleStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("3")). // Yellow/Gold
			PaddingLeft(1).
			MarginLeft(2).
			Foreground(lipgloss.Color("3"))
)

type convoModel struct {
	client       llm.LLMClient
	progress     string
	history      []llm.Message
	lastCorr     []llm.Correction
	textInput    textinput.Model
	viewport     viewport.Model
	loading      bool
	spinner      spinner.Model
	width        int
	height       int
	err          error
	hostLang     string
	targetLang   string
}

func InitialConvoModel(client llm.LLMClient, progress string) convoModel {
	InitStyles()
	ti := textinput.New()
	ti.Placeholder = i18n.T(i18n.ConvoPlaceholder)
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	s := GetSpinner()

	return convoModel{
		client:     client,
		progress:   progress,
		textInput:  ti,
		spinner:    s,
		hostLang:   i18n.GetLanguageName(i18n.GetHostLanguage()),
		targetLang: i18n.GetLanguageName(i18n.GetTargetLanguage()),
	}
}

func (m convoModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

type chatMsg struct {
	response *llm.ChatResponse
	history  []llm.Message
	err      error
}

func (m convoModel) sendMessage() tea.Cmd {
	return func() tea.Msg {
		msg := m.textInput.Value()
		res, history, err := m.client.Chat(context.Background(), m.progress, m.history, msg, m.hostLang, m.targetLang)
		return chatMsg{response: res, history: history, err: err}
	}
}

func (m convoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Adjust height to leave room for title and fixed input area
		m.viewport = viewport.New(msg.Width, msg.Height-9)
		m.textInput.Width = msg.Width - 10
		if len(m.history) == 0 {
			m.viewport.SetContent(i18n.T(i18n.ConvoStarted))
		} else {
			m.viewport.SetContent(m.renderHistory())
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if !m.loading && m.textInput.Value() != "" {
				m.loading = true
				cmd := m.sendMessage()
				m.textInput.Reset()
				return m, cmd
			}
		}

	case chatMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.history = msg.history
		m.lastCorr = msg.response.Corrections
		m.viewport.SetContent(m.renderHistory())
		m.viewport.GotoBottom()

	case spinner.TickMsg:
		m.spinner, spCmd = m.spinner.Update(msg)
		return m, spCmd
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	if _, ok := msg.(tea.KeyMsg); !ok {
		m.viewport, vpCmd = m.viewport.Update(msg)
	}

	return m, tea.Batch(tiCmd, vpCmd, spCmd)
}

func (m convoModel) renderHistory() string {
	var historyView strings.Builder
	maxBubbleWidth := m.width - 20
	if maxBubbleWidth < 20 {
		maxBubbleWidth = 20
	}

	for i, h := range m.history {
		if h.Role == "user" {
			// Calculate if we need to wrap
			contentWidth := lipgloss.Width(h.Content)
			uStyle := userBubbleStyle
			if contentWidth > maxBubbleWidth {
				uStyle = uStyle.Width(maxBubbleWidth)
			}
			
			bubble := uStyle.Render(h.Content)
			// Manually align to the right
			paddedBubble := lipgloss.PlaceHorizontal(m.width-4, lipgloss.Right, bubble)
			historyView.WriteString(paddedBubble + "\n\n")

			// If this is the latest user message and we have corrections, show them here
			if i == len(m.history)-2 && len(m.lastCorr) > 0 {
				var corrBuilder strings.Builder
				corrBuilder.WriteString(i18n.T(i18n.ConvoCorrections) + "\n")
				for _, c := range m.lastCorr {
					corrBuilder.WriteString(fmt.Sprintf("• %s → %s\n", c.Incorrect, c.Correct))
				}
				corrStr := strings.TrimSpace(corrBuilder.String())

				cWidth := lipgloss.Width(corrStr)
				cStyle := corrBubbleStyle
				if cWidth > maxBubbleWidth {
					cStyle = cStyle.Width(maxBubbleWidth)
				}

				corr := cStyle.Render(corrStr)
				paddedCorr := lipgloss.PlaceHorizontal(m.width-4, lipgloss.Right, corr)
				historyView.WriteString(paddedCorr + "\n\n")
			}
		} else if h.Role == "model" {
			bubble := coachBubbleStyle.Width(maxBubbleWidth).Render(h.Content)
			historyView.WriteString(bubble + "\n\n")
		}
	}
	return historyView.String()
}

func (m convoModel) View() string {
	if m.err != nil {
		errStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Width(m.width - 4).
			Padding(1, 2)
		return fmt.Sprintf("\n%s\n\n%s\n\n%s", 
			titleStyle.Render(i18n.T(i18n.ConvoError)),
			errStyle.Render(m.err.Error()),
			hintStyle.Render(i18n.T(i18n.ConvoExitHint)))
	}

	title := titleStyle.Render(i18n.T(i18n.ConvoCoachTitle))
	history := m.viewport.View()
	
	// Stabilize input area height
	inputAreaHeight := 3
	var inputContent string
	if m.loading {
		inputContent = m.spinner.View() + " " + i18n.T(i18n.ConvoThinking)
	} else {
		inputContent = m.textInput.View()
	}

	input := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 0).
		Width(m.width).
		Height(inputAreaHeight).
		Render(inputContent)

	return fmt.Sprintf(
		"\n %s\n\n%s\n%s",
		title,
		history,
		input,
	)
}

func RunConvo(client llm.LLMClient, progress string) error {
	p := tea.NewProgram(InitialConvoModel(client, progress), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
