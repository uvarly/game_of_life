package main

import (
	"bytes"
	"math"
	"math/rand"
	"sync"
	"time"
)

type GameOfLife struct {
	wg sync.WaitGroup

	board   [][]uint8
	workers []chan<- struct{}
}

func NewGameOfLife() *GameOfLife {
	return &GameOfLife{}
}

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

func splitWorkload(tasks, workers int) []int {
	var (
		workloads  = make([]int, workers)
		breakpoint = workers - (tasks % workers)
	)

	for i := range workloads {
		workloads[i] = tasks / workers
		if i >= breakpoint {
			workloads[i]++
		}
	}

	return workloads
}
