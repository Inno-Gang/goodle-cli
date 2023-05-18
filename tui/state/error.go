package state

import (
	"github.com/Inno-Gang/goodle-cli/icon"
	"github.com/Inno-Gang/goodle-cli/style"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/tuiutil"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type errorKeyMap struct {
	quit key.Binding
}

func (e errorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		e.quit,
	}
}

func (e errorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{e.ShortHelp()}
}

type Error struct {
	error  error
	keyMap errorKeyMap
}

func NewError(err error) *Error {
	return &Error{
		error: err,
		keyMap: errorKeyMap{
			quit: tuiutil.Bind("quit", "q"),
		},
	}
}

func (*Error) Intermediate() bool {
	return true
}

func (*Error) Title() string {
	return "Error"
}

func (e *Error) KeyMap() help.KeyMap {
	return e.keyMap
}

func (*Error) Resize(base.Size) {}

func (e *Error) Update(_ base.Model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, e.keyMap.quit):
			return tea.Quit
		}
	}

	return nil
}

func (e *Error) View(base.Model) string {
	return style.Failure(icon.Cross + " " + e.error.Error())
}

func (*Error) Init(base.Model) tea.Cmd {
	return nil
}
