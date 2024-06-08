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

	listDocStyle  = lipgloss.NewStyle().Margin(1, 5, 0, 1)
	listContainer = lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Width(60)
	moreInfoStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69")).
			Width(90).
			Height(25).MarginTop(2).MarginLeft(5)

	infoEntryLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("49")).
			Bold(true)

	infoEntry = lipgloss.NewStyle().Margin(1)

	dialogContainerStyle = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center)

	windowTooSmallStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(highlightColor).
				Align(lipgloss.Center, lipgloss.Center).
				Padding(2, 2)
	containerRunningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("41"))
	containerExitedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("172"))
	containerCreatedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("118"))
)
