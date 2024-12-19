package tuicomposeable

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	Tabs struct {
		layout
		tabButtons *Row
		openTab    int
	}

	tabButton struct {
		HoveredFocused
		name string
		open bool
	}
)

func (t *tabButton) FocusableElements() []Composeable {
	return []Composeable{t}
}

func (t *tabButton) CanFocus() bool { return true }

func (t *tabButton) Init() tea.Cmd {
	return nil
}

func (t *tabButton) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (t *tabButton) View() string {
	if !t.Focused() && !t.Hovered() && !t.open {
		return lg.NewStyle().Border(lg.Border{
			Bottom:      "_",
			BottomLeft:  "_",
			BottomRight: "_",
		}).Render(t.name)
	} else if t.Focused() || t.Hovered() {
		return applyFocusStyling(
			lg.NewStyle(), t.Focused(), t.Hovered(),
		).Render(t.name + "↵")
	} else { // t.open
		return applyFocusStyling(
			lg.NewStyle(), true, t.Hovered(),
		).Render(t.name)
	}
}

func NewTabs(style lg.Style, names []string, components ...Composeable) *Tabs {
	buttons := make([]Composeable, len(names))
	for i, name := range names {
		buttons[i] = &tabButton{name: name, open: i == 0}
	}
	return &Tabs{
		layout: layout{
			components: components,
			style:      lg.NewStyle(),
		},
		tabButtons: NewRow(lg.NewStyle(), buttons...),
	}
}

func (t *Tabs) FocusableElements() []Composeable {
	rv := []Composeable{}
	rv = append(rv, t.tabButtons.FocusableElements()...)
	rv = append(rv, t.layout.components[t.openTab].FocusableElements()...)
	return rv
}

func (t *Tabs) Focus() {
	// intentional noop - layout elements cannot be focused on
}
func (t *Tabs) Hover() {
	// intentional noop - layout elements cannot be hovered on
}
func (t *Tabs) UnHover() {
	// intentional noop - layout elements cannot be hovered on
}
func (t *Tabs) Blur() {
	// intentional noop - layout elements cannot be focused on
}
func (t *Tabs) Focused() bool  { return false }
func (t *Tabs) Hovered() bool  { return false }
func (t *Tabs) CanFocus() bool { return false }

func (t *Tabs) Init() tea.Cmd {
	return t.layout.Init()
}

func (t *Tabs) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Width -= t.style.GetHorizontalFrameSize()
		msg.Height -= t.style.GetVerticalFrameSize() + 3 // height of tab button

		t.layout, cmd = t.layout.Update(msg)
		cmds = append(cmds, cmd)
		return tea.Batch(cmds...)
	case FocusEvent:
		// Check if we are focused on a tab. If we are then update the currently
		// open tab and send a cmd to remove the focus from the tab button.
		focusedOnTab := false
		for i, tab := range t.tabButtons.layout.components {
			if tab.Focused() {
				t.tabButtons.layout.components[t.openTab].(*tabButton).open = false
				t.openTab = i
				t.tabButtons.layout.components[t.openTab].(*tabButton).open = true
				focusedOnTab = true
				break
			}
		}
		if focusedOnTab {
			cmds = append(
				cmds, NewRemoveFocusCmd(), NewRecomputeFocusableElements(),
			)
			return tea.Batch(cmds...)
		}
	}

	t.layout, cmd = t.layout.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (t *Tabs) View() string {
	return lg.JoinVertical(
		lg.Top,
		t.layout.style.Render(t.tabButtons.View()),
		t.layout.components[t.openTab].View(),
	)
}
