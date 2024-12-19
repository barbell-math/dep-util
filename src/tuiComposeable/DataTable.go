package composeable

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	DataTable[T any] struct {
		HoveredFocused
		table.Model

		Style        lg.Style
		id           int
		maxWidths    []int
		customUpdate func(tea.Msg) tea.Cmd
	}
)

func NewDataTable[T any](style lg.Style, id int) *DataTable[T] {
	rv := DataTable[T]{
		Style: style,
		id:    id,
	}
	return &rv
}

func (d *DataTable[T]) SetHeadersAndData(
	headers []string,
	data [][]T,
	maxWidths []int,
) *DataTable[T] {
	rows := d.getRowsFromData(data)

	columns := make([]table.Column, len(headers))
	for i, h := range headers {
		columns[i].Title = h
	}
	d.calculateColumnWidths(columns, rows, maxWidths)

	d.Model = table.New(table.WithColumns(columns), table.WithRows(rows))
	d.maxWidths = maxWidths
	return d
}

func (d *DataTable[T]) SetCustomUpdate(
	customUpdate func(tea.Msg) tea.Cmd,
) *DataTable[T] {
	d.customUpdate = customUpdate
	return d
}

func (d *DataTable[T]) getRowsFromData(data [][]T) []table.Row {
	rv := make([]table.Row, len(data))
	for i, r := range data {
		rv[i] = make([]string, len(r))
		for j, cell := range r {
			rv[i][j] = fmt.Sprint(cell)
		}
	}
	return rv
}

func (d *DataTable[T]) calculateColumnWidths(
	cols []table.Column,
	rows []table.Row,
	maxWidths []int,
) {
	for i, h := range cols {
		cols[i].Width = min(maxWidths[i], len(h.Title))
	}
	for _, r := range rows {
		for j, cell := range r {
			cols[j].Width = min(maxWidths[j], max(cols[j].Width, len(cell)))
		}
	}
}

func (d *DataTable[T]) FocusableElements() []Composeable {
	return []Composeable{d}
}

func (d *DataTable[T]) Focus() {
	d.HoveredFocused.Focus()
	d.Model.Focus()
}
func (d *DataTable[T]) Blur() {
	d.HoveredFocused.Blur()
	d.Model.Blur()
}
func (d *DataTable[T]) Focused() bool  { return d.HoveredFocused.Focused() }
func (d *DataTable[T]) CanFocus() bool { return true }

func (d *DataTable[T]) Init() tea.Cmd {
	return nil
}

func (d *DataTable[T]) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Width -= d.Style.GetHorizontalFrameSize() + len(d.Columns()) + 3
		msg.Height -= d.Style.GetVerticalFrameSize() + 3
		d.Model.SetWidth(msg.Width)
		d.Model.SetHeight(msg.Height)
		return nil
	case UpdateTableData[T]:
		if msg.TableID == d.id {
			rows := d.getRowsFromData(msg.Data)
			cols := d.Model.Columns()
			d.calculateColumnWidths(cols, rows, d.maxWidths)

			d.Model.SetRows(rows)
			d.Model.SetColumns(cols)
			d.Model.UpdateViewport()
		}
	}

	cmds = append(cmds, d.customUpdate(msg))

	if d.Focused() {
		var cmd tea.Cmd
		d.Model, cmd = d.Model.Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func (d *DataTable[T]) View() string {
	return applyFocusStyling(
		d.Style, d.Focused(), d.Hovered(),
	).Render(d.Model.View())
}
