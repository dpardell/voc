package ui

import (
	"math/rand/v2"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	subtleColor    lipgloss.Color
	highlightColor lipgloss.Color
	titleBgColor   lipgloss.Color

	titleStyle        lipgloss.Style
	inputBoxStyle     lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
	previewStyle      lipgloss.Style
	subtleStyle       lipgloss.Style
	hintStyle         lipgloss.Style
	savedStyle        lipgloss.Style

	correctStyle lipgloss.Style
	wrongStyle   lipgloss.Style
)

type Theme struct {
	Subtle    lipgloss.Color
	Highlight lipgloss.Color
	TitleBg   lipgloss.Color
}

var DefaultTheme = Theme{
	Subtle:    lipgloss.Color("241"),
	Highlight: lipgloss.Color("2"),
	TitleBg:   lipgloss.Color("4"),
}

var NightTheme = Theme{
	Subtle:    lipgloss.Color("241"),
	Highlight: lipgloss.Color("229"),
	TitleBg:   lipgloss.Color("244"),
}

func GetTheme() Theme {
	if os.Getenv("VOC_THEME") == "night" {
		return NightTheme
	}
	// Fallback to time-based for backward compatibility if not explicitly set
	hour := time.Now().Hour()
	if hour >= 21 || hour < 7 {
		return NightTheme
	}
	return DefaultTheme
}

func InitStyles() {
	theme := GetTheme()
	subtleColor = theme.Subtle
	highlightColor = theme.Highlight
	titleBgColor = theme.TitleBg
	// ... rest of the styles use these variables

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")). // White
		Background(titleBgColor).
		Padding(0, 1).
		MarginLeft(1).
		Bold(true)

	inputBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(titleBgColor).
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

	subtleStyle = lipgloss.NewStyle().
		Foreground(subtleColor).
		MarginLeft(2)

	hintStyle = lipgloss.NewStyle().
		Foreground(subtleColor).
		MarginTop(1)

	savedStyle = lipgloss.NewStyle().
		Foreground(highlightColor).
		Bold(true).
		MarginLeft(2)

	correctStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Bold(true)

	wrongStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true)
}

func GetSpinner() spinner.Model {
	hour := time.Now().Hour()
	isNight := hour >= 21 || hour < 7

	s := spinner.New()
	if isNight {
		s.Spinner = spinner.Spinner{
			Frames: []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"},
			FPS:    time.Second / 8,
		}
	} else {
		if rand.Float64() < 0.1 { // 10% chance of Sun spinner during the day
			s.Spinner = spinner.Spinner{
				Frames: []string{"☀️", "🌤️", "⛅️", "🌥️", "☁️", "🌥️", "⛅️", "🌤️", "☀️"},
				FPS:    time.Second / 4,
			}
		} else {
			s.Spinner = spinner.Spinner{
				Frames: []string{"🌍", "🌎", "🌏"},
				FPS:    time.Second / 4,
			}
		}
	}
	return s
}
