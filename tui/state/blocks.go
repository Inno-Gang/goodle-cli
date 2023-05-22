package state

import (
	"fmt"
	"github.com/Inno-Gang/goodle-cli/color"
	"github.com/Inno-Gang/goodle-cli/stringutil"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/util"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/inno-gang/goodle"
	"github.com/skratchdot/open-golang/open"
)

type blocksItem struct {
	goodle.Block
}

func (b blocksItem) FilterValue() string {
	return b.Block.Title()
}

func (b blocksItem) Title() string {
	return b.FilterValue()
}

func (b blocksItem) Description() string {
	return b.Block.Type().Name()
}

type blocksKeyMap struct {
	openBrowser, open key.Binding
	list              list.KeyMap
}

func (b blocksKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		b.openBrowser,
		b.open,
		b.list.Filter,
		b.list.CursorUp,
		b.list.CursorDown,
	}
}

func (b blocksKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{b.ShortHelp()}
}

type Blocks struct {
	section goodle.Section
	list    list.Model
	keyMap  blocksKeyMap
}

func NewBlocks(section goodle.Section) *Blocks {
	blocks := section.Blocks()

	var items = make([]list.Item, len(blocks))
	for i, block := range blocks {
		items[i] = blocksItem{block}
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

	return &Blocks{
		section: section,
		list:    l,
		keyMap: blocksKeyMap{
			openBrowser: util.Bind("open browser", "o"),
			open:        util.Bind("open", "enter"),
			list:        l.KeyMap,
		},
	}
}

func (b *Blocks) Intermediate() bool {
	return false
}

func (b *Blocks) KeyMap() help.KeyMap {
	return b.keyMap
}

func (b *Blocks) Title() string {
	return b.section.Title()
}

func (b *Blocks) Status() string {
	paginator := b.list.Paginator.View()
	text := stringutil.Quantify(
		len(b.list.VisibleItems()),
		"block",
		"blocks",
	)

	return fmt.Sprintf("%s %s", paginator, text)
}

func (b *Blocks) Backable() bool {
	return b.list.FilterState() == list.Unfiltered
}

func (b *Blocks) Resize(size base.Size) {
	b.list.SetSize(size.Width, size.Height)
}

func (b *Blocks) openSelected() tea.Cmd {
	return b.openSelectedInBrowser()

	// TODO: waiting for moodle api library update

	//item, ok := b.list.SelectedItem().(blocksItem)
	//if !ok {
	//	return nil
	//}
	//
	//switch item.Block.Type() {
	//case goodle.BlockTypeFile:
	//	block := item.Block.(goodle.BlockFile)
	//	mime := mimetype.Lookup(block.MimeType())
	//	if mime == nil {
	//		return b.openSelectedInBrowser()
	//	}
	//
	//	path := filepath.Join(where.Temp(), strconv.Itoa(block.Id())+mime.Extension())
	//
	//	return tea.Sequence(
	//		func() tea.Msg {
	//			return NewLoading(fmt.Sprintf("Downloading %s", block.Title()))
	//		},
	//		func() tea.Msg {
	//			log.Info("downloading file", "url", block.DownloadUrl())
	//			res, err := http.Get(block.DownloadUrl())
	//			if err != nil {
	//				return nil
	//			}
	//
	//			if res.StatusCode != http.StatusOK {
	//				return res.Status
	//			}
	//
	//			err = filesystem.Api().WriteReader(path, res.Body)
	//			if err != nil {
	//				return nil
	//			}
	//
	//			err = open.Run(path)
	//			if err != nil {
	//				return nil
	//			}
	//
	//			return base.MsgBack{}
	//		},
	//	)
	//default:
	//	return b.openSelectedInBrowser()
	//}
}

func (b *Blocks) openSelectedInBrowser() tea.Cmd {
	item, ok := b.list.SelectedItem().(blocksItem)
	if !ok {
		return nil
	}

	return tea.Sequence(
		func() tea.Msg {
			return NewLoading("opening...")
		},
		func() tea.Msg {
			err := open.Run(item.MoodleUrl())
			if err != nil {
				return nil
			}

			return base.MsgBack{}
		},
	)
}

func (b *Blocks) Update(_ base.Model, msg tea.Msg) (cmd tea.Cmd) {
	isFiltering := b.list.FilterState() == list.Filtering
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case !isFiltering && key.Matches(msg, b.keyMap.openBrowser):
			return b.openSelectedInBrowser()
		case !isFiltering && key.Matches(msg, b.keyMap.open):
			return b.openSelected()
		}
	}
	b.list, cmd = b.list.Update(msg)
	return cmd
}

func (b *Blocks) View(base.Model) string {
	return b.list.View()
}

func (b *Blocks) Init(base.Model) tea.Cmd {
	return nil
}
