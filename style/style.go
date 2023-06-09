package style

import (
	"github.com/Inno-Gang/goodle-cli/color"
	"github.com/charmbracelet/lipgloss"
)

var (
	Success = lipgloss.NewStyle().Foreground(color.Green)
	Failure = lipgloss.NewStyle().Foreground(color.Red)
	Warning = lipgloss.NewStyle().Foreground(color.Yellow)
)

var (
	Accent    = lipgloss.NewStyle().Bold(true).Foreground(color.Purple)
	Secondary = lipgloss.NewStyle().Faint(true).Italic(true)
)
