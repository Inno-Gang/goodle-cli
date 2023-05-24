package util

import (
	"github.com/Inno-Gang/goodle-cli/color"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func NewList[T any](items []T, transform func(T) list.Item) list.Model {
	var listItems = make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = transform(item)
	}

	border := lipgloss.ThickBorder()

	delegate := list.NewDefaultDelegate()

	delegate.Styles.NormalTitle.Foreground(color.HiWhite).Bold(true)
	delegate.Styles.SelectedTitle.Foreground(color.Accent).Bold(true)
	delegate.Styles.SelectedTitle.Border(border, false, false, false, true)
	delegate.Styles.SelectedDesc.
		Border(border, false, false, false, true).
		Foreground(delegate.Styles.NormalDesc.GetForeground())

	delegate.Styles.SelectedTitle.BorderLeftForeground(color.Accent)
	delegate.Styles.SelectedDesc.BorderLeftForeground(color.Accent)
	delegate.SetHeight(3)

	l := list.New(listItems, delegate, 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	l.KeyMap.CancelWhileFiltering = Bind("cancel", "esc")

	return l
}
