package tuicomposeable

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	InputValue[T any] struct {
		textinput.Model
		HoveredFocused
		Err             *error
		Data            *T
		inputChangedCmd tea.Cmd

		Style  lg.Style
		Prompt string
		Parse  func(raw string) ValueErrPair[T]
		Format func(v T) string
	}
)

func NewInputValue[T any](style lg.Style) *InputValue[T] {
	return &InputValue[T]{
		Style: style,
	}
}

func (i *InputValue[T]) SetData(d *T, err *error) *InputValue[T] {
	i.Data = d
	i.Err = err
	return i
}

func (i *InputValue[T]) SetPrompt(prompt string) *InputValue[T] {
	i.Prompt = prompt
	return i
}

func (i *InputValue[T]) SetParser(
	parser func(raw string) ValueErrPair[T],
) *InputValue[T] {
	i.Parse = parser
	return i
}

func (i *InputValue[T]) SetFormater(fmter func(v T) string) *InputValue[T] {
	i.Format = fmter
	return i
}

func (i *InputValue[T]) SetInputChanedCmd(cmd tea.Cmd) *InputValue[T] {
	i.inputChangedCmd = cmd
	return i
}

func (i *InputValue[T]) SetTextInput(
	placeholder string,
	charLimit int,
) *InputValue[T] {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = charLimit
	ti.Width = charLimit
	i.Model = ti
	return i
}

func (i *InputValue[T]) FocusableElements() []Composeable {
	return []Composeable{i}
}

func (i *InputValue[T]) Focus() {
	i.HoveredFocused.Focus()
	i.Model.Focus()
}
func (i *InputValue[T]) Blur() {
	i.HoveredFocused.Blur()
	i.Model.Blur()
}
func (i *InputValue[T]) Focused() bool  { return i.HoveredFocused.Focused() }
func (i *InputValue[T]) Hovered() bool  { return i.HoveredFocused.Hovered() }
func (i *InputValue[T]) CanFocus() bool { return true }

func (i *InputValue[T]) Init() tea.Cmd {
	if i.Focused() {
		return textinput.Blink
	}
	return nil
}

func (i *InputValue[T]) Update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}

	if i.Focused() {
		switch msg := msg.(type) {
		case ValueErrPair[T]:
			*i.Data = msg.Value
			*i.Err = msg.Err
			return tea.Batch(cmds...)
		}

		var cmd tea.Cmd
		oldValue := i.Model.Value()

		i.Model, cmd = i.Model.Update(msg)
		cmds = append(cmds, cmd)
		if i.Model.Value() != oldValue {
			cmds = append(
				cmds,
				NewValueErrPairCmd[T](i.Model.Value(), i.Parse),
				i.inputChangedCmd,
			)
		}
	}

	return tea.Batch(cmds...)
}

func (i *InputValue[T]) View() string {
	var sb strings.Builder
	sb.WriteString(i.Prompt)
	sb.WriteString("\n\n")
	sb.WriteString(i.Model.View())
	sb.WriteString("\nStatus: ")
	if *(i.Err) == nil {
		sb.WriteString(DefaultGoodStatusStyle.Render("Valid"))
	} else {
		sb.WriteString(DefaultBadStatusStyle.Render((*i.Err).Error()))
	}
	sb.WriteString("\n\nValue:")
	sb.WriteString(i.Format(*i.Data))

	return applyFocusStyling(
		i.Style, i.Focused(), i.Hovered(),
	).Render(sb.String())
}

func (i *InputValue[T]) GetValue() *T {
	return i.Data
}
