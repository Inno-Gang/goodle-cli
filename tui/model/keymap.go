package model

import (
	"github.com/Inno-Gang/goodle-cli/tui/tuiutil"
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Back, Quit key.Binding
}

func newKeyMap() *keyMap {
	return &keyMap{
		Back: tuiutil.Bind("back", "esc"),
		Quit: tuiutil.Bind("quit", "ctrl+c", "ctrl+d"),
	}
}
