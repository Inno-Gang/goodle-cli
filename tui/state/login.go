package state

import (
	"github.com/Inno-Gang/goodle-cli/color"
	configKey "github.com/Inno-Gang/goodle-cli/key"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/util"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	loggedIn        bool
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
			confirm:   util.Bind("confirm", "enter"),
			focusNext: util.Bind("focus next", "tab"),
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

func (l *Login) focusNext() tea.Cmd {
	var cmds []tea.Cmd

	fields := l.textFields()

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

	return tea.Batch(cmds...)
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

func (*Login) Title() base.Title {
	return base.Title{Text: "Login"}
}

func (l *Login) Status() string {
	for _, f := range l.textFields() {
		if f.Err != nil {
			return lipgloss.NewStyle().Foreground(color.Red).Render(f.Err.Error())
		}
	}

	return ""
}

func (*Login) Backable() bool {
	return true
}

func (l *Login) KeyMap() help.KeyMap {
	return l.keyMap
}

func (l *Login) Resize(size base.Size) {
	for _, f := range l.textFields() {
		f.Width = size.Width
	}
}

func (l *Login) Update(m base.Model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, l.keyMap.confirm):
			return tea.Sequence(
				func() tea.Msg {
					return NewLoading("Logging in...")
				},
				func() tea.Msg {
					client, err := l.Client(m)
					if err != nil {
						return err
					}

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
			return l.focusNext()
		}
	}

	var cmds []tea.Cmd
	for _, f := range l.textFields() {
		var cmd tea.Cmd
		*f, cmd = f.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (l *Login) View(base.Model) string {
	return l.email.View() + "\n\n" + l.password.View()
}

func (l *Login) Init(model base.Model) tea.Cmd {
	email := viper.GetString(configKey.AuthEmail)
	l.SetEmail(email)

	password := viper.GetString(configKey.AuthPassword)
	l.SetPassword(password)

	if !l.loggedIn && email != "" && password != "" {
		l.loggedIn = true

		return tea.Sequence(
			func() tea.Msg {
				return NewLoading("Welcome back, " + lipgloss.NewStyle().Italic(true).Render(email))
			},
			func() tea.Msg {
				client, err := l.Client(model)
				if err != nil {
					return err
				}

				courses, err := NewCourses(model.Context(), client)
				if err != nil {
					return err
				}

				return courses
			},
		)
	}

	if !l.email.Focused() && !l.password.Focused() {
		return l.email.Focus()
	}

	return nil
}
