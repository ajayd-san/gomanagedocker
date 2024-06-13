package tui

import (
	"github.com/ajayd-san/gomanagedocker/dockercmd"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyLabel   = "label"
	columnKeyImgName = "imgName"
	columnKeyCrit    = "crit"
	columnKeyhigh    = "high"
	columnKeyMed     = "med"
	columnKeyLow     = "low"
	columnKeyUnknown = "unknown"
)

const (
	columnLabelWidth   = 25
	columnImgNameWidth = 25
	columnNumTypeWidth = 10
)

var (
	colorNormal = lipgloss.NewStyle().Foreground(lipgloss.Color("#fa0"))
	colorCrit   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	colorHigh   = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	colorMed    = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	colorLow    = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	colorUnkown = lipgloss.NewStyle().Foreground(lipgloss.Color("129"))
)

var (
	styleBase = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a7a")).
		BorderForeground(lipgloss.Color("#a38")).Align(lipgloss.Center)
)

type TableModel struct {
	inner table.Model
}

func makeRow(label, name, crit, high, med, low, unknown string) table.Row {
	return table.NewRow(table.RowData{
		columnKeyLabel:   label,
		columnKeyImgName: table.NewStyledCell(name, colorNormal),
		columnKeyCrit:    table.NewStyledCell(crit, colorCrit),
		columnKeyhigh:    table.NewStyledCell(high, colorHigh),
		columnKeyMed:     table.NewStyledCell(med, colorMed),
		columnKeyLow:     table.NewStyledCell(low, colorLow),
		columnKeyUnknown: table.NewStyledCell(unknown, colorUnkown),
	})
}

func NewTable(scoutData dockercmd.ScoutData) *TableModel {
	rows := make([]table.Row, 0, len(scoutData.ImageVulEntries))

	for _, scoutEntry := range scoutData.ImageVulEntries {
		row := makeRow(
			scoutEntry.Label,
			scoutEntry.ImageName,
			scoutEntry.Critical,
			scoutEntry.High,
			scoutEntry.Medium,
			scoutEntry.Low,
			scoutEntry.UnknownSeverity,
		)
		rows = append(rows, row)
	}
	return &TableModel{
		inner: table.New([]table.Column{
			table.NewColumn(columnKeyLabel, "", columnLabelWidth),
			table.NewColumn(columnKeyImgName, "Image Name", columnImgNameWidth),
			table.NewColumn(columnKeyCrit, "Critical", columnNumTypeWidth),
			table.NewColumn(columnKeyhigh, "High", columnNumTypeWidth),
			table.NewColumn(columnKeyMed, "Medium", columnNumTypeWidth),
			table.NewColumn(columnKeyLow, "Low", columnNumTypeWidth),
			table.NewColumn(columnKeyUnknown, "Unknown", columnNumTypeWidth),
		}).WithRows(rows).
			BorderRounded().
			WithBaseStyle(styleBase),
	}
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.inner, cmd = m.inner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m TableModel) View() string {
	view := m.inner.View()

	return lipgloss.NewStyle().MarginLeft(1).Render(view)
}

func getDummyTable() *TableModel {
	return &TableModel{
		table.New([]table.Column{
			table.NewColumn(columnKeyLabel, "", 25),
			table.NewColumn(columnKeyImgName, "Image Name", 25),
			table.NewColumn(columnKeyCrit, "Critical", 10),
			table.NewColumn(columnKeyhigh, "High", 10),
			table.NewColumn(columnKeyMed, "Medium", 10),
			table.NewColumn(columnKeyLow, "Low", 10),
			table.NewColumn(columnKeyUnknown, "Unknown", 10),
		}).WithRows([]table.Row{
			makeRow("target", "nginx ", "2300648", "21", "84", "1", "0"),
			makeRow("target", "nginx ", "2300648", "21", "84", "1", "0"),
			makeRow("target", "nginx ", "2300648", "21", "84", "1", "0"),
			makeRow("target", "nginx ", "2300648", "21", "84", "1", "0"),
			makeRow("target", "nginx ", "2300648", "21", "84", "1", "0"),
			makeRow("target", "nginx ", "2300648", "21", "84", "1", "0"),
			makeRow("target", "nginx ", "2300648", "21", "84", "1", "0"),
		}).
			BorderRounded().
			WithBaseStyle(styleBase),
	}
}
