package state

import (
	"fmt"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/tuiutil"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type loginKeyMap struct {
	confirm, focusNext key.Binding
}

func (l loginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		l.confirm,
		l.focusNext,
	}
}

func (l loginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{l.ShortHelp()}
}

type Login struct {
	username, password textinput.Model
	keyMap             loginKeyMap
}

func NewLogin() *Login {
	username := textinput.New()
	username.Placeholder = "Username"
	username.Validate = func(s string) error {
		if strings.Contains(s, " ") {
			return fmt.Errorf("whitespaces is not permitted")
		}

		return nil
	}

	password := textinput.New()
	password.EchoMode = textinput.EchoPassword
	password.Placeholder = "Password"

	return &Login{
		username: username,
		password: password,
		keyMap: loginKeyMap{
			confirm:   tuiutil.Bind("confirm", "enter"),
			focusNext: tuiutil.Bind("focus next", "tab"),
		},
	}
}

func (l *Login) textFields() []*textinput.Model {
	return []*textinput.Model{
		&l.username,
		&l.password,
	}
}

func (*Login) Intermediate() bool {
	return false
}

func (*Login) Title() string {
	return "Login"
}

func (l *Login) KeyMap() help.KeyMap {
	return l.keyMap
}

func (*Login) Resize(base.Size) {}

func (l *Login) Update(model base.Model, msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	fields := l.textFields()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, l.keyMap.confirm):
			return tea.Sequence(
				func() tea.Msg {
					return NewLoading("Launching nukes...")
				},
				func() tea.Msg {
					// TODO: remove this. for testing purposes only
					select {
					case <-model.Context().Done():
					case <-time.After(time.Second * 2):
						return fmt.Errorf("not implemented")
					}

					return nil
				})
		case key.Matches(msg, l.keyMap.focusNext):
			for i, curr := range fields[:len(fields)-1] {
				next := fields[i+1]

				if next.Focused() {
					curr = next
					next = fields[0]
				}

				if curr.Focused() {
					curr.Blur()
					cmds = append(cmds, next.Focus())
				}
			}

			return nil
		}
	}

	for _, f := range fields {
		var cmd tea.Cmd
		*f, cmd = f.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (l *Login) View(base.Model) string {
	return l.username.View() + "\n\n" + l.password.View()
}

func (l *Login) Init(base.Model) tea.Cmd {
	return l.username.Focus()
}
