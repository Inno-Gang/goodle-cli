package state

import (
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type loadingKeyMap struct{}

func (loadingKeyMap) ShortHelp() []key.Binding {
	return nil
}

func (loadingKeyMap) FullHelp() [][]key.Binding {
	return nil
}

type Loading struct {
	message string
	spinner spinner.Model
	keyMap  loadingKeyMap
}

func NewLoading(message string) *Loading {
	return &Loading{
		message: message,
		spinner: spinner.New(),
		keyMap:  loadingKeyMap{},
	}
}

func (*Loading) Intermediate() bool {
	return true
}

func (l *Loading) KeyMap() help.KeyMap {
	return l.keyMap
}

func (*Loading) Title() string {
	return "Loading"
}

func (l *Loading) Resize(base.Size) {}

func (l *Loading) Update(_ base.Model, msg tea.Msg) (cmd tea.Cmd) {
	l.spinner, cmd = l.spinner.Update(msg)
	return cmd
}

func (l *Loading) View(base.Model) string {
	return l.spinner.View() + " " + l.message
}

func (l *Loading) Init(base.Model) tea.Cmd {
	// TODO: this message is lost somewhere, wtf
	return l.spinner.Tick
}
