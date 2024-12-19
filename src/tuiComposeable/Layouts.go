package composeable

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	layout struct {
		components []Composeable
		style      lg.Style

		height int
		width  int
	}
)

func (l *layout) FocusableElements() []Composeable {
	rv := []Composeable{}
	for _, iterC := range l.components {
		rv = append(rv, iterC.FocusableElements()...)
	}
	return rv
}

func (l *layout) Focus() {
	// intentional noop - layout elements cannot be focused on
}
func (l *layout) Hover() {
	// intentional noop - layout elements cannot be hovered on
}
func (l *layout) UnHover() {
	// intentional noop - layout elements cannot be hovered on
}
func (l *layout) Blur() {
	// intentional noop - layout elements cannot be focused on
}
func (l *layout) Focused() bool  { return false }
func (l *layout) Hovered() bool  { return false }
func (l *layout) CanFocus() bool { return false }

func (l *layout) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	for _, v := range l.components {
		cmds = append(cmds, v.Init())
	}
	return tea.Batch(cmds...)
}

func (l layout) Update(msg tea.Msg) (layout, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Width -= l.style.GetHorizontalFrameSize()
		msg.Height -= l.style.GetVerticalFrameSize()
		l.width = msg.Width
		l.height = msg.Height

		for _, iterComp := range l.components {
			cmds = append(cmds, iterComp.Update(msg))
		}
		return l, tea.Batch(cmds...)
	}

	for _, iterComp := range l.components {
		cmds = append(cmds, iterComp.Update(msg))
	}
	return l, tea.Batch(cmds...)
}

func (l *layout) View(s lg.Style) []string {
	renderedComponents := make([]string, len(l.components))
	if len(l.components) > 0 {

		for i, iterComp := range l.components {
			renderedComponents[i] = s.Render(iterComp.View())
		}
	}
	return renderedComponents
}
