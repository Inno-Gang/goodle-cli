package state

import (
	"fmt"
	configKey "github.com/Inno-Gang/goodle-cli/key"
	"github.com/Inno-Gang/goodle-cli/stringutil"
	"github.com/Inno-Gang/goodle-cli/tui/base"
	"github.com/Inno-Gang/goodle-cli/tui/util"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/gabriel-vasile/mimetype"
	"github.com/inno-gang/goodle"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"strings"
)

const moodleFormat = "Monday, 02 Jan 2006, 15:04"

type blocksItem struct {
	goodle.Block
}

func (b blocksItem) FilterValue() string {
	return b.Block.Title()
}

func (b blocksItem) Title() string {
	title := b.FilterValue()

	if !viper.GetBool(configKey.TUIShowEmoji) || b.Type() == goodle.BlockTypeUnknown {
		return title
	}

	var emoji string
	switch b.Type() {
	case goodle.BlockTypeQuiz:
		emoji = "📊"
	case goodle.BlockTypeLink:
		emoji = "🖇️"
	case goodle.BlockTypeFile:
		emoji = "📃"
	case goodle.BlockTypeFolder:
		emoji = "📁"
	case goodle.BlockTypeAssignment:
		emoji = "📥"
	}

	return title + " " + emoji
}

func (b blocksItem) Description() string {
	const whitespace = ' '
	var info strings.Builder

	switch b.Type() {
	case goodle.BlockTypeFile:
		blockFile := b.Block.(goodle.BlockFile)

		var fileType string

		if mime := mimetype.Lookup(blockFile.MimeType()); mime != nil {
			extension := mime.Extension()

			if extension != "" {
				fileType = strings.TrimPrefix(extension, ".")
			} else {
				fileType = mime.String()
			}
		} else {
			fileType = blockFile.MimeType()
		}

		info.WriteString(humanize.Bytes(uint64(blockFile.SizeBytes())))
		info.WriteRune(whitespace)
		info.WriteString(fileType)
	case goodle.BlockTypeAssignment:
		blockAssignment := b.Block.(goodle.BlockAssignment)

		deadline := blockAssignment.DeadlineAt()

		info.WriteString("Deadline")
		info.WriteRune(whitespace)
		info.WriteString(lo.If(deadline.IsZero(), "unknown").Else(deadline.Format(moodleFormat)))
	case goodle.BlockTypeQuiz:
		blockQuiz := b.Block.(goodle.BlockQuiz)

		opens := blockQuiz.OpensAt()
		closes := blockQuiz.ClosesAt()

		info.WriteString("Opens")
		info.WriteRune(whitespace)
		info.WriteString(lo.If(opens.IsZero(), "unknown").Else(opens.Format(moodleFormat)))
		info.WriteRune(',')
		info.WriteRune(whitespace)
		info.WriteString("closes")
		info.WriteRune(whitespace)
		info.WriteString(lo.If(closes.IsZero(), "unknown").Else(closes.Format(moodleFormat)))
	case goodle.BlockTypeLink:
		blockLink := b.Block.(goodle.BlockLink)

		info.WriteString(blockLink.Url())
	default:
		info.WriteString("No information")
	}

	return info.String() + "\n" + b.Block.Type().Name()
}

type blocksKeyMap struct {
	openBrowser,
	open,
	reverseItemsOrder key.Binding

	list list.KeyMap
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
	return [][]key.Binding{
		b.ShortHelp(),
		{b.reverseItemsOrder},
	}
}

type Blocks struct {
	section goodle.Section
	list    list.Model
	size    base.Size
	keyMap  blocksKeyMap
}

func NewBlocks(section goodle.Section) *Blocks {
	blocks := section.Blocks()

	l := util.NewList(3, blocks, func(block goodle.Block) list.Item {
		return blocksItem{block}
	})

	return &Blocks{
		list:    l,
		section: section,
		keyMap: blocksKeyMap{
			openBrowser:       util.Bind("open browser", "o"),
			open:              util.Bind("open", "enter"),
			reverseItemsOrder: util.Bind("reverse items", "r"),
			list:              l.KeyMap,
		},
	}
}

func (b *Blocks) Intermediate() bool {
	return false
}

func (b *Blocks) KeyMap() help.KeyMap {
	return b.keyMap
}

func (b *Blocks) Title() base.Title {
	return base.Title{Text: b.section.Title()}
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
	item, ok := b.list.SelectedItem().(blocksItem)
	if !ok {
		return nil
	}

	switch item.Type() {
	case goodle.BlockTypeLink:
		block := item.Block.(goodle.BlockLink)

		return openWithDefaultApp(block.Url())
	default:
		return openWithDefaultApp(item.MoodleUrl())
	}

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

func (b *Blocks) Update(_ base.Model, msg tea.Msg) (cmd tea.Cmd) {
	isFiltering := b.list.FilterState() == list.Filtering
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if isFiltering {
			goto end
		}

		switch {
		case key.Matches(msg, b.keyMap.openBrowser):
			item, ok := b.list.SelectedItem().(blocksItem)
			if !ok {
				return nil
			}

			return openWithDefaultApp(item.MoodleUrl())
		case key.Matches(msg, b.keyMap.open):
			return b.openSelected()
		case key.Matches(msg, b.keyMap.reverseItemsOrder):
			return b.list.SetItems(lo.Reverse(b.list.Items()))
		}
	}

end:
	b.list, cmd = b.list.Update(msg)
	return cmd
}

func (b *Blocks) View(base.Model) string {
	return b.list.View()
}

func (b *Blocks) Init(base.Model) tea.Cmd {
	if viper.GetBool(configKey.TUIReverseBlocks) {
		return b.list.SetItems(lo.Reverse(b.list.Items()))
	}

	return nil
}
