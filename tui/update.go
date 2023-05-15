package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Quit) && m.state.SoftQuit():
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.GoBack):
			m.back()
			return m, nil
		}
	case state:
		return m, m.pushState(msg)
	case error:
		if errors.Is(msg, context.Canceled) {
			return m, nil
		}

		cmd := m.pushState(&stateError{
			error: msg,
		})

		return m, cmd
	}

	return m.state.Update(m, msg)
}

func (s *stateError) Update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (s *stateLoading) Update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	s.spinner, cmd = s.spinner.Update(msg)
	return m, cmd
}

func (s *stateLogin) Update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	inputs := s.inputs()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Confirm):
			return m, func() tea.Msg {
				return fmt.Errorf("not implemented")
			}
		case key.Matches(msg, m.keyMap.FocusNext):
			for i, curr := range inputs[:len(inputs)-1] {
				next := inputs[i+1]

				if next.Focused() {
					curr = next
					next = inputs[0]
				}

				if curr.Focused() {
					curr.Blur()
					cmd = next.Focus()
				}
			}
		}
	}

	for _, n := range inputs {
		*n, cmd = n.Update(msg)
	}

	return m, cmd
}

func (s *stateCourseSelection) Update(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	s.list, cmd = s.list.Update(msg)
	return m, cmd
}
