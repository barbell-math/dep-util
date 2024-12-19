package tuicomposeable

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	FocusEvent  struct{}
	BlurEvent   struct{}
	RemoveFocus struct{}

	RecomputeFocusableElements struct{}

	ValueErrPair[T any] struct {
		Value T
		Err   error
	}

	UpdateStatusSpinner struct {
		StatusSpinnerID int
		Status
		Err error
	}

	UpdateTableData[T any] struct {
		TableID int
		Data    [][]T
	}

	UpdateSelectorItems[T any] struct {
		SelectorID int
		Data       []T
	}
)

func NewFocusEventCmd() tea.Cmd {
	return func() tea.Msg {
		return FocusEvent{}
	}
}

func NewBlurEventCmd() tea.Cmd {
	return func() tea.Msg {
		return BlurEvent{}
	}
}

func NewRemoveFocusCmd() tea.Cmd {
	return func() tea.Msg {
		return RemoveFocus{}
	}
}

func NewRecomputeFocusableElements() tea.Cmd {
	return func() tea.Msg {
		return RecomputeFocusableElements{}
	}
}

func NewValueErrPairCmd[T any](
	raw string,
	op func(raw string) ValueErrPair[T],
) tea.Cmd {
	return func() tea.Msg {
		return op(raw)
	}
}

func (v ValueErrPair[T]) Error() string {
	return v.Err.Error()
}

func NewUpdateStatusSpinnerCmd(id int, status Status, err error) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusSpinner{
			StatusSpinnerID: id,
			Status:          status,
			Err:             err,
		}
	}
}

func (s UpdateStatusSpinner) Error() string {
	return s.Err.Error()
}

func NewUpdateTableDataCmd[T any](id int, data [][]T) tea.Cmd {
	return func() tea.Msg {
		return UpdateTableData[T]{TableID: id, Data: data}
	}
}

func NewUpdateSelectorItems[T any](id int, data []T) tea.Cmd {
	return func() tea.Msg {
		return UpdateSelectorItems[T]{SelectorID: id, Data: data}
	}
}
