package tui

import (
	"github.com/Inno-Gang/goodle-cli/tui/model"
	tea "github.com/charmbracelet/bubbletea"
)

func Run() error {
	program := tea.NewProgram(model.New(), tea.WithAltScreen())

	_, err := program.Run()
	return err
}
