package state

import (
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type errorKeyMap struct{}

func (errorKeyMap) ShortHelp() []key.Binding {
	return nil
}

func (e errorKeyMap) FullHelp() [][]key.Binding {
	return nil
}

type Error struct {
	error  error
	keyMap errorKeyMap
}

func NewError(err error) *Error {
	return &Error{
		error:  err,
		keyMap: errorKeyMap{},
	}
}

func (*Error) Intermediate() bool {
	return true
}

func (*Error) Header() string {
	return "error"
}

func (e *Error) KeyMap() help.KeyMap {
	return e.keyMap
}

func (*Error) Resize(base.Size) {}

func (*Error) Update(base.Model, tea.Msg) tea.Cmd {
	return nil
}

func (e *Error) View(base.Model) string {
	return e.error.Error()
}

func (*Error) Init(base.Model) tea.Cmd {
	return nil
}
