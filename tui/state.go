package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type state interface {
	View(*model) string
	Update(*model, tea.Msg) (tea.Model, tea.Cmd)
	Init(*model) tea.Cmd

	Resize(width, height int)

	Header() string
	SoftQuit() bool
	Intermediate() bool
}

type stateError struct {
	error error
}

func (s *stateError) SoftQuit() bool {
	return true
}

type stateLoading struct {
	spinner spinner.Model
	text    string
}

func (s *stateLoading) SoftQuit() bool {
	return false
}

type stateLogin struct {
	inputUsername, inputPassword textinput.Model
}

func (s *stateLogin) inputs() []*textinput.Model {
	return []*textinput.Model{
		&s.inputUsername,
		&s.inputPassword,
	}
}

func (s *stateLogin) SoftQuit() bool {
	return false
}

type stateCourseSelection struct {
	list list.Model
}

func (s *stateCourseSelection) SoftQuit() bool {
	return true
}
