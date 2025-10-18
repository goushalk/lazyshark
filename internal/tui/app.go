package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"log"
	"packetB/internal/analyzer"
)

type AppModel struct {
	filePath    string
	currentView string
	packetList  packListModel
	hexView     hexViewModel

	help help.Model
	keys keyMap
}

type keyMap struct {
	Quit   key.Binding
	Select key.Binding
	Back   key.Binding
}

type PacketSelectedMsg struct {
	Index int
}

type BackToListMsg struct{}

func newKeyMap() keyMap {
	return keyMap{
		Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
		Select: key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "view hex")),
		Back:   key.NewBinding(key.WithKeys("backspace"), key.WithHelp("←", "back")),
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.packetList.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.currentView {
	case "list":
		// delegate update to packetList
		var updated tea.Model
		updated, cmd = m.packetList.Update(msg)
		m.packetList = updated.(packListModel)

		// check if PacketSelectedMsg came from packetList
		if packetMsg, ok := msg.(PacketSelectedMsg); ok {
			m.currentView = "hex"
			m.hexView.Load(m.filePath, packetMsg.Index) // give the hex view the data
		}

	case "hex":
		// delegate update to hexView
		var updated tea.Model
		updated, cmd = m.hexView.Update(msg)
		m.hexView = updated.(hexViewModel)

		// check if BackToListMsg came from hexView
		if _, ok := msg.(BackToListMsg); ok {
			m.currentView = "list"
		}
	}

	return m, cmd
}

func (m AppModel) View() string {

	var content string
	switch m.currentView {
	case "hex":
		content = m.hexView.View()

	default:
		content = m.packetList.View()
	}

	footer := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("5")).
		MarginTop(0).
		MarginBottom(0).
		Render(m.help.ShortHelpView([]key.Binding{
			m.keys.Select, m.keys.Back, m.keys.Quit,
		}))

	return lipgloss.JoinVertical(lipgloss.Center, content, footer)

}

func NewAppModel(filePath string) (AppModel, error) {
	analyzerResult, err := analyzer.Analyzer(filePath)
	if err != nil {
		log.Fatalf("Error analyzing pcap file: %v", err)
		return AppModel{}, err
	}

	packetListModel := NewPacketListModel(*analyzerResult)

	return AppModel{
		filePath:    filePath,
		currentView: "list",
		packetList:  packetListModel,
		hexView:     hexViewModel{},

		help: help.New(),
		keys: newKeyMap(),
	}, nil
}

func StartTUI(filePath string) error {
	initialModel, err := NewAppModel(filePath)
	if err != nil {
		return err
	}

	p := tea.NewProgram(&initialModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
	return nil
}
