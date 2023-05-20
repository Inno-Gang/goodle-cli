package base

import (
	"github.com/Inno-Gang/goodle-cli/color"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Title    lipgloss.Style
	TitleBar lipgloss.Style
	HelpBar  lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.
			NewStyle().
			Bold(true).
			Background(color.Background).
			Foreground(color.Foreground).
			Padding(0, 1),
		TitleBar: lipgloss.
			NewStyle().
			Padding(0, 0, 1, 2),
		HelpBar: lipgloss.
			NewStyle().
			Padding(0, 1),
	}
}
