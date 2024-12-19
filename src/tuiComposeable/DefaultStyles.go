package composeable

import lg "github.com/charmbracelet/lipgloss"

var (
	DefaultWindowStyle = lg.NewStyle().Border(lg.RoundedBorder()).Margin(1, 1)

	DefaultBadStatusStyle = lg.NewStyle().
				Foreground(lg.AdaptiveColor{Light: "#ad0000", Dark: "#820000"})
	DefaultGoodStatusStyle = lg.NewStyle().
				Foreground(lg.AdaptiveColor{Light: "#00db0f", Dark: "#005706"})
	DefaultWarnStatusStyle = lg.NewStyle().
				Foreground(lg.AdaptiveColor{Light: "#fcba03", Dark: "#fcba03"})
)

func applyFocusStyling(s lg.Style, focused bool, hovered bool) lg.Style {
	if focused {
		s = s.Bold(true).Border(lg.DoubleBorder())
	} else if hovered {
		s = s.Italic(true).Border(lg.NormalBorder())
	} else {
		s = s.Border(lg.HiddenBorder())
	}
	return s
}
