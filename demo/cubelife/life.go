package main

import (
	"sync"

	"github.com/acifani/vita/lib/game"
)

var (
	multirows = 2
	multicols = 3

	height      uint32 = 32
	width       uint32 = 32
	population         = 35
	gamebuffers [][]byte
)

func createUniverses() []*game.ParallelUniverse {
	multi := []*game.ParallelUniverse{}
	for i := 0; i < 6; i++ {
		u := game.NewParallelUniverse(height, width)
		u.Randomize(population)

		multi = append(multi, u)
		gamebuffers = append(gamebuffers, make([]byte, height*width))
	}

	return multi
}

func connectUniverses(multi []*game.ParallelUniverse) {
	multi[0].SetTopNeighbor(multi[3])
	multi[0].SetBottomNeighbor(multi[3])
	multi[0].SetLeftNeighbor(multi[2])
	multi[0].SetRightNeighbor(multi[1])

	multi[1].SetTopNeighbor(multi[4])
	multi[1].SetBottomNeighbor(multi[4])
	multi[1].SetLeftNeighbor(multi[0])
	multi[1].SetRightNeighbor(multi[2])

	multi[2].SetTopNeighbor(multi[5])
	multi[2].SetBottomNeighbor(multi[5])
	multi[2].SetLeftNeighbor(multi[1])
	multi[2].SetRightNeighbor(multi[0])

	multi[3].SetTopNeighbor(multi[0])
	multi[3].SetBottomNeighbor(multi[0])
	multi[3].SetLeftNeighbor(multi[5])
	multi[3].SetRightNeighbor(multi[4])

	multi[4].SetTopNeighbor(multi[1])
	multi[4].SetBottomNeighbor(multi[1])
	multi[4].SetLeftNeighbor(multi[3])
	multi[4].SetRightNeighbor(multi[5])

	multi[5].SetTopNeighbor(multi[2])
	multi[5].SetBottomNeighbor(multi[2])
	multi[5].SetLeftNeighbor(multi[4])
	multi[5].SetRightNeighbor(multi[3])
}

func runUniverses(multi []*game.ParallelUniverse) {
	var wg sync.WaitGroup
	for _, u := range multi {
		callMultiTick(&wg, u)
	}
	wg.Wait()
}

func callMultiTick(wg *sync.WaitGroup, u *game.ParallelUniverse) {
	wg.Add(1)
	go func() {
		u.MultiTick()
		wg.Done()
	}()
}

func resetUniverses(multi []*game.ParallelUniverse) {
	for _, u := range multi {
		u.Reset()
	}
}

func randomizeUniverses(multi []*game.ParallelUniverse) {
	for _, u := range multi {
		u.Randomize(population)
	}
}
