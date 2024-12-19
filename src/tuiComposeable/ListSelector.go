package tuicomposeable

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	Selector[T any, U interface {
		*T
		list.DefaultItem
	}] struct {
		HoveredFocused
		list.Model
		selectedItem **T
		id           int
		customUpdate func(tea.Msg) tea.Cmd
		selectionCmd tea.Cmd

		Style           lg.Style
		GoodStatusStyle lg.Style
		BadStatusStyle  lg.Style
	}
)

func NewSelector[T any, U interface {
	*T
	list.DefaultItem
}](
	style lg.Style,
	id int,
	title string,
	initialData []T,
) *Selector[T, U] {
	items := make([]list.Item, len(initialData))
	for i, v := range initialData {
		var tmp U
		tmp = &v
		items[i] = tmp
	}

	rv := Selector[T, U]{
		id:              id,
		Model:           list.New(items, list.NewDefaultDelegate(), 0, 0),
		Style:           style.Border(lg.HiddenBorder()),
		GoodStatusStyle: DefaultGoodStatusStyle,
		BadStatusStyle:  DefaultBadStatusStyle,
	}
	rv.Model.Title = title
	rv.Model.KeyMap.Quit = key.NewBinding()
	return &rv
}

func (s *Selector[T, U]) SetCustomUpdate(
	customUpdate func(tea.Msg) tea.Cmd,
) *Selector[T, U] {
	s.customUpdate = customUpdate
	return s
}

func (s *Selector[T, U]) SetSelectionCmd(cmd tea.Cmd) *Selector[T, U] {
	s.selectionCmd = cmd
	return s
}

func (s *Selector[T, U]) SetSelectedItem(v **T) *Selector[T, U] {
	s.selectedItem = v
	return s
}

func (s *Selector[T, U]) FocusableElements() []Composeable {
	return []Composeable{s}
}

func (s *Selector[T, U]) CanFocus() bool { return true }

func (s *Selector[T, U]) Init() tea.Cmd {
	s.Model.NewStatusMessage(s.BadStatusStyle.Render(
		"\nYou must select an item.",
	))
	return nil
}

func (s *Selector[T, U]) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}

	if msg, ok := msg.(tea.KeyMsg); ok && s.Focused() && s.Model.FilterState() != list.Filtering {
		switch msg.String() {
		case "enter":
			s.setChosenItem()
			if s.selectionCmd != nil {
				cmds = append(cmds, s.selectionCmd)
			}
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Width -= s.Style.GetHorizontalFrameSize()
		msg.Height -= s.Style.GetVerticalFrameSize()
		s.Model.SetSize(msg.Width, msg.Height-1)
		return nil
	case UpdateSelectorItems[T]:
		items := make([]list.Item, len(msg.Data))
		for i, v := range msg.Data {
			var tmp U
			tmp = &v
			items[i] = tmp
		}
		return s.Model.SetItems(items)
	}

	if s.customUpdate != nil {
		cmds = append(cmds, s.customUpdate(msg))
	}

	if s.Focused() {
		chosenItem := U(*s.selectedItem)
		if chosenItem != nil {
			s.Model.NewStatusMessage(s.GoodStatusStyle.Render(
				fmt.Sprintf("\nSelected: %s", chosenItem.Title()),
			))
		} else {
			s.Model.NewStatusMessage(s.BadStatusStyle.Render(
				"\nYou must select an item.",
			))
		}
		var cmd tea.Cmd
		s.Model, cmd = s.Model.Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func (s *Selector[T, U]) View() string {
	return applyFocusStyling(
		s.Style, s.Focused(), s.Hovered(),
	).Render(s.Model.View())
}

func (s *Selector[T, U]) setChosenItem() {
	if i, ok := s.Model.SelectedItem().(U); ok {
		*s.selectedItem = i
	}
}
