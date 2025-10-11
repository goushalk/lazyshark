package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"packetB/internal/analyzer"
)

// model to view hex data
type hexViewModel struct {
	hexData string
}

// load packet data into hexViewModel

func (m *hexViewModel) Load(filePath string, packetIndex int) {
	hexData, err := analyzer.HexDumper(filePath, packetIndex)

	if err != nil {
		m.hexData = fmt.Sprintf("error : %v", err)
	} else {
		m.hexData = hexData // put the ascii + hex string to the hexViewModel
	}
}

func (m hexViewModel) Init() tea.Cmd {
	return nil
}

func (m hexViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "backspace":
			return m, func() tea.Msg {
				return BackToListMsg{} // tells appmodel to go to list packets view
			}
		}
	}
	return m, nil
}

func (m hexViewModel) View() string {
	if m.hexData == "" {
		return "no ascii + hex"
	}
	return m.hexData
}
