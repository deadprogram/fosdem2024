package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tinygo.org/x/bluetooth"
)

func (m model) updateScanning(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case bluetooth.ScanResult:
		r := table.Row{
			msg.Address.String(),
			strconv.Itoa(int(msg.RSSI)),
			msg.LocalName(),
		}

		devices := m.devices.Rows()
		if len(devices) > maxRows {
			devices = devices[1:]
		}
		devices = append(devices, r)

		m.devices.SetRows(devices)
		return m, m.waitForDevice()

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.devices.Focused() {
				m.devices.Blur()
			} else {
				m.devices.Focus()
			}
		case " ":
			adapter.StopScan()
		case "enter":
			adapter.StopScan()
			m.discover = true

			// now go discover services
			// m.devices.SelectedRow()[0]
			return m.updateDiscover(connectDeviceMsg{})
		}
	}

	m.devices, cmd = m.devices.Update(msg)
	return m, cmd
}

func (m model) scanView() string {
	return baseStyle.Render(m.devices.View()) + "\n"
}

const maxRows = 36

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
)

func initDevicesTable() table.Model {
	columns := []table.Column{
		{Title: "MAC", Width: 18},
		{Title: "RSSI", Width: 4},
		{Title: "Name", Width: 30},
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(maxRows+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	tbl.SetStyles(s)

	return tbl
}
