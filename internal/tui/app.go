package tui

import (

	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/goushalk/lazyshark/internal/analyzer"
)

type AppModel struct {
	filePath    string
	currentView string
	packetList  packListModel
	hexView     hexViewModel
}

type PacketSelectedMsg struct {
	Index int
}

type BackToListMsg struct{}

func (m AppModel) Init() tea.Cmd {
	return m.packetList.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	switch m.currentView {

	case "list":
		var updated tea.Model
		updated, cmd = m.packetList.Update(msg)
		m.packetList = updated.(packListModel)

		if pktMsg, ok := msg.(PacketSelectedMsg); ok {
			if pktMsg.Index >= 0 && pktMsg.Index < len(m.packetList.allpacket) {
				pkt := m.packetList.allpacket[pktMsg.Index]
				m.hexView.Load(pkt.RawData)
				m.currentView = "hex"
			}
		}

	case "hex":
		var updated tea.Model
		updated, cmd = m.hexView.Update(msg)
		m.hexView = updated.(hexViewModel)

		if _, ok := msg.(BackToListMsg); ok {
			m.currentView = "list"
		}
	}

	return m, cmd
}

func (m AppModel) View() string {
	if m.currentView == "hex" {
		return m.hexView.View()
	}
	return m.packetList.View()
}

func NewAppModel(filePath string) (AppModel, error) {
	result, err := analyzer.Analyzer(filePath)
	if err != nil {
		log.Fatalf("Error analyzing pcap file: %v", err)
	}

	packetListModel := NewPacketListModel(*result)

	return AppModel{
		filePath:    filePath,
		currentView: "list",
		packetList:  packetListModel,
		hexView:     hexViewModel{},
	}, nil
}

func StartTUI(filePath string) error {
	initialModel, err := NewAppModel(filePath)
	if err != nil {
		return err
	}

	p := tea.NewProgram(&initialModel, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
