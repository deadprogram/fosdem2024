//go:build metro_m4_airlift

package main

import (
	"machine"

	"github.com/tinygo-org/things/hub75"
)

var display = hub75.New(hub75.Config{
	SPI:          machine.SPI1,
	Data:         machine.NoPin,
	Clock:        machine.NoPin,
	Latch:        machine.D10,
	OutputEnable: machine.D9,
	A:            machine.D8,
	B:            machine.D7,
	C:            machine.D6,
	D:            machine.D5,
	NumScreens:   6,
	Brightness:   10,
})
