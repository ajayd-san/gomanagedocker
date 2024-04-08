package tui

import "github.com/charmbracelet/lipgloss"

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	fillerStyle       = lipgloss.NewStyle().Foreground(highlightColor)
	windowStyle       = lipgloss.NewStyle().
				BorderForeground(highlightColor).
				Padding(2, 0).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder()).
				UnsetBorderTop()
)
