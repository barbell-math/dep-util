package tuicomposeable

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	Button struct {
		HoveredFocused
		name     string
		clickCmd tea.Cmd

		Style    lg.Style
		canClick bool
	}
)

func NewButton(style lg.Style, name string, clickCmd tea.Cmd) *Button {
	return &Button{
		name:     name,
		clickCmd: clickCmd,
		Style:    style,
	}
}

func (b *Button) FocusableElements() []Composeable {
	return []Composeable{b}
}

func (b *Button) CanFocus() bool { return true }

func (b *Button) Init() tea.Cmd {
	return nil
}

func (b *Button) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case FocusEvent:
		if b.Focused() {
			return tea.Batch(NewRemoveFocusCmd(), b.clickCmd)
		}
	}
	return nil
}

func (b *Button) View() string {
	if b.Hovered() {
		return applyFocusStyling(
			b.Style, b.Focused(), b.Hovered(),
		).Render(b.name + "↵")
	}
	return applyFocusStyling(b.Style, b.Focused(), b.Hovered()).Render(b.name)
}
