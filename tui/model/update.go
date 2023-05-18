package model

import (
	"context"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/state"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(base.Size{
			Width:  msg.Width,
			Height: msg.Height,
		})

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Back):
			m.back()
			return m, nil
		}
	case base.MsgBack:
		m.back()
		return m, nil
	case base.State:
		return m, m.pushState(msg)
	case error:
		if errors.Is(msg, context.Canceled) {
			return m, nil
		}

		return m, m.pushState(state.NewError(msg))
	}

	cmd := m.state.Update(m, msg)
	return m, cmd
}
