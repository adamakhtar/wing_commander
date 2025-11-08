package styles

import "github.com/charmbracelet/lipgloss"

var (
	White = lipgloss.Color("#FFFFFF")
	Black = lipgloss.Color("#000000")

	Gray50 = lipgloss.Color("#f9fafb")
	Gray100 = lipgloss.Color("#f3f4f6")
	Gray200 = lipgloss.Color("#e5e7eb")
	Gray300 = lipgloss.Color("#d1d5db")
	Gray400 = lipgloss.Color("#9ca3af")
	Gray500 = lipgloss.Color("#6b7280")
	Gray600 = lipgloss.Color("#4b5563")
	Gray700 = lipgloss.Color("#374151")
	Gray800 = lipgloss.Color("#1f2937")
	Gray900 = lipgloss.Color("#111827")
	Gray950 = lipgloss.Color("#030712")

	Yellow50 = lipgloss.Color("#fefce8")
	Yellow100 = lipgloss.Color("#fef9c3")
	Yellow200 = lipgloss.Color("#fef08a")
	Yellow300 = lipgloss.Color("#fde047")
	Yellow400 = lipgloss.Color("#facc15")
	Yellow500 = lipgloss.Color("#eab308")
	Yellow600 = lipgloss.Color("#ca8a04")
	Yellow700 = lipgloss.Color("#a16207")
	Yellow800 = lipgloss.Color("#854d0e")
	Yellow900 = lipgloss.Color("#713f12")
	Yellow950 = lipgloss.Color("#422006")

	Amber50 = lipgloss.Color("#fffbeb")
	Amber100 = lipgloss.Color("#fef3c7")
	Amber200 = lipgloss.Color("#fde68a")
	Amber300 = lipgloss.Color("#fcd34d")
	Amber400 = lipgloss.Color("#fbbf24")
	Amber500 = lipgloss.Color("#f59e0b")
	Amber600 = lipgloss.Color("#d97706")
	Amber700 = lipgloss.Color("#b45309")
	Amber800 = lipgloss.Color("#92400e")
	Amber900 = lipgloss.Color("#78350f")
	Amber950 = lipgloss.Color("#451a03")

	Orange50 = lipgloss.Color("#fff7ed")
	Orange100 = lipgloss.Color("#ffedd5")
	Orange200 = lipgloss.Color("#fed7aa")
	Orange300 = lipgloss.Color("#fdba74")
	Orange400 = lipgloss.Color("#fb923c")
	Orange500 = lipgloss.Color("#f97316")
	Orange600 = lipgloss.Color("#ea580c")
	Orange700 = lipgloss.Color("#c2410c")
	Orange800 = lipgloss.Color("#9a3412")
	Orange900 = lipgloss.Color("#7c2d12")
	Orange950 = lipgloss.Color("#431407")

	Red50 = lipgloss.Color("#fef2f2")
	Red100 = lipgloss.Color("#fee2e2")
	Red200 = lipgloss.Color("#fecaca")
	Red300 = lipgloss.Color("#fca5a5")
	Red400 = lipgloss.Color("#f87171")
	Red500 = lipgloss.Color("#ef4444")
	Red600 = lipgloss.Color("#dc2626")
	Red700 = lipgloss.Color("#b91c1c")
	Red800 = lipgloss.Color("#991b1b")
	Red900 = lipgloss.Color("#7f1d1d")
	Red950 = lipgloss.Color("#450a0a")

	Pink50 = lipgloss.Color("#fdf2f8")
	Pink100 = lipgloss.Color("#fce7f3")
	Pink200 = lipgloss.Color("#fbcfe8")
	Pink300 = lipgloss.Color("#f9a8d4")
	Pink400 = lipgloss.Color("#f472b6")
	Pink500 = lipgloss.Color("#ec4899")
	Pink600 = lipgloss.Color("#db2777")
	Pink700 = lipgloss.Color("#be185d")
	Pink800 = lipgloss.Color("#9d174d")
	Pink900 = lipgloss.Color("#831843")
	Pink950 = lipgloss.Color("#500724")

	Purple50 = lipgloss.Color("#faf5ff")
	Purple100 = lipgloss.Color("#f3e8ff")
	Purple200 = lipgloss.Color("#e9d5ff")
	Purple300 = lipgloss.Color("#d8b4fe")
	Purple400 = lipgloss.Color("#c084fc")
	Purple500 = lipgloss.Color("#a855f7")
	Purple600 = lipgloss.Color("#9333ea")
	Purple700 = lipgloss.Color("#7c3aed")
	Purple800 = lipgloss.Color("#6b21a8")
	Purple900 = lipgloss.Color("#581c87")
	Purple950 = lipgloss.Color("#3b0764")

	Indigo50 = lipgloss.Color("#eef2ff")
	Indigo100 = lipgloss.Color("#e0e7ff")
	Indigo200 = lipgloss.Color("#c7d2fe")
	Indigo300 = lipgloss.Color("#a5b4fc")
	Indigo400 = lipgloss.Color("#818cf8")
	Indigo500 = lipgloss.Color("#6366f1")
	Indigo600 = lipgloss.Color("#4f46e5")
	Indigo700 = lipgloss.Color("#4338ca")
	Indigo800 = lipgloss.Color("#3730a3")
	Indigo900 = lipgloss.Color("#312e81")
	Indigo950 = lipgloss.Color("#1e1b4b")
)

