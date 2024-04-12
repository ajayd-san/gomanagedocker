package tui

import "github.com/charmbracelet/lipgloss"

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(0, 1, 0, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	fillerStyle       = lipgloss.NewStyle().Foreground(highlightColor)
	windowStyle       = lipgloss.NewStyle().
				BorderForeground(highlightColor).
				Border(lipgloss.NormalBorder()).
				UnsetBorderTop()

	listDocStyle  = lipgloss.NewStyle().Margin(1, 2)
	moreInfoStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69")).
			Width(90).
			Height(35).
			Margin(0, 0, 0, 30)
)
