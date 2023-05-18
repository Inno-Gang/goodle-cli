package base

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Title    lipgloss.Style
	TitleBar lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.
			NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1),
		TitleBar: lipgloss.
			NewStyle().
			Padding(0, 0, 1, 2),
	}
}
