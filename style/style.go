package style

import (
	"github.com/Inno-Gang/goodle-cli/color"
	"github.com/charmbracelet/lipgloss"
)

var (
	Success = lipgloss.NewStyle().Foreground(color.Green).Render
	Failure = lipgloss.NewStyle().Foreground(color.Red).Render
	Warning = lipgloss.NewStyle().Foreground(color.Yellow).Render
)

var (
	Accent    = lipgloss.NewStyle().Bold(true).Foreground(color.Purple)
	Secondary = lipgloss.NewStyle().Faint(true).Italic(true)
)
