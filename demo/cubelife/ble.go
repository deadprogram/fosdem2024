package main

import (
	"image/color"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var ledColor = [3]byte{0x00, 0x00, 0x00} // controls which program is running on the LED cube

var (
	serviceUUID = [16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
	charUUID    = [16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}

	ledCubeCharacteristic bluetooth.Characteristic
)

var connected bool
var disconnected bool = true

func initBluetooth() {
	adapter.SetConnectHandler(func(d bluetooth.Device, c bool) {
		connected = c

		if !connected && !disconnected {
			//clearLEDS()
			disconnected = true
		}

		if connected {
			disconnected = false
		}
	})

	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "TinyGo LifeCube",
	}))
	must("start adv", adv.Start())

	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.NewUUID(serviceUUID),
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &ledCubeCharacteristic,
				UUID:   bluetooth.NewUUID(charUUID),
				Value:  ledColor[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) != 3 {
						return
					}
					ledColor[0] = value[0]
					ledColor[1] = value[1]
					ledColor[2] = value[2]

					alive = color.RGBA{ledColor[0], ledColor[1], ledColor[2], 0}
					resetCube()
				},
			},
		},
	}))
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
