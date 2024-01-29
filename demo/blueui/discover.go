package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	isService   bool
	uuid        string
	description string
}

func (i item) Title() string {
	typ := "Characteristic"
	if i.isService {
		typ = "Service"
	}

	return fmt.Sprintf("%s %s", typ, i.uuid)
}

func (i item) Description() string {
	return i.description
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
		name := m.devices.SelectedRow()[0]
		if m.devices.SelectedRow()[2] != "" {
			name = m.devices.SelectedRow()[2]
		}
		m.servicesList = initServicesList(name, msg.items)

	case tea.WindowSizeMsg:
		w, h := docStyle.GetFrameSize()
		m.w, m.h = msg.Width, msg.Height
		m.servicesList.SetSize(msg.Width-w, msg.Height-h)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.servicesList, cmd = m.servicesList.Update(msg)
	return m, cmd
}

func (m model) discoveringView() string {
	if m.err != nil {
		return fmt.Sprintf("\nERROR: %v\n\n", m.err)
	}

	str := fmt.Sprintf("\n\n   %s Discovering services...\n\n", m.spinner.View())
	return str
}

func (m model) discoverView() string {
	if m.err != nil {
		return fmt.Sprintf("\nERROR: %v\n\n", m.err)
	}

	return docStyle.Render(m.servicesList.View())
}

func initServicesList(title string, items []list.Item) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 80, 40)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return l
}
