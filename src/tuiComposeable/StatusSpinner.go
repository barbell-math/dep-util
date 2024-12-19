package composeable

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	StatusSpinner struct {
		Status
		spinner.Model

		Err        error
		GoodMsg    string
		WaitingMsg string
		ErrMsg     string

		Style              lg.Style
		GoodStatusStyle    lg.Style
		BadStatusStyle     lg.Style
		WaitingStatusStyle lg.Style

		id int
	}
)

func NewStatusSpinner(style lg.Style, id int) *StatusSpinner {
	return &StatusSpinner{
		id:                 id,
		Model:              spinner.New(),
		Status:             Ok,
		Style:              style,
		GoodStatusStyle:    DefaultGoodStatusStyle,
		BadStatusStyle:     DefaultBadStatusStyle,
		WaitingStatusStyle: DefaultWarnStatusStyle,
	}
}

func (s *StatusSpinner) SetGoodMsg(m string) *StatusSpinner {
	s.GoodMsg = m
	return s
}
func (s *StatusSpinner) SetErrMsg(m string) *StatusSpinner {
	s.ErrMsg = m
	return s
}
func (s *StatusSpinner) SetWaitingMsg(m string) *StatusSpinner {
	s.WaitingMsg = m
	return s
}

func (s *StatusSpinner) FocusableElements() []Composeable {
	return []Composeable{}
}

func (s *StatusSpinner) Focus() {
	// intentional noop - status spinner cannot be focused
}
func (s *StatusSpinner) Blur() {
	// intentional noop - status spinner cannot be focused
}
func (s *StatusSpinner) Hover() {
	// intentional noop - status spinner cannot be hovered
}
func (s *StatusSpinner) UnHover() {
	// intentional noop - status spinner cannot be hovered
}
func (s *StatusSpinner) Focused() bool  { return false }
func (s *StatusSpinner) Hovered() bool  { return false }
func (s *StatusSpinner) CanFocus() bool { return false }

func (s *StatusSpinner) Init() tea.Cmd {
	return nil
}

func (s *StatusSpinner) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case UpdateStatusSpinner:
		if msg.StatusSpinnerID == s.id {
			s.Spinner = spinner.Pulse
			s.Status = msg.Status
			s.Err = msg.Err
			return s.Model.Tick
		}
	}
	if s.Status == Waiting {
		var cmd tea.Cmd
		s.Model, cmd = s.Model.Update(msg)
		return cmd
	}
	return nil
}

func (s *StatusSpinner) View() string {
	var style lg.Style
	var message string

	switch s.Status {
	case Ok:
		style = s.GoodStatusStyle
		message = string(s.Symbol()) + " " + s.GoodMsg
	case Error:
		style = s.BadStatusStyle
		message = string(s.Symbol()) + " " + s.ErrMsg + ": " + s.Err.Error()
	case Waiting:
		style = s.WaitingStatusStyle
		message = s.Model.View() + " " + s.WaitingMsg
	}

	return style.Render(message)
}
