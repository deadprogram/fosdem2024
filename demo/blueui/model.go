package main

import (
	"errors"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tinygo.org/x/bluetooth"
)

type model struct {
	w, h    int
	err     error
	spinner spinner.Model

	discover bool
	devices  table.Model

	sr chan bluetooth.ScanResult

	connected bool
	device    *bluetooth.Device

	servicesList list.Model
	services     []bluetooth.DeviceService
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

			for _, c := range chars {
				items = append(items, item{
					isService: false,
					uuid:      c.UUID().String(),
				})
			}
		}

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

	if m.discover {
		return m.updateDiscover(msg)
	}

	return m.updateScanning(msg)
}

func (m model) View() string {
	if m.discover {
		return m.discoverView()
	}

	return m.scanView()
}
