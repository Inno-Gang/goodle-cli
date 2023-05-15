package tui

import (
	"fmt"
	"github.com/Inno-Gang/goodle-cli/icon"
	"github.com/Inno-Gang/goodle-cli/style"
)

func (m *model) View() string {
	if header := m.state.Header(); header != "" {
		return m.state.Header() + "\n\n" + m.state.View(m)
	}

	return m.state.View(m)
}

func (s *stateError) View(_ *model) string {
	return s.error.Error()
}

func (s *stateLoading) View(_ *model) string {
	return s.spinner.View() + " " + s.text
}

func (s *stateLogin) View(_ *model) string {
	return fmt.Sprintf(
		`%s
%s

%s`,

		s.inputUsername.View(),
		s.inputPassword.View(),
		style.Secondary.Render(icon.Info+" All your data is stored locally"),
	)
}

func (s *stateCourseSelection) View(_ *model) string {
	return s.list.View()
}
