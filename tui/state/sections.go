package state

import (
	"context"
	"fmt"
	"github.com/Inno-Gang/goodle-cli/color"
	"github.com/Inno-Gang/goodle-cli/stringutil"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/tuiutil"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/inno-gang/goodle"
	"github.com/inno-gang/goodle/moodle"
)

type sectionsItem struct {
	goodle.Section
}

func (s sectionsItem) FilterValue() string {
	return s.Section.Title()
}

func (s sectionsItem) Title() string {
	return s.FilterValue()
}

func (s sectionsItem) Description() string {
	return stringutil.Quantify(
		len(s.Section.Blocks()),
		"block",
		"blocks",
	)
}

type sectionsKeyMap struct {
	confirm key.Binding
	list    list.KeyMap
}

func (s sectionsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		s.confirm,
		s.list.Filter,
		s.list.CursorUp,
		s.list.CursorDown,
	}
}

func (s sectionsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{s.ShortHelp()}
}

type Sections struct {
	course goodle.Course
	list   list.Model
	keyMap sectionsKeyMap
}

func NewSections(
	ctx context.Context,
	client *moodle.Client,
	course goodle.Course,
) (*Sections, error) {
	sections, err := client.GetCourseSections(ctx, course.Id())
	if err != nil {
		return nil, err
	}

	var items = make([]list.Item, len(sections))
	for i, section := range sections {
		items[i] = sectionsItem{section}
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

	l.KeyMap.CancelWhileFiltering = tuiutil.Bind("cancel", "esc")

	return &Sections{
		course: course,
		list:   l,
		keyMap: sectionsKeyMap{
			confirm: tuiutil.Bind("confirm", "enter"),
			list:    l.KeyMap,
		},
	}, nil
}

func (s *Sections) Intermediate() bool {
	return false
}

func (s *Sections) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *Sections) Title() string {
	return s.course.Title()
}

func (s *Sections) Status() string {
	paginator := s.list.Paginator.View()
	text := stringutil.Quantify(
		len(s.list.VisibleItems()),
		"section",
		"sections",
	)

	return fmt.Sprintf("%s %s", paginator, text)
}

func (s *Sections) Backable() bool {
	return s.list.FilterState() == list.Unfiltered
}

func (s *Sections) Resize(size base.Size) {
	s.list.SetSize(size.Width, size.Height)
}

func (s *Sections) Update(_ base.Model, msg tea.Msg) (cmd tea.Cmd) {
	isFiltering := s.list.FilterState() == list.Filtering
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case !isFiltering && key.Matches(msg, s.keyMap.confirm):
			item, ok := s.list.SelectedItem().(sectionsItem)
			if !ok {
				return nil
			}

			return func() tea.Msg {
				return NewBlocks(item.Section)
			}
		}
	}
	s.list, cmd = s.list.Update(msg)
	return cmd
}

func (s *Sections) View(base.Model) string {
	return s.list.View()
}

func (s *Sections) Init(base.Model) tea.Cmd {
	return nil
}