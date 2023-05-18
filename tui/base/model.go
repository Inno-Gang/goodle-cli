package base

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
)

type Size struct {
	Width, Height int
}

type Model interface {
	tea.Model

	Size() Size
	Context() context.Context
}
