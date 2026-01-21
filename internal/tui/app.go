package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"packetB/internal/analyzer"
)

type AppModel struct {
	result analyzer.AnalyzerResult
	table  table.Model

	selectedIndex int
	detail        string
	hex           string

	showHelp bool

	width  int
	height int
}

/* ───────── Styles ───────── */

var (
	outerBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	separator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	header = lipgloss.NewStyle().
		Bold(true)

	dim = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	helpStyle = lipgloss.NewStyle().
			Padding(1).
			Foreground(lipgloss.Color("252"))

	keybar = lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))
)

/* ───────── Init ───────── */

func NewAppModel(filePath string) (AppModel, error) {
	result, err := analyzer.Analyzer(filePath)
	if err != nil {
		return AppModel{}, err
	}

	t := table.New(table.WithFocused(true))
	t.SetColumns(defaultColumns(80))

	s := table.DefaultStyles()
	s.Header = s.Header.Bold(true)
	s.Selected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("57"))
	t.SetStyles(s)

	m := AppModel{
		result:        *result,
		table:         t,
		selectedIndex: 0,
	}

	m.buildRows()
	m.sync()
	return m, nil
}

func (m AppModel) Init() tea.Cmd { return nil }

/* ───────── Update ───────── */

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = !m.showHelp
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	if idx := m.table.Cursor(); idx != m.selectedIndex {
		m.selectedIndex = idx
		m.sync()
	}

	return m, cmd
}

/* ───────── View ───────── */

func (m AppModel) View() string {
	// Terminal too small
	if m.width < 60 || m.height < 20 {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			"Terminal too small\nResize to continue",
		)
	}

	innerWidth := m.width - 2
	innerHeight := m.height - 2

	sep := separator.Render(strings.Repeat("─", innerWidth))

	keybinds := keybar.Render(
		"↑/↓ navigate   ? help   q quit",
	)

	// Heights
	listH := m.table.Height()
	remaining := innerHeight - listH - 3 // separators + keybar

	if remaining < 6 {
		remaining = 6
	}

	detailH := remaining / 2
	hexH := remaining - detailH

	var middle string
	if m.showHelp {
		middle = helpStyle.Render(
			header.Render("Help") + "\n\n" +
				"↑ / ↓  Navigate packets\n" +
				"?      Toggle help\n" +
				"q      Quit\n",
		)
	} else {
		middle = lipgloss.NewStyle().
			Width(innerWidth).
			Height(detailH).
			Render(header.Render("Packet Details") + "\n" + dim.Render(m.detail))
	}

	bottom := lipgloss.NewStyle().
		Width(innerWidth).
		Height(hexH).
		Render(header.Render("Packet Bytes") + "\n" + dim.Render(m.hex))

	content := strings.Join([]string{
		m.table.View(),
		sep,
		middle,
		sep,
		bottom,
		sep,
		keybinds,
	}, "\n")

	// Render border ONLY if it fits
	if lipgloss.Width(content) <= m.width && lipgloss.Height(content) <= m.height {
		return outerBorder.
			Width(m.width).
			Height(m.height).
			Render(content)
	}

	// Fallback: no border
	return content
}

/* ───────── Layout ───────── */

func (m *AppModel) resize() {
	if m.width == 0 || m.height == 0 {
		return
	}

	usableWidth := m.width - 2

	// Packet list height (header + rows)
	listHeight := 12
	if listHeight < 6 {
		listHeight = 6
	}

	m.table.SetWidth(usableWidth)
	m.table.SetHeight(listHeight)

	m.table.SetColumns(defaultColumns(usableWidth))
	m.buildRows()
}

/* ───────── Columns ───────── */

func defaultColumns(totalWidth int) []table.Column {
	noW := 4
	timeW := 16
	srcW := 14
	dstW := 14
	protoW := 8
	lenW := 6

	used := noW + timeW + srcW + dstW + protoW + lenW
	infoW := totalWidth - used - 6
	if infoW < 10 {
		infoW = 10
	}

	return []table.Column{
		{Title: "No.", Width: noW},
		{Title: "Time", Width: timeW},
		{Title: "Source", Width: srcW},
		{Title: "Destination", Width: dstW},
		{Title: "Proto", Width: protoW},
		{Title: "Len", Width: lenW},
		{Title: "Info", Width: infoW},
	}
}

/* ───────── Rows ───────── */

func (m *AppModel) buildRows() {
	rows := make([]table.Row, 0, len(m.result.Packets))
	for _, p := range m.result.Packets {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", p.Number),
			p.TimeStamp,
			p.SrcIp,
			p.DstIp,
			p.Protocol,
			fmt.Sprintf("%d", p.Length),
			p.Info,
		})
	}
	m.table.SetRows(rows)
}

/* ───────── Sync ───────── */

func (m *AppModel) sync() {
	if m.selectedIndex < 0 || m.selectedIndex >= len(m.result.Packets) {
		return
	}
	p := m.result.Packets[m.selectedIndex]
	m.detail = buildDetails(p)
	m.hex = analyzer.DumpHex(p.RawData)
}

/* ───────── Details Builder ───────── */

func buildDetails(p analyzer.PacketSummary) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Frame %d\n", p.Number)
	fmt.Fprintf(&b, "  Length: %d bytes\n\n", p.Length)

	if p.SrcIp != "N/A" {
		fmt.Fprintf(&b, "Internet Protocol\n")
		fmt.Fprintf(&b, "  Source: %s\n", p.SrcIp)
		fmt.Fprintf(&b, "  Destination: %s\n\n", p.DstIp)
	}

	fmt.Fprintf(&b, "Protocol\n")
	fmt.Fprintf(&b, "  %s\n\n", p.Protocol)

	fmt.Fprintf(&b, "Info\n")
	fmt.Fprintf(&b, "  %s\n", p.Info)

	return b.String()
}

/* ───────── Runner ───────── */

func StartTUI(filePath string) error {
	model, err := NewAppModel(filePath)
	if err != nil {
		return err
	}
	p := tea.NewProgram(&model, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
