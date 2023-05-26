package state

import (
	"github.com/Inno-Gang/goodle-cli/color"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type loadingKeyMap struct{}

func (loadingKeyMap) ShortHelp() []key.Binding {
	return nil
}

func (loadingKeyMap) FullHelp() [][]key.Binding {
	return nil
}

type Loading struct {
	message string
	spinner spinner.Model
	keyMap  loadingKeyMap
}

func NewLoading(message string) *Loading {
	return &Loading{
		message: message,
		spinner: spinner.New(spinner.WithSpinner(spinner.Spinner{
			Frames: []string{
				"Moodle",
				"mOodle",
				"moOdle",
				"mooDle",
				"moodLe",
				"moodlE",
			},
			FPS: time.Second / 7,
		})),
		keyMap: loadingKeyMap{},
	}
}

func (*Loading) Intermediate() bool {
	return true
}

func (l *Loading) KeyMap() help.KeyMap {
	return l.keyMap
}

func (*Loading) Title() base.Title {
	return base.Title{
		Text:       "Loading",
		Background: color.Yellow,
	}
}

func (l *Loading) Status() string {
	return lipgloss.NewStyle().Foreground(color.Accent).Render(l.spinner.View())
}

func (*Loading) Backable() bool {
	return true
}

func (l *Loading) Resize(base.Size) {}

func (l *Loading) Update(_ base.Model, msg tea.Msg) (cmd tea.Cmd) {
	l.spinner, cmd = l.spinner.Update(msg)
	return cmd
}

func (l *Loading) View(base.Model) string {
	return l.message
}

func (l *Loading) Init(base.Model) tea.Cmd {
	return l.spinner.Tick
}
