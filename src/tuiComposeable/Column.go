package composeable

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	Column struct {
		layout
	}
)

func NewColumn(style lg.Style, components ...Composeable) *Column {
	return &Column{
		layout: layout{
			components: components,
			style:      style,
		},
	}
}

func (c *Column) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Width -= c.layout.style.GetHorizontalFrameSize()
		msg.Height -= c.layout.style.GetVerticalFrameSize()
		c.layout.width = msg.Width
		c.layout.height = msg.Height

		msg.Height /= len(c.layout.components)
		for _, iterComp := range c.layout.components {
			cmds = append(cmds, iterComp.Update(msg))
		}
		return tea.Batch(cmds...)
	}

	for _, iterComp := range c.components {
		cmds = append(cmds, iterComp.Update(msg))
	}
	return tea.Batch(cmds...)
}

func (c *Column) View() string {
	style := c.layout.style
	if len(c.layout.components) > 0 {
		// TODO - in the future this could be calculated in different ways
		style = style.
			Height(c.layout.height / len(c.layout.components)).
			Width(c.layout.width)
	}
	return c.style.Render(lg.JoinVertical(lg.Top, c.layout.View(style)...))
}
