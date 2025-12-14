package tui

import (
	"github.com/goushalk/lazyshark/internal/analyzer"
	tea "github.com/charmbracelet/bubbletea"
)

type hexViewModel struct {
	hexData string
}

func (m *hexViewModel) Load(data []byte) {
	m.hexData = analyzer.DumpHex(data)
}

func (m hexViewModel) Init() tea.Cmd {
	return nil
}

func (m hexViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "backspace" {
			return m, func() tea.Msg {
				return BackToListMsg{}
			}
		}
	}
	return m, nil
}

func (m hexViewModel) View() string {
	if m.hexData == "" {
		return "no hex data"
	}
	return m.hexData
}
