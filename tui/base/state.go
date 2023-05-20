package base

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type State interface {
	Intermediate() bool
	KeyMap() help.KeyMap
	Title() string
	Status() string

	Resize(size Size)

	Update(model Model, msg tea.Msg) tea.Cmd
	View(model Model) string
	Init(model Model) tea.Cmd
}
