package main

import (
	"image/color"
	"time"

	"github.com/acifani/vita/lib/game"
)

var (
	multiverse []*game.ParallelUniverse

	dead  = color.RGBA{0, 0, 0, 255}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}

	lifecolors  = []color.RGBA{red, green, blue}
	currentlife = 1
	alive       = green

	fullRefreshes  uint
	previousSecond int64
	restartTime    int64

	start time.Time

	// controls which color is used on the LED cube
	ledColor = [3]byte{0x00, 0x00, 0x00}

	// tracks the current FPS of the LED cube for notifications
	fpsValue = [2]byte{0x00, 0x00}
)

func main() {
	initBluetooth()

	multiverse = createUniverses()
	connectUniverses(multiverse)

	// set cube starting color
	nextLifeColor()

	for {
		start = time.Now()

		drawCube()
		display.Display()

		for i := 0; i < len(multiverse); i++ {
			multiverse[i].Read(gamebuffers[i])
		}

		runUniverses(multiverse)

		displayStatsEverySecond()
		resetCubeEveryMinute()
	}
}

func drawCube() {
	for i := range multiverse {
		drawSide(int16(i), multiverse[i], gamebuffers[i])
	}
}

func drawSide(side int16, u *game.ParallelUniverse, gamebuffer []byte) {
	var rows, cols uint32
	c := dead

	for rows = 0; rows < height; rows++ {
		for cols = 0; cols < width; cols++ {
			idx := u.GetIndex(rows, cols)

			switch {
			case u.Cell(idx) == gamebuffer[idx]:
				// no change, so skip
				continue
			case u.Cell(idx) == game.Alive:
				c = alive
			default: // game.Dead
				c = dead
			}

			display.SetPixel(int16(cols)+side*int16(width), int16(rows), c)
		}
	}
}

func nextLifeColor() {
	currentlife = (currentlife + 1) % len(lifecolors)
	alive = lifecolors[currentlife]

	ledColor[0] = alive.R
	ledColor[1] = alive.G
	ledColor[2] = alive.B
}

func displayStatsEverySecond() {
	second := (start.UnixNano() / int64(time.Second))
	if second != previousSecond {
		previousSecond = second
		newFullRefreshes := getFullRefreshes()
		animationTime := time.Since(start)
		animationFPS := int64(10 * time.Second / animationTime)
		print("#", second, " screen=", newFullRefreshes-fullRefreshes, "fps animation=", animationTime.String(), "/", (animationFPS / 10), ".", animationFPS%10, "fps\r\n")

		// update FPS characteristic
		fpsValue[0] = byte(newFullRefreshes - fullRefreshes)
		fpsValue[1] = byte(animationFPS / 10)
		fpsCharacteristic.Write(fpsValue[:])

		fullRefreshes = newFullRefreshes
	}
}

func resetCubeEveryMinute() {
	minute := (start.UnixNano() / int64(time.Minute))
	if minute != restartTime {
		restartTime = minute
		resetUniverses(multiverse)

		drawCube()
		display.Display()

		time.Sleep(time.Second)

		randomizeUniverses(multiverse)
		nextLifeColor()
	}
}

func resetCube() {
	restartTime = (start.UnixNano() / int64(time.Minute))
	resetUniverses(multiverse)
	randomizeUniverses(multiverse)
}
