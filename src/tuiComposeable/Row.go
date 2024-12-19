package tuicomposeable

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	Row struct {
		layout
	}
)

func NewRow(style lg.Style, components ...Composeable) *Row {
	return &Row{
		layout: layout{
			components: components,
			style:      style,
		},
	}
}

func (r *Row) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Width -= r.layout.style.GetHorizontalFrameSize()
		msg.Height -= r.layout.style.GetVerticalFrameSize()
		r.layout.width = msg.Width
		r.layout.height = msg.Height

		msg.Width /= len(r.layout.components)
		for _, iterComp := range r.layout.components {
			cmds = append(cmds, iterComp.Update(msg))
		}
		return tea.Batch(cmds...)
	}

	for _, iterComp := range r.components {
		cmds = append(cmds, iterComp.Update(msg))
	}
	return tea.Batch(cmds...)
}

func (r *Row) View() string {
	style := r.layout.style
	if len(r.layout.components) > 0 {
		// TODO - in the future this could be calculated in different ways
		style = style.
			Height(r.layout.height).
			Width(r.layout.width / len(r.layout.components))
	}
	return r.style.Render(lg.JoinHorizontal(lg.Top, r.layout.View(style)...))
}
