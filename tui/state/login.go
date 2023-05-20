package state

import (
	configKey "github.com/Inno-Gang/goodle-cli/key"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/tuiutil"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/inno-gang/goodle/auth"
	"github.com/inno-gang/goodle/moodle"
	"github.com/spf13/viper"
	"strings"
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
	email, password textinput.Model
	keyMap          loginKeyMap
	authenticator   *auth.IuAuthenticator
}

func NewLogin() *Login {
	authenticator := auth.NewIuAuthenticator()

	username := textinput.New()
	username.Placeholder = "Email"

	password := textinput.New()
	password.EchoMode = textinput.EchoPassword
	password.Placeholder = "Password"

	return &Login{
		authenticator: authenticator,
		email:         username,
		password:      password,
		keyMap: loginKeyMap{
			confirm:   tuiutil.Bind("confirm", "enter"),
			focusNext: tuiutil.Bind("focus next", "tab"),
		},
	}
}

func (l *Login) SetEmail(email string) {
	l.email.SetValue(email)
}

func (l *Login) SetPassword(password string) {
	l.password.SetValue(password)
}

func (l *Login) textFields() []*textinput.Model {
	return []*textinput.Model{
		&l.email,
		&l.password,
	}
}

func (l *Login) credentials() auth.IuCredentials {
	return auth.IuCredentials{
		Email:    strings.TrimSpace(l.email.Value()),
		Password: l.password.Value(),
	}
}

func (l *Login) saveCredentials() error {
	credentials := l.credentials()

	viper.Set(configKey.AuthEmail, credentials.Email)

	if viper.GetBool(configKey.AuthRemember) {
		viper.Set(configKey.AuthPassword, credentials.Password)
	}

	switch err := viper.WriteConfig(); err.(type) {
	case viper.ConfigFileNotFoundError:
		return viper.SafeWriteConfig()
	default:
		return err
	}
}

func (l *Login) Client(m base.Model) (*moodle.Client, error) {
	return l.authenticator.Authenticate(m.Context(), l.credentials())
}

func (*Login) Intermediate() bool {
	return false
}

func (*Login) Title() string {
	return "Login"
}

func (l *Login) Status() string {
	return ""
}

func (l *Login) KeyMap() help.KeyMap {
	return l.keyMap
}

func (*Login) Resize(base.Size) {}

func (l *Login) Update(m base.Model, msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	fields := l.textFields()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, l.keyMap.confirm):
			return tea.Sequence(
				func() tea.Msg {
					return NewLoading("Hacking the mainframe...")
				},
				func() tea.Msg {
					client, err := l.Client(m)
					if err != nil {
						return err
					}

					// TODO: log errors
					err = l.saveCredentials()
					if err != nil {
						log.Error("failed to login", "err", err.Error())
					}

					newState, err := NewCourses(m.Context(), client)
					if err != nil {
						return err
					}

					return newState
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
	return l.email.View() + "\n\n" + l.password.View()
}

func (l *Login) Init(base.Model) tea.Cmd {
	return l.email.Focus()
}
