package tui

import (
	"fmt"
	"packetB/internal/analyzer"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type packListModel struct {
	table     table.Model
	allpacket []analyzer.PacketSummary
	width     int
	height    int
}

var (
	cursorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4"))

	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))
)

func (m packListModel) Init() tea.Cmd {
	return nil
}

func (m packListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(m.width - baseStyle.GetHorizontalFrameSize())
		m.table.SetHeight(m.height - baseStyle.GetVerticalFrameSize()) // Adjust height for border and potential status line

	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			selected := m.table.Cursor()
			return m, func() tea.Msg {
				return PacketSelectedMsg{
					Index: selected,
				}
			}
		}

	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m packListModel) View() string {
	return baseStyle.Render(m.table.View())
}

func NewPacketListModel(analyzerResult analyzer.AnalyzerResult) packListModel {
	columns := []table.Column{
		{Title: "No.", Width: 5},
		{Title: "Time", Width: 15},
		{Title: "Source", Width: 15},
		{Title: "Destination", Width: 15},
		{Title: "Protocol", Width: 10},
		{Title: "Length", Width: 8},
		{Title: "Info", Width: 40},
	}

	rows := []table.Row{}
	for _, pkt := range analyzerResult.Packets {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", pkt.Number),
			pkt.TimeStamp,
			pkt.SrcIp,
			pkt.DstIp,
			pkt.Protocol,
			fmt.Sprintf("%d", pkt.Length),
			pkt.Info,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10), // Initial height, will be adjusted by WindowSizeMsg
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4"))
	t.SetStyles(s)

	return packListModel{
		table:     t,
		allpacket: analyzerResult.Packets,
	}
}
