package tuicomposeable

import (
	"github.com/charmbracelet/bubbles/key"
)

type (
	// Defines a set of keybindings that split the UI into two modes: navigation
	// and intaractive mode.
	WindowKeyBindings struct {
		// The keyboard short cuts that will navigate to the next focusable
		// UI components.
		Next key.Binding
		// The keyboard short cuts that will navigate to the previous focusable
		// UI components.
		Prev key.Binding
		// The keyboard shortcuts to navigate to the first focusable UI
		// component.
		First key.Binding
		// The keyboard shortcuts to navigate to the last focusable UI
		// component.
		Last key.Binding
		// The keyboard short cuts that will leave navigation mode and enter
		// the interactive mode for interacting with UI components.
		Focus key.Binding
		// The keyboard short cuts that will leave interactive mode and enter
		// the navigation mode for navigating the UI components.
		Blur           key.Binding
		NavigationHelp key.Binding
		Exit           key.Binding
	}

	WindowKeyBindingModes int
)

const (
	Navigation WindowKeyBindingModes = iota
	Interactive
)

func NewWindowKeyBindings() WindowKeyBindings {
	return WindowKeyBindings{
		Next: key.NewBinding(
			key.WithKeys("tab", "shift+n", "n", "right", "l", "down", "j"),
			key.WithHelp("tab/shift+n/n/→/l/↓/j", "Next UI component"),
		),
		Prev: key.NewBinding(
			key.WithKeys("shift+tab", "shift+p", "N", "left", "h", "up", "k"),
			key.WithHelp("shift+tab/shift+n/N/←/h/↑/k", "Previous UI component"),
		),
		First: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home", "Goto first UI component"),
		),
		Last: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end", "Goto last UI component"),
		),
		Focus: key.NewBinding(
			key.WithKeys("enter", "i", "a"),
			key.WithHelp("enter/i/a", "Enter interactive mode"),
		),
		Blur: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Enter navigation mode"),
		),
		NavigationHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Extended Help"),
		),
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Exit"),
		),
	}
}

func (w *WindowKeyBindings) SetMode(mode WindowKeyBindingModes) {
	switch mode {
	case Navigation:
		w.Next.SetEnabled(true)
		w.Prev.SetEnabled(true)
		w.First.SetEnabled(true)
		w.Last.SetEnabled(true)
		w.Focus.SetEnabled(true)
		w.Blur.SetEnabled(false)
		w.NavigationHelp.SetEnabled(true)
		w.Exit.SetEnabled(true)
	case Interactive:
		w.Next.SetEnabled(false)
		w.Prev.SetEnabled(false)
		w.First.SetEnabled(false)
		w.Last.SetEnabled(false)
		w.Focus.SetEnabled(false)
		w.Blur.SetEnabled(true)
		w.NavigationHelp.SetEnabled(true)
		w.Exit.SetEnabled(true)
	}
}

func (w WindowKeyBindings) ShortHelp() []key.Binding {
	return []key.Binding{
		w.Next, w.Focus, w.Blur, w.NavigationHelp, w.Exit,
	}
}

func (w WindowKeyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		[]key.Binding{w.Next, w.Prev},
		[]key.Binding{w.First, w.Last},
		[]key.Binding{w.Focus, w.Blur},
		[]key.Binding{w.NavigationHelp, w.Exit},
	}
}

func (w WindowKeyBindings) ShortHelpHeight() int {
	return 1
}

func (w WindowKeyBindings) FullHelpHeight() int {
	rv := 0
	for _, col := range w.FullHelp() {
		rv = max(rv, len(col))
	}
	return rv
}
