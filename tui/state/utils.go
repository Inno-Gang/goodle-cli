package state

import (
	"github.com/Inno-Gang/goodle-cli/tui/base"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skratchdot/open-golang/open"
)

func openWithDefaultApp(input string) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return NewLoading("opening...")
		},
		func() tea.Msg {
			err := open.Run(input)
			if err != nil {
				return nil
			}

			return base.MsgBack{}
		},
	)
}
