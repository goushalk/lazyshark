package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"packetB/internal/analyzer"
)

type AppModel struct {
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

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	var updatedPacketList tea.Model
	updatedPacketList, cmd = m.packetList.Update(msg)
	m.packetList = updatedPacketList.(packListModel)
	return m, cmd
}

func (m AppModel) View() string {
	return m.packetList.View()
}

func NewAppModel(filePath string) (AppModel, error) {
	analyzerResult, err := analyzer.Analyzer(filePath)
	if err != nil {
		log.Fatalf("Error analyzing pcap file: %v", err)
		return AppModel{}, err
	}

	packetListModel := NewPacketListModel(*analyzerResult)

	return AppModel{
		packetList: packetListModel,
	}, nil
}

func StartTUI(filePath string) error {
	initialModel, err := NewAppModel(filePath)
	if err != nil {
		return err
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
	return nil
}
