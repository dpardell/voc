package ui

import (
	"strings"

	"voc/internal/i18n"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuOption int

const (
	OptionSearch MenuOption = iota
	OptionQuiz
	OptionConvo
	OptionList
	OptionInstall
	OptionQuit
)

func getMenuLabels() map[MenuOption]string {
	return map[MenuOption]string{
		OptionSearch:  i18n.T(i18n.MenuSearch),
		OptionQuiz:    i18n.T(i18n.MenuQuiz),
		OptionConvo:   i18n.T(i18n.MenuConvo),
		OptionList:    i18n.T(i18n.MenuVocab),
		OptionInstall: i18n.T(i18n.MenuInstall),
		OptionQuit:    i18n.T(i18n.MenuQuit),
	}
}

type splashModel struct {
	errorMsg      string
	dictInstalled bool
	cursor        MenuOption
	quitted       bool
	selected      MenuOption
	chosen        bool
	width         int
	height        int
}

func InitialSplashModel(errorMsg string, dictInstalled bool) splashModel {
	InitStyles()
	return splashModel{
		errorMsg:      errorMsg,
		dictInstalled: dictInstalled,
		cursor:        OptionSearch,
		selected:      -1,
	}
}

func (m splashModel) Init() tea.Cmd {
	return nil
}

func (m splashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// Clear error on any keypress
		m.errorMsg = ""

		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitted = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				// Skip Install if already installed
				if m.cursor == OptionInstall && m.dictInstalled {
					m.cursor--
				}
			}
		case "down", "j":
			if m.cursor < OptionQuit {
				m.cursor++
				// Skip Install if already installed
				if m.cursor == OptionInstall && m.dictInstalled {
					m.cursor++
				}
			}
		case "enter":
			if m.cursor == OptionQuit {
				m.quitted = true
				return m, tea.Quit
			}
			m.selected = m.cursor
			m.chosen = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m splashModel) View() string {
	if m.quitted || m.chosen || m.width == 0 {
		return ""
	}

	// Stylized Logo
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("4")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	logo := `
  _   _    ___     ____ 
 | | | |  / _ \   / ___|
 | | | | | | | | | |    
 \ \_/ / | |_| | | |___ 
  \___/   \___/   \____|
`
	// Tagline
	tagline := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		MarginBottom(2).
		Render(i18n.T(i18n.SplashTagline))

	// Error message if any
	var errBox string
	if m.errorMsg != "" {
		errBox = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("1")).
			Padding(0, 2).
			MarginBottom(1).
			Render("❌ " + m.errorMsg)
	}

	// Menu
	var menu strings.Builder
	labels := getMenuLabels()
	for i := OptionSearch; i <= OptionQuit; i++ {
		// Only show Install if NOT installed
		if i == OptionInstall && m.dictInstalled {
			continue
		}

		label := labels[i]

		icon := "  "
		switch i {
		case OptionSearch:
			icon = "🔍"
		case OptionQuiz:
			icon = "🎯"
		case OptionConvo:
			icon = "💬"
		case OptionList:
			icon = "📚"
		case OptionInstall:
			icon = "📥"
		case OptionQuit:
			icon = "👋"
		}

		if i == m.cursor {
			menu.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Bold(true).
				Render(" > "+icon+" "+label) + "\n")
		} else {
			menu.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Render("   "+icon+" "+label) + "\n")
		}
	}

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(2).
		Render(i18n.T(i18n.SplashFooter))

	var content string
	if errBox != "" {
		content = lipgloss.JoinVertical(lipgloss.Center,
			logoStyle.Render(logo),
			tagline,
			errBox,
			menu.String(),
			footer,
		)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center,
			logoStyle.Render(logo),
			tagline,
			menu.String(),
			footer,
		)
	}

	// Perfectly center the entire content block in the terminal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func RunSplash(errorMsg string, dictInstalled bool) (MenuOption, error) {
	p := tea.NewProgram(InitialSplashModel(errorMsg, dictInstalled), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return -1, err
	}
	res := m.(splashModel)
	if res.quitted {
		return OptionQuit, nil
	}
	return res.selected, nil
}
