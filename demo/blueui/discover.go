package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	isService bool
	uuid      string
}

func (i item) Title() string { return i.uuid }

func (i item) Description() string {
	if i.isService {
		return "Service"
	} else {
		return "Characteristic"
	}
}

func (i item) FilterValue() string { return i.uuid }

func (m *model) updateDiscover(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.err
		return m, nil

	case connectDeviceMsg:
		return m, m.connectToDevice(m.devices.SelectedRow()[0])

	case discoverServicesMsg:
		return m, m.discoverServices()

	case servicesDiscoveredMsg:
		m.servicesList = initServicesList(msg.items)

	case tea.WindowSizeMsg:
		w, h := docStyle.GetFrameSize()
		m.w, m.h = msg.Width, msg.Height
		m.servicesList.SetSize(msg.Width-w, msg.Height-h)
	}

	var cmd tea.Cmd
	m.servicesList, cmd = m.servicesList.Update(msg)
	return m, cmd
}

func (m model) discoverView() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	return docStyle.Render(m.servicesList.View())
}

func initServicesList(items []list.Item) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 80, 40)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return l
}
