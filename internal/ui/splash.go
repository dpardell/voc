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
	OptionQuit
)

func getMenuLabels() map[MenuOption]string {
	return map[MenuOption]string{
		OptionSearch: i18n.T(i18n.MenuSearch),
		OptionQuiz:   i18n.T(i18n.MenuQuiz),
		OptionConvo:  i18n.T(i18n.MenuConvo),
		OptionList:   i18n.T(i18n.MenuVocab),
		OptionQuit:   i18n.T(i18n.MenuQuit),
	}
}

type splashModel struct {
	saying   string
	cursor   MenuOption
	quitted  bool
	selected MenuOption
	chosen   bool
	width    int
	height   int
}

func InitialSplashModel(saying string) splashModel {
	InitStyles()
	return splashModel{
		saying:   saying,
		cursor:   OptionSearch,
		selected: -1,
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
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitted = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < OptionQuit {
				m.cursor++
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
  _   _   ___   ____ 
 | | | | / _ \ / ___|
 | | | || | | || |    
 \ \_/ /| |_| || |___ 
  \___/  \___/  \____|
`
	// Tagline
	tagline := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		MarginBottom(2).
		Render(i18n.T(i18n.SplashTagline))

	// Saying Card
	sayingCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		MarginBottom(2).
		Width(50).
		Align(lipgloss.Center).
		Render(lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("246")).Render(m.saying))

	// Menu
	var menu strings.Builder
	labels := getMenuLabels()
	for i := OptionSearch; i <= OptionQuit; i++ {
		label := labels[i]
		
		icon := "  "
		switch i {
		case OptionSearch: icon = "🔍"
		case OptionQuiz:   icon = "🎯"
		case OptionConvo:  icon = "💬"
		case OptionList:   icon = "📚"
		case OptionQuit:   icon = "👋"
		}

		if i == m.cursor {
			menu.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Bold(true).
				Render(" > " + icon + " " + label) + "\n")
		} else {
			menu.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Render("   " + icon + " " + label) + "\n")
		}
	}

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(2).
		Render(i18n.T(i18n.SplashFooter))

	content := lipgloss.JoinVertical(lipgloss.Center,
		logoStyle.Render(logo),
		tagline,
		sayingCard,
		menu.String(),
		footer,
	)

	// Perfectly center the entire content block in the terminal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func RunSplash(saying string) (MenuOption, error) {
	p := tea.NewProgram(InitialSplashModel(saying), tea.WithAltScreen())
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
