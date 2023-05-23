package util

import (
	"github.com/Inno-Gang/goodle-cli/color"
	"github.com/charmbracelet/bubbles/list"
)

func NewList[T any](items []T, transform func(T) list.Item) list.Model {
	var listItems = make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = transform(item)
	}

	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle.Foreground(color.Accent)
	delegate.Styles.SelectedDesc.Foreground(color.AccentDarken)

	delegate.Styles.SelectedTitle.BorderLeftForeground(color.Accent)
	delegate.Styles.SelectedDesc.BorderLeftForeground(color.Accent)

	l := list.New(listItems, delegate, 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	l.KeyMap.CancelWhileFiltering = Bind("cancel", "esc")

	return l
}