type Theme struct {
	BodyText lipgloss.Color
	BodyTextLight lipgloss.Color
	PrimaryText lipgloss.Color
	SecondaryText lipgloss.Color
	ErrorText lipgloss.Color
	ErrorBackground lipgloss.Color
	PrimaryBorderColor lipgloss.Color
	PrimaryBorderActiveColor lipgloss.Color
	TableBorderColor lipgloss.Color
	TableRowTextColor lipgloss.Color
	TableHeaderTextColor lipgloss.Color
	TableSelectedBackground lipgloss.Color
}

var DefaultTheme Theme = Theme {
	BodyText: lipgloss.Color(White),
	BodyTextLight: lipgloss.Color(Gray400),
	PrimaryText: lipgloss.Color(Pink500),
	SecondaryText: lipgloss.Color(Pink600),
	ErrorText: lipgloss.Color(Red400),
	ErrorBackground: lipgloss.Color(Red950),
	PrimaryBorderColor: lipgloss.Color(Purple700),
	PrimaryBorderActiveColor: lipgloss.Color(Purple500),
	TableBorderColor: lipgloss.Color(Purple900),
	TableHeaderTextColor: lipgloss.Color(Pink500),
	TableRowTextColor: lipgloss.Color(Gray100),
	TableSelectedBackground: lipgloss.Color(Purple700),
}

type Styles struct {
	BodyText lipgloss.Style
	BodyTextLight lipgloss.Style
	HeadingTextStyle lipgloss.Style
	Border lipgloss.Style
	BorderMuted lipgloss.Style
	BorderActive lipgloss.Style
	TestDefinitionErrorBadge lipgloss.Style
	ProductionCodeErrorBadge lipgloss.Style
	AssertionErrorBadge lipgloss.Style
	Preview struct {
		AlertStyle lipgloss.Style
	}
	ResultsSection struct {
		TableBaseStyle lipgloss.Style
		TableBorderColor lipgloss.Color
		TableHeaderTextColor lipgloss.Color
		TableRowTextColor lipgloss.Color
		TableHighlight lipgloss.Style
	}
	PreviewSection struct {
		BacktracePath lipgloss.Style
		CodeLine lipgloss.Style
		HighlightedCodeLine lipgloss.Style
		SnippetBorder lipgloss.Style
	}
	// FaintTextStyle lipgloss.Style
	// Results struct {

	// }
}


func BuildStyles(theme Theme) Styles {
	s := Styles{
		BodyText: lipgloss.NewStyle().Foreground(theme.BodyText),
		BodyTextLight: lipgloss.NewStyle().Foreground(theme.BodyTextLight),
		HeadingTextStyle: lipgloss.NewStyle().Bold(true).Foreground(theme.PrimaryText),
		Border: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
		BorderMuted: lipgloss.NewStyle().BorderForeground(theme.PrimaryBorderColor),
		BorderActive: lipgloss.NewStyle().BorderForeground(theme.PrimaryBorderActiveColor),
		TestDefinitionErrorBadge: lipgloss.NewStyle().Background(Yellow600).Foreground(White).Align(lipgloss.Center),
		ProductionCodeErrorBadge: lipgloss.NewStyle().Background(Black).Foreground(Red500).Align(lipgloss.Center),
		AssertionErrorBadge: lipgloss.NewStyle().Background(Orange600).Foreground(White).Align(lipgloss.Center),
	}

	s.Preview.AlertStyle = lipgloss.NewStyle().Background(theme.ErrorBackground).Foreground(theme.ErrorText)

	s.ResultsSection.TableBaseStyle = lipgloss.NewStyle().
									Align(lipgloss.Left).
									BorderForeground(theme.TableBorderColor)
	s.ResultsSection.TableHeaderTextColor = theme.TableHeaderTextColor
	s.ResultsSection.TableRowTextColor = theme.TableRowTextColor
	s.ResultsSection.TableHighlight = lipgloss.NewStyle().Background(theme.TableSelectedBackground).Foreground(theme.TableRowTextColor)

	s.PreviewSection.BacktracePath = lipgloss.NewStyle().Underline(true).Foreground(theme.BodyText).Bold(true)
	s.PreviewSection.CodeLine = lipgloss.NewStyle().Foreground(Gray400)
	s.PreviewSection.HighlightedCodeLine = lipgloss.NewStyle().Background(Pink800).Foreground(White)
	s.PreviewSection.SnippetBorder = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(theme.PrimaryBorderColor)

	return s
}