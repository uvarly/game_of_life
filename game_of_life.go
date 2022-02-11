package main

import (
	"bytes"
	"math"
	"math/rand"
	"sync"
	"time"
)

// GameOfLife is an implementation of John Conway's Game of Life.
//
// The simulation is run in a configurable rectangular window through standart
// output.
type GameOfLife struct {
	wg      sync.WaitGroup
	board   [][]uint8
	workers []chan<- struct{}
}

// NewGameOfLife returns a new GameOfLife that will run a simulation
// determined by the rules of John Conway's Game of Life
func NewGameOfLife() *GameOfLife {
	return &GameOfLife{}
}

// Populate fills the hidden board of configurable scale with cells.
//
// The probability for a cell to appear is determined by the passed
// floating point value ranging between 0 and 1.
func (g *GameOfLife) Populate(h, w int, d float64) {
	threshold := int(100 * d)

	rand.Seed(time.Now().UnixNano())

	g.board = make([][]uint8, h)
	for i := range g.board {
		g.board[i] = make([]uint8, w)
		for j := range g.board[i] {
			if rand.Intn(100) < threshold {
				g.board[i][j] = 1
			}
		}
	}

	g.generateWorkers()
}

// Step progresses simulation by one step.
//
// It takes into account the neighbour count of each individual cell
// and determines the new state of each cell simultaneously.
func (g *GameOfLife) Step() {
	g.wg.Add(len(g.workers))
	for _, exec := range g.workers {
		exec <- struct{}{}
	}
	g.wg.Wait()
}

func (g *GameOfLife) String() string {
	var buf bytes.Buffer

	for i := -1; i <= len(g.board[0]); i++ {
		if i == -1 || i == len(g.board[0]) {
			buf.WriteByte(byte('+'))
		} else {
			buf.WriteByte(byte('-'))
		}
	}
	buf.WriteByte('\n')

	for i := 0; i < len(g.board); i++ {
		for j := -1; j <= len(g.board[i]); j++ {
			if j == -1 || j == len(g.board[i]) {
				buf.WriteByte(byte('|'))
				continue
			}
			if g.board[i][j]&1 == 1 {
				buf.WriteByte(byte('*'))
			} else {
				buf.WriteByte(byte(' '))
			}
		}
		buf.WriteByte('\n')
	}

	for i := 0; i < len(g.board[0])+2; i++ {
		if i == 0 || i == len(g.board[0])+1 {
			buf.WriteByte(byte('+'))
		} else {
			buf.WriteByte(byte('-'))
		}
	}
	buf.WriteByte('\n')

	return buf.String()
}

// willLive tells whether cell survives the next emulation step or not.
func (g *GameOfLife) willLive(i, j int) bool {
	var neighbourCount uint8

	for _, k := range [3]int{-1, 0, 1} {
		for _, l := range [3]int{-1, 0, 1} {
			if k == 0 && l == 0 {
				continue
			}
			if i+k >= 0 && i+k < len(g.board) && j+l >= 0 && j+l < len(g.board[i]) {
				neighbourCount += g.board[i+k][j+l] & 1
			}
		}
	}

	if g.board[i][j]&1 == 1 && (neighbourCount == 2 || neighbourCount == 3) {
		return true
	}

	if g.board[i][j]&1 == 0 && neighbourCount == 3 {
		return true
	}

	return false
}

// work evaluates a designated sector of a board by the rules of Game of Life.
//
// It is designed to run concurrently
func (g *GameOfLife) work(exec <-chan struct{}, begin, end int) {
	for range exec {
		for i := begin; i < end; i++ {
			for j := range g.board[i] {
				if g.willLive(i, j) {
					g.board[i][j] |= 2
				}
			}
		}

		for i := begin; i < end; i++ {
			for j := range g.board[i] {
				g.board[i][j] >>= 1
			}
		}

		g.wg.Done()
	}
}

// generateWorkers slpits the board into roughly equal chunks and allocates
// a worker to handle each.
func (g *GameOfLife) generateWorkers() {
	var (
		taskCount   = len(g.board)
		workerCount = int(math.Log2(float64(taskCount))) + 1
		workloads   = splitWorkload(taskCount, workerCount)
		i           = 0
	)

	for _, wl := range workloads {
		exec := make(chan struct{})

		go g.work(exec, i, i+wl)
		g.workers = append(g.workers, exec)

		i += wl
	}
}

// splitWorkload returns a slice containing a workload to be given to
// each of the workers.
//
// It is designed to split a given amount of work such that difference between
// the smallest and the largest workload is minimal
func splitWorkload(workAmount, workerAmount int) []int {
	var (
		workloads  = make([]int, workerAmount)
		breakpoint = workerAmount - (workAmount % workerAmount)
	)

	for i := range workloads {
		workloads[i] = workAmount / workerAmount
		if i >= breakpoint {
			workloads[i]++
		}
	}

	return workloads
}
