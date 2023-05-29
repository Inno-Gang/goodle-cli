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
	"github.com/inno-gang/goodle"
	"github.com/inno-gang/goodle/moodle"
	"github.com/inno-gang/goodle/richtext"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

var hiddenCourses = cache.New("hidden_courses")

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

	model *Courses
}

func (c coursesItem) FilterValue() string {
	return c.Course.Title()
}

func (c coursesItem) Title() string {
	title := c.FilterValue()
	titleWidth := int(float32(c.model.size.Width) * 0.75)

	if viper.GetBool(configKey.TUIShowEmoji) {
		return stringutil.Trim(title, titleWidth) + " " + stringutil.Correlate(title, icon.Emojis)
	}

	return stringutil.Trim(title, titleWidth)
}

func (c coursesItem) Description() string {
	return "Description\n" + c.Course.MoodleUrl()
}

type coursesKeyMap struct {
	list list.KeyMap
	OpenBrowser,
	OpenGrades,
	Confirm,
	ToggleHide,
	ToggleShowHidden key.Binding
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
			c.OpenGrades,
			c.ToggleHide,
			c.ToggleShowHidden,
		},
	}
}

type Courses struct {
	size       base.Size
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
		isHidden, _ := hiddenCourses.Get(fmt.Sprint(course.Id()), &cache.Empty{})

		return !isHidden
	})

	c := &Courses{
		client:  client,
		courses: courses,
		keyMap: coursesKeyMap{
			OpenBrowser:      util.Bind("open browser", "o"),
			OpenGrades:       util.Bind("open grades", "d"),
			Confirm:          util.Bind("confirm", "enter"),
			ToggleHide:       util.Bind("hide", "backspace"),
			ToggleShowHidden: util.Bind("show hidden", "H"),
		},
	}

	c.list = util.NewList(3, showable, func(course goodle.Course) list.Item {
		return coursesItem{course, c}
	})
	c.keyMap.list = c.list.KeyMap

	return c, nil
}

func (c *Courses) Intermediate() bool {
	return false
}

func (c *Courses) KeyMap() help.KeyMap {
	return c.keyMap
}

func (c *Courses) Title() base.Title {
	if c.showHidden {
		return base.Title{Text: "Hidden Courses", Background: color.Yellow}
	}

	return base.Title{Text: "Courses"}
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
	c.size = size
	c.list.SetSize(size.Width, size.Height)
}

func (c *Courses) toggleShowHidden() tea.Cmd {
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
		isHidden, _ := hiddenCourses.Get(fmt.Sprint(course.Id()), &cache.Empty{})

		show := (c.showHidden && isHidden) || (!c.showHidden && !isHidden)

		if show {
			courses = append(courses, coursesItem{course, c})
		}
	}

	c.list.Select(0)
	return c.list.SetItems(courses)
}

func (c *Courses) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	isFiltering := c.list.FilterState() == list.Filtering
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if isFiltering {
			goto end
		}

		switch {
		case key.Matches(msg, c.keyMap.OpenGrades):
			// https://moodle.innopolis.university/grade/report/user/index.php?id=<ID>
			item, ok := c.list.SelectedItem().(coursesItem)
			if !ok {
				return nil
			}

			URL := fmt.Sprintf(
				"https://moodle.innopolis.university/grade/report/user/index.php?id=%d",
				item.Id(),
			)

			return openWithDefaultApp(URL)
		case key.Matches(msg, c.keyMap.ToggleShowHidden):
			return c.toggleShowHidden()
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
		case key.Matches(msg, c.keyMap.OpenBrowser):
			item, ok := c.list.SelectedItem().(coursesItem)
			if !ok {
				return nil
			}

			return openWithDefaultApp(item.Course.MoodleUrl())
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

end:
	c.list, cmd = c.list.Update(msg)
	return cmd
}

func (c *Courses) View(base.Model) string {
	return c.list.View()
}

func (*Courses) Init(base.Model) tea.Cmd {
	return nil
}
