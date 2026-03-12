package ui

import (
	"os"
	"testing"
)

func TestGetTheme(t *testing.T) {
	os.Setenv("VOC_THEME", "night")
	defer os.Unsetenv("VOC_THEME")

	theme := GetTheme()
	if theme.Highlight == "" || theme.Subtle == "" || theme.TitleBg == "" {
		t.Error("night theme has empty color fields")
	}

	os.Unsetenv("VOC_THEME")
	theme = GetTheme()
	if theme.Highlight == "" || theme.Subtle == "" || theme.TitleBg == "" {
		t.Error("default (time-based) theme has empty color fields")
	}
}

func TestGetSpinner(t *testing.T) {
	s := GetSpinner()
	if len(s.Spinner.Frames) == 0 {
		t.Error("GetSpinner returned a spinner with no frames")
	}
}
