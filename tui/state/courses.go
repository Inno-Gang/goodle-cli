package state

import (
	"context"
	"fmt"
	"github.com/Inno-Gang/goodle-cli/cache"
	"github.com/Inno-Gang/goodle-cli/color"
	"github.com/Inno-Gang/goodle-cli/icon"
	configKey "github.com/Inno-Gang/goodle-cli/key"
	"github.com/Inno-Gang/goodle-cli/stringutil"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/util"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inno-gang/goodle"
	"github.com/inno-gang/goodle/moodle"
	"github.com/inno-gang/goodle/richtext"
	"github.com/samber/lo"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
)

var (
	favoriteCourses = cache.New("favorite_courses")
	hiddenCourses   = cache.New("hidden_courses")
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
	title := c.FilterValue()

	if c.IsFavorite() {
		title += " " + lipgloss.NewStyle().Foreground(color.Yellow).Render(icon.Star)
	}

	return title
}

func (c coursesItem) IsFavorite() bool {
	found, err := favoriteCourses.Get(fmt.Sprint(c.Course.Id()), &cache.Empty{})
	return err == nil && found
}

func (c coursesItem) ToggleFavorite() error {
	id := fmt.Sprint(c.Course.Id())
	if c.IsFavorite() {
		return favoriteCourses.Delete(id)
	}

	return favoriteCourses.Set(id, cache.Empty{})
}

func (c coursesItem) Description() string {
	return c.Course.MoodleUrl()
}

type coursesKeyMap struct {
	list list.KeyMap
	OpenBrowser,
	Confirm,
	ToggleHide,
	ToggleShowHidden,
	ToggleFavorite key.Binding
}

func (c coursesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		c.Confirm,
		c.list.Filter,
		c.list.CursorUp,
		c.list.CursorDown,
	}
}

func (c coursesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		c.ShortHelp(),
		{
			c.OpenBrowser,
			c.ToggleHide,
			c.ToggleFavorite,
			c.ToggleShowHidden,
		},
	}
}

type Courses struct {
	client     *moodle.Client
	list       list.Model
	courses    []goodle.Course
	showHidden bool
	keyMap     coursesKeyMap
}

func NewCourses(ctx context.Context, client *moodle.Client) (*Courses, error) {
	courses, err := client.GetRecentCourses(ctx)
	if err != nil {
		return nil, err
	}

	// Filter out hidden courses
	showable := lo.Filter(courses, func(course goodle.Course, _ int) bool {
		found, _ := hiddenCourses.Get(fmt.Sprint(course.Id()), &cache.Empty{})
		return found
	})

	l := util.NewList(showable, func(course goodle.Course) list.Item {
		return coursesItem{course}
	})

	return &Courses{
		client:  client,
		list:    l,
		courses: courses,
		keyMap: coursesKeyMap{
			OpenBrowser:      util.Bind("open browser", "o"),
			Confirm:          util.Bind("confirm", "enter"),
			ToggleFavorite:   util.Bind("toggle favorite", "f"),
			ToggleHide:       util.Bind("hide", "backspace"),
			ToggleShowHidden: util.Bind("show hidden", "H"),
			list:             l.KeyMap,
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
	if c.showHidden {
		return "Hidden Courses"
	}

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

	if !c.showHidden {
		text += fmt.Sprintf(" %d hidden", len(c.courses)-len(c.list.Items()))
	}

	var filterValue string
	if value := c.list.FilterValue(); value != "" {
		filterValue = fmt.Sprintf(`"%s"`, value)
	}

	return fmt.Sprintf("%s %s %s", paginator, text, filterValue)
}

func (c *Courses) Resize(size base.Size) {
	c.list.SetSize(size.Width, size.Height)
}

func (c *Courses) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	isFiltering := c.list.FilterState() == list.Filtering
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !isFiltering {
			switch {
			case key.Matches(msg, c.keyMap.ToggleShowHidden):
				c.showHidden = !c.showHidden

				if c.showHidden {
					c.keyMap.ToggleHide.SetHelp(c.keyMap.ToggleHide.Keys()[0], "show")
					c.keyMap.ToggleShowHidden.SetHelp(c.keyMap.ToggleShowHidden.Keys()[0], "show visible")
				} else {
					c.keyMap.ToggleHide.SetHelp(c.keyMap.ToggleHide.Keys()[0], "hide")
					c.keyMap.ToggleShowHidden.SetHelp(c.keyMap.ToggleShowHidden.Keys()[0], "show hidden")
				}

				var courses []list.Item
				for _, course := range c.courses {
					found, _ := hiddenCourses.Get(fmt.Sprint(course.Id()), &cache.Empty{})

					show := (c.showHidden && !found) || (!c.showHidden && found)

					if show {
						courses = append(courses, coursesItem{course})
					}
				}

				return c.list.SetItems(courses)
			case key.Matches(msg, c.keyMap.ToggleHide):
				item, ok := c.list.SelectedItem().(coursesItem)
				if !ok {
					return nil
				}

				id := fmt.Sprint(item.Course.Id())
				found, _ := hiddenCourses.Get(id, &cache.Empty{})

				var err error
				if found {
					err = hiddenCourses.Delete(id)
				} else {
					err = hiddenCourses.Set(id, cache.Empty{})
				}

				if err != nil {
					return func() tea.Msg {
						return err
					}
				}

				index := c.list.Index()
				c.list.RemoveItem(index)

				visibleItems := len(c.list.VisibleItems())
				if visibleItems != 0 {
					if index == visibleItems {
						c.list.Select(index - 1)
					}
				} else {
					c.list.Select(0)
				}

				return nil
			case key.Matches(msg, c.keyMap.ToggleFavorite):
				item, ok := c.list.SelectedItem().(coursesItem)
				if !ok {
					return nil
				}

				err := item.ToggleFavorite()
				return func() tea.Msg {
					return err
				}
			case key.Matches(msg, c.keyMap.OpenBrowser):
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
			case key.Matches(msg, c.keyMap.Confirm):
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
