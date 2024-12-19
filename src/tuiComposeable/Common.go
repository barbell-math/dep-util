package tuicomposeable

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	Status int

	HoveredFocused struct {
		value uint
	}

	Composeable interface {
		Init() tea.Cmd
		Update(msg tea.Msg) tea.Cmd
		View() string

		Focus()
		Hover()
		UnHover()
		Blur()
		Focused() bool
		Hovered() bool
		CanFocus() bool
		FocusableElements() []Composeable
	}
)

const (
	Waiting Status = iota
	Ok
	Error
)

func (h *HoveredFocused) Focus() {
	h.value |= uint(1)
}
func (h *HoveredFocused) Blur() {
	h.value &= ^uint(1)
}
func (h *HoveredFocused) Hover() {
	h.value |= uint(2)
}
func (h *HoveredFocused) UnHover() {
	h.value &= ^uint(2)
}
func (h *HoveredFocused) Focused() bool { return h.value&uint(1) != 0 }
func (h *HoveredFocused) Hovered() bool { return h.value&uint(2) != 0 }

func (s Status) Symbol() rune {
	switch s {
	case Ok:
		return '✓'
	case Error:
		return '❌'
	default:
		return ' '
	}
}
