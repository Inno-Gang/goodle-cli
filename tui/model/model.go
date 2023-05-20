package model

import (
	"context"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/state"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/zyedidia/generic/stack"
	"golang.org/x/term"
	"os"
)

type Model struct {
	state   base.State
	history *stack.Stack[base.State]

	context           context.Context
	contextCancelFunc context.CancelFunc

	size base.Size

	styles base.Styles

	keyMap *keyMap
	help   help.Model
}

func (m *Model) ShortHelp() []key.Binding {
	keys := []key.Binding{m.keyMap.Back}
	return append(keys, m.state.KeyMap().ShortHelp()...)
}

func (m *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

func New() *Model {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 40
	}

	model := &Model{
		state:   state.NewLogin(),
		history: stack.New[base.State](),
		size: base.Size{
			Width:  width,
			Height: height,
		},
		keyMap: newKeyMap(),
		help:   help.New(),
		styles: base.DefaultStyles(),
	}

	defer model.resize(model.StateSize())

	model.context, model.contextCancelFunc = context.WithCancel(context.Background())

	return model
}

func (m *Model) StateSize() base.Size {
	return base.Size{
		Width: m.size.Width,
		// TODO: should be calculated dynamically
		Height: m.size.Height - 3,
	}
}

func (m *Model) Context() context.Context {
	return m.context
}

func (m *Model) cancel() {
	m.contextCancelFunc()
	m.context, m.contextCancelFunc = context.WithCancel(context.Background())
}

func (m *Model) resize(size base.Size) {
	m.size = size
	m.state.Resize(m.StateSize())
}

func (m *Model) back() tea.Cmd {
	// do not pop the last state
	if m.history.Size() == 0 {
		return nil
	}

	log.Info("going to the previous state", "state", m.history.Peek().Title())

	m.cancel()
	m.state = m.history.Pop()

	// update size for old models
	m.state.Resize(m.StateSize())

	return m.state.Init(m)
}

func (m *Model) pushState(state base.State) tea.Cmd {
	log.Info("new state", "state", state.Title())
	if !m.state.Intermediate() {
		m.history.Push(m.state)
	}

	m.state = state
	m.state.Resize(m.StateSize())

	return m.state.Init(m)
}
