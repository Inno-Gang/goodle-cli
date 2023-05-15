package tui

import (
	"context"
	"fmt"
	"github.com/Inno-Gang/goodle-cli/util"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zyedidia/generic/stack"
	"golang.org/x/term"
	"os"
	"strings"
)

type model struct {
	state   state
	history *stack.Stack[state]

	context           context.Context
	contextCancelFunc context.CancelFunc

	width, height int
	keyMap        *keyMap
}

func newModel() *model {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 40
	}

	var initialState state

	{
		inputUsername := textinput.New()
		inputUsername.Placeholder = "Username"
		inputUsername.Validate = func(s string) error {

			if strings.Contains(s, " ") {
				return fmt.Errorf("whitespaces not allowed")
			}

			if !util.IsASCII(s) {
				return fmt.Errorf("ASCII only")
			}

			return nil
		}

		inputPassword := textinput.New()
		inputPassword.EchoMode = textinput.EchoPassword
		inputPassword.Placeholder = "Password"

		initialState = &stateLogin{
			inputUsername: inputUsername,
			inputPassword: inputPassword,
		}
	}

	m := &model{
		width:   width,
		height:  height,
		state:   initialState,
		history: stack.New[state](),
	}

	m.context, m.contextCancelFunc = context.WithCancel(context.Background())
	m.keyMap = newKeyMap()

	return m
}

func (m *model) cancel() {
	m.contextCancelFunc()
	m.context, m.contextCancelFunc = context.WithCancel(context.Background())
}

func (m *model) back() {
	// do not pop the last state
	if m.history.Size() == 0 {
		return
	}

	m.cancel()
	m.state = m.history.Pop()

	// update size for old models
	m.resize(m.size())
}

func (m *model) pushState(s state) tea.Cmd {
	return func() tea.Msg {
		if !m.state.Intermediate() {
			m.history.Push(m.state)
		}

		m.state = s
		return m.state.Init(m)
	}
}
