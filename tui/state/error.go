package state

import (
	"fmt"
	"github.com/Inno-Gang/goodle-cli/app"
	"github.com/Inno-Gang/goodle-cli/filesystem"
	"github.com/Inno-Gang/goodle-cli/icon"
	"github.com/Inno-Gang/goodle-cli/style"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/util"
	"github.com/Inno-Gang/goodle-cli/where"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/skratchdot/open-golang/open"
	"net/url"
	"runtime"
)

type errorKeyMap struct {
	quit, openIssue key.Binding
}

func (e errorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		e.quit,
		e.openIssue,
	}
}

func (e errorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{e.ShortHelp()}
}

type Error struct {
	error  error
	keyMap errorKeyMap
}

func NewError(err error) *Error {
	return &Error{
		error: err,
		keyMap: errorKeyMap{
			quit:      util.Bind("quit", "q"),
			openIssue: util.Bind("open issue", "o"),
		},
	}
}

func (*Error) Intermediate() bool {
	return true
}

func (*Error) Title() string {
	return "Error"
}

func (*Error) Status() string {
	return ""
}

func (*Error) Backable() bool {
	return true
}

func (e *Error) KeyMap() help.KeyMap {
	return e.keyMap
}

func (*Error) Resize(base.Size) {}

func (e *Error) Update(_ base.Model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, e.keyMap.quit):
			return tea.Quit
		case key.Matches(msg, e.keyMap.openIssue):
			URL := lo.Must(url.Parse("https://github.com/Inno-Gang/goodle-cli/issues/new"))
			values := url.Values{}

			logs, _ := filesystem.Api().ReadFile(where.LogFile())

			// TODO: add logs
			body := fmt.Sprintf(`%s

OS: %s
App version: %s

<details>
<summary>Logs</summary>

%s
%s
%s
</details>`,
				e.error.Error(),
				runtime.GOOS,
				app.Version,
				"```",
				logs,
				"```",
			)

			values.Set("title", "Error: "+errors.Cause(e.error).Error())
			values.Set("body", body)
			values.Set("labels", "bug")

			URL.RawQuery = values.Encode()

			return tea.Sequence(
				func() tea.Msg {
					return NewLoading("Opening...")
				},
				func() tea.Msg {
					err := open.Run(URL.String())
					if err != nil {
						return err
					}

					return base.MsgBack{}
				},
			)
		}
	}

	return nil
}

func (e *Error) View(base.Model) string {
	return style.Failure.Render(icon.Cross + " " + e.error.Error())
}

func (*Error) Init(base.Model) tea.Cmd {
	return nil
}
