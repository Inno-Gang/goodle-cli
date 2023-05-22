package state

import (
	"context"
	"fmt"
	"github.com/Inno-Gang/goodle-cli/color"
	configKey "github.com/Inno-Gang/goodle-cli/key"
	"github.com/Inno-Gang/goodle-cli/stringutil"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/util"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/inno-gang/goodle"
	"github.com/inno-gang/goodle/moodle"
	"github.com/inno-gang/goodle/richtext"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
)

type globalSection struct {
	sections []goodle.Section
	course   goodle.Course
}

func newGlobalSection(ctx context.Context, client *moodle.Client, course goodle.Course) (*globalSection, error) {
	sections, err := client.GetCourseSections(ctx, course.Id())
	if err != nil {
		return nil, err
	}

	return &globalSection{
		sections: sections,
		course:   course,
	}, nil
}

func (g globalSection) Id() int {
	return g.course.Id()
}

func (g globalSection) Title() string {
	return g.course.Title()
}

func (g globalSection) Description() *richtext.RichText {
	return g.course.Description()
}

func (g globalSection) Blocks() (blocks []goodle.Block) {
	for _, section := range g.sections {
		for _, block := range section.Blocks() {
			blocks = append(blocks, block)
		}
	}

	return blocks
}

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
	list                 list.KeyMap
	OpenBrowser, Confirm key.Binding
}

func (c coursesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{c.OpenBrowser, c.Confirm, c.list.Filter, c.list.CursorUp, c.list.CursorDown}
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

	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle.Foreground(color.Accent)
	delegate.Styles.SelectedDesc.Foreground(color.AccentDarken)

	delegate.Styles.SelectedTitle.BorderLeftForeground(color.Accent)
	delegate.Styles.SelectedDesc.BorderLeftForeground(color.Accent)

	l := list.New(items, delegate, 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	l.KeyMap.CancelWhileFiltering = util.Bind("cancel", "esc")

	return &Courses{
		client: client,
		list:   l,
		keyMap: coursesKeyMap{
			OpenBrowser: util.Bind("open browser", "o"),
			Confirm:     util.Bind("confirm", "enter"),
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

func (c *Courses) Backable() bool {
	return c.list.FilterState() == list.Unfiltered
}

func (c *Courses) Status() string {
	paginator := c.list.Paginator.View()
	text := stringutil.Quantify(
		len(c.list.VisibleItems()),
		"course",
		"courses",
	)

	return fmt.Sprintf("%s %s", paginator, text)
}

func (c *Courses) Resize(size base.Size) {
	c.list.SetSize(size.Width, size.Height)
}

func (c *Courses) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	isFiltering := c.list.FilterState() == list.Filtering
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
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
		case !isFiltering && key.Matches(msg, c.keyMap.Confirm):
			item, ok := c.list.SelectedItem().(coursesItem)
			if !ok {
				return nil
			}

			return tea.Sequence(
				func() tea.Msg {
					return NewLoading("Getting sections...")
				},
				func() tea.Msg {
					if viper.GetBool(configKey.TUIShowSections) {
						sections, err := NewSections(model.Context(), c.client, item.Course)
						if err != nil {
							return err
						}

						return sections
					}

					section, err := newGlobalSection(model.Context(), c.client, item.Course)
					if err != nil {
						return nil
					}

					return NewBlocks(section)
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
