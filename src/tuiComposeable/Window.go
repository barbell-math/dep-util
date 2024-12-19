package tuicomposeable

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	Window struct {
		Root  Composeable
		Style lg.Style

		hoveredComponent    int
		focusableComponents []Composeable

		keyBindings WindowKeyBindings
		help        help.Model
	}
)

func NewWindow(style lg.Style, root Composeable) Window {
	rv := Window{
		Root:                root,
		Style:               style,
		focusableComponents: root.FocusableElements(),
		keyBindings:         NewWindowKeyBindings(),
		help:                help.New(),
	}
	rv.keyBindings.SetMode(Navigation)
	if len(rv.focusableComponents) > 0 {
		rv.focusableComponents[rv.hoveredComponent].Hover()
	}
	return rv
}

func (w Window) Init() tea.Cmd {
	return w.Root.Init()
}

func (w Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, w.keyBindings.Exit):
			return w, tea.Quit
		case key.Matches(msg, w.keyBindings.Next):
			if len(w.focusableComponents) == 0 {
				return w, nil
			}
			if !w.focusableComponents[w.hoveredComponent].Focused() {
				w.incHoveredComponent(1)
				w.keyBindings.SetMode(Navigation)
				return w, nil
			}
		case key.Matches(msg, w.keyBindings.Prev):
			if len(w.focusableComponents) == 0 {
				return w, nil
			}
			if !w.focusableComponents[w.hoveredComponent].Focused() {
				w.incHoveredComponent(-1)
				w.keyBindings.SetMode(Navigation)
				return w, nil
			}
		case key.Matches(msg, w.keyBindings.First):
			if len(w.focusableComponents) == 0 {
				return w, nil
			}
			if !w.focusableComponents[w.hoveredComponent].Focused() {
				w.incHoveredComponent(-w.hoveredComponent)
				w.keyBindings.SetMode(Navigation)
				return w, nil
			}
		case key.Matches(msg, w.keyBindings.Last):
			if len(w.focusableComponents) == 0 {
				return w, nil
			}
			if !w.focusableComponents[w.hoveredComponent].Focused() {
				w.incHoveredComponent(
					len(w.focusableComponents) - w.hoveredComponent - 1,
				)
				w.keyBindings.SetMode(Navigation)
				return w, nil
			}
		case key.Matches(msg, w.keyBindings.Focus):
			if len(w.focusableComponents) == 0 {
				return w, nil
			}
			if !w.focusableComponents[w.hoveredComponent].Focused() {
				w.focusableComponents[w.hoveredComponent].Focus()
				w.keyBindings.SetMode(Interactive)
				return w, NewFocusEventCmd()
			}
		case key.Matches(msg, w.keyBindings.Blur):
			if len(w.focusableComponents) == 0 {
				return w, nil
			}
			if w.focusableComponents[w.hoveredComponent].Focused() {
				w.focusableComponents[w.hoveredComponent].Blur()
				w.keyBindings.SetMode(Navigation)
				return w, NewBlurEventCmd()
			}
		case key.Matches(msg, w.keyBindings.NavigationHelp):
			if !w.focusableComponents[w.hoveredComponent].Focused() {
				w.help.ShowAll = !w.help.ShowAll
				cmds = append(cmds, tea.WindowSize())
			}
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Width -= w.Style.GetHorizontalFrameSize()
		msg.Height -= w.Style.GetVerticalFrameSize()
		if w.help.ShowAll {
			msg.Height -= w.keyBindings.FullHelpHeight()
		} else {
			msg.Height -= w.keyBindings.ShortHelpHeight()
		}

		cmds = append(cmds, w.Root.Update(msg))
		return w, tea.Batch(cmds...)
	case RemoveFocus:
		if len(w.focusableComponents) == 0 {
			return w, nil
		}
		if w.focusableComponents[w.hoveredComponent].Focused() {
			w.focusableComponents[w.hoveredComponent].Blur()
			w.keyBindings.SetMode(Navigation)
			return w, nil
		}
	case RecomputeFocusableElements:
		w.focusableComponents = w.Root.FocusableElements()
		if w.hoveredComponent >= len(w.focusableComponents) {
			w.hoveredComponent = len(w.focusableComponents) - 1
			w.focusableComponents[w.hoveredComponent].Hover()
		}
	}

	cmds = append(cmds, w.Root.Update(msg))
	return w, tea.Batch(cmds...)
}

func (w *Window) incHoveredComponent(inc int) {
	w.focusableComponents[w.hoveredComponent].UnHover()
	w.hoveredComponent += inc
	if w.hoveredComponent < 0 {
		w.hoveredComponent += len(w.focusableComponents)
	}
	w.hoveredComponent %= len(w.focusableComponents)
	w.focusableComponents[w.hoveredComponent].Hover()
}

func (w Window) View() string {
	return lg.JoinVertical(
		lg.Top,
		w.Style.Render(w.Root.View()),
		w.help.View(w.keyBindings),
	)
}

func (w Window) Main() error {
	p := tea.NewProgram(w, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
