package state

import (
	"context"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/tuiutil"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/inno-gang/goodle"
	"github.com/inno-gang/goodle/moodle"
	"github.com/skratchdot/open-golang/open"
)

type coursesItem struct {
	goodle.Course
}

func (c coursesItem) FilterValue() string {
	return c.Course.Title()
}

func (c coursesItem) Title() string {
	return c.FilterValue()
}

func (c coursesItem) Description() string {
	return c.Course.MoodleUrl()
}

type coursesKeyMap struct {
	list        list.KeyMap
	OpenBrowser key.Binding
}

func (c coursesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{c.OpenBrowser, c.list.Filter, c.list.CursorUp, c.list.CursorDown}
}

func (c coursesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{c.ShortHelp()}
}

type Courses struct {
	client *moodle.Client
	list   list.Model
	keyMap coursesKeyMap
}

func NewCourses(ctx context.Context, client *moodle.Client) (*Courses, error) {
	courses, err := client.GetRecentCourses(ctx)
	if err != nil {
		return nil, err
	}

	var items = make([]list.Item, len(courses))
	for i, course := range courses {
		items[i] = coursesItem{course}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)

	l.KeyMap.Filter = tuiutil.Bind("filter", "ctrl+f")

	return &Courses{
		client: client,
		list:   l,
		keyMap: coursesKeyMap{
			OpenBrowser: tuiutil.Bind("open browser", "o"),
			list:        l.KeyMap,
		},
	}, nil
}

func (c *Courses) Intermediate() bool {
	return false
}

func (c *Courses) KeyMap() help.KeyMap {
	return c.keyMap
}

func (c *Courses) Title() string {
	return "Courses"
}

func (c *Courses) Resize(size base.Size) {
	c.list.SetSize(size.Width, size.Height)
}

func (c *Courses) Update(_ base.Model, msg tea.Msg) (cmd tea.Cmd) {
	isFiltering := c.list.FilterState() == list.Filtering
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keyMap.list.Filter):
			if c.list.FilterState() != list.Unfiltered {
				c.list.ResetFilter()
				return nil
			}
		case !isFiltering && key.Matches(msg, c.keyMap.OpenBrowser):
			item, ok := c.list.SelectedItem().(coursesItem)
			if !ok {
				return nil
			}

			return tea.Sequence(
				func() tea.Msg {
					return NewLoading("Opening...")
				},
				func() tea.Msg {
					err := open.Start(item.Course.MoodleUrl())
					if err != nil {
						return err
					}

					return base.MsgBack{}
				},
			)
		}
	}

	c.list, cmd = c.list.Update(msg)
	return cmd
}

func (c *Courses) View(base.Model) string {
	return c.list.View()
}

func (*Courses) Init(base.Model) tea.Cmd {
	return nil
}
