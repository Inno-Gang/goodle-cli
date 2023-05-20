package model

import (
	"github.com/Inno-Gang/goodle-cli/key"
	"github.com/Inno-Gang/goodle-cli/tui/state"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

func (m *Model) Init() tea.Cmd {
	// TODO: move it to a separate function or so
	if login, ok := m.state.(*state.Login); ok {
		email := viper.GetString(key.AuthEmail)
		login.SetEmail(email)

		password := viper.GetString(key.AuthPassword)
		login.SetPassword(password)

		if email != "" && password != "" {
			return tea.Sequence(
				func() tea.Msg {
					return state.NewLoading("Hacking the mainframe...")
				},
				func() tea.Msg {
					client, err := login.Client(m)
					if err != nil {
						return err
					}

					courses, err := state.NewCourses(m.Context(), client)
					if err != nil {
						return err
					}

					return courses
				},
			)
		}
	}

	return m.state.Init(m)
}
