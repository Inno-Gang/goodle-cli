package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit, ForceQuit, GoBack, Confirm,
	OpenBrowser, FocusNext,
	Up, Down key.Binding
}

func newKeyMap() *keyMap {
	bind := func(help string, keys ...string) key.Binding {
		return key.NewBinding(
			key.WithKeys(keys...),
			key.WithHelp(keys[0], help),
		)
	}

	return &keyMap{
		Quit:        bind("quit", "q"),
		ForceQuit:   bind("force quite", "ctrl+c", "ctrl+d"),
		GoBack:      bind("back", "esc"),
		Confirm:     bind("confirm", "enter"),
		OpenBrowser: bind("open", "o"),
		FocusNext:   bind("focus next", "tab"),
		Up:          bind("up", "up", "k"),
		Down:        bind("down", "down", "j"),
	}
}
