package model

import (
	"github.com/Inno-Gang/goodle-cli/tui/util"
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Back, Quit, Help key.Binding
}

func newKeyMap() *keyMap {
	return &keyMap{
		Back: util.Bind("back", "esc"),
		Quit: util.Bind("quit", "ctrl+c", "ctrl+d"),
		Help: util.Bind("help", "?"),
	}
}
