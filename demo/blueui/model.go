package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tinygo.org/x/bluetooth"
)

type model struct {
	w, h int
	err  error

	state        string
	spinner      spinner.Model
	devices      table.Model
	servicesList list.Model

	connected bool
	sr        chan bluetooth.ScanResult
	device    *bluetooth.Device
	services  []bluetooth.DeviceService
}

type connectDeviceMsg struct{}
type discoverServicesMsg struct{}
type servicesDiscoveredMsg struct{ items []list.Item }
type errMsg struct{ err error }

func newModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner: s,
		devices: initDevicesTable(),
		sr:      make(chan bluetooth.ScanResult),
	}
}

func (m model) Init() tea.Cmd {
	m.state = "scanning"

	return tea.Batch(
		m.spinner.Tick,
		m.scanForDevices(),
		m.waitForDevice(),
	)
}

func (m model) scanForDevices() tea.Cmd {
	return func() tea.Msg {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			m.sr <- device
		})
		must("start scan", err)
		return nil
	}
}

func (m model) waitForDevice() tea.Cmd {
	return func() tea.Msg {
		return <-m.sr
	}
}

func (m *model) connectToDevice(address string) tea.Cmd {
	return func() tea.Msg {
		a := bluetooth.MACAddress{}
		a.Set(address)
		addr := bluetooth.Address{MACAddress: a}

		d, err := adapter.Connect(addr, bluetooth.ConnectionParams{})
		if err != nil {
			return errMsg{err}
		}

		m.device = &d
		m.connected = true
		m.state = "discovering"
		return discoverServicesMsg{}
	}
}

func (m *model) discoverServices() tea.Cmd {
	return func() tea.Msg {
		var err error
		if m.device == nil {
			return errMsg{errors.New("no device")}
		}
		m.services, err = m.device.DiscoverServices(nil)
		if err != nil {
			println(err.Error())
			return errMsg{err}
		}

		items := make([]list.Item, 0)
		for _, s := range m.services {
			items = append(items, item{
				isService: true,
				uuid:      s.UUID().String(),
			})
			chars, err := s.DiscoverCharacteristics(nil)
			if err != nil {
				return errMsg{err}
			}

			// buffer to retrieve characteristic data
			buf := make([]byte, 255)
			for _, c := range chars {
				description := ""

				mtu, err := c.GetMTU()
				if err != nil {
					description = fmt.Sprintf(" mtu: %s", err.Error())
				} else {
					description = fmt.Sprintf(" mtu: %d", mtu)
				}

				n, err := c.Read(buf)
				if err != nil {
					description += fmt.Sprintf(" data: %s", err.Error())
				} else {
					description += fmt.Sprintf(" data(%d): %s", n, string(buf[:n]))
				}

				items = append(items, item{
					isService:   false,
					uuid:        c.UUID().String(),
					description: description,
				})
			}
		}

		m.state = "discovered"

		return servicesDiscoveredMsg{items}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			if m.connected {
				m.device.Disconnect()
			}

			return m, tea.Quit
		}
	}

	switch m.state {
	case "scanning":
		return m.updateScan(msg)
	case "discovering":
		return m.updateDiscover(msg)
	case "discovered":
		return m.updateDiscover(msg)
	}

	return m.updateScan(msg)
}

func (m model) View() string {
	switch m.state {
	case "scanning":
		return m.scanView()
	case "discovering":
		return m.discoveringView()
	case "discovered":
		return m.discoverView()
	}

	return m.scanView()
}
