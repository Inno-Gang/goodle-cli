package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *model) Init() tea.Cmd {
	return m.state.Init(m)
}

func (s *stateError) Init(_ *model) tea.Cmd {
	return nil
}

func (s *stateLoading) Init(_ *model) tea.Cmd {
	return s.spinner.Tick
}

func (s *stateLogin) Init(_ *model) tea.Cmd {
	return s.inputUsername.Focus()
}
