package main

import (
	"bytes"
	"math/rand"
	"time"
)

type gameOfLife struct {
	board [][]int
}

func NewGameOfLife() *gameOfLife {
	return &gameOfLife{}
}

func (g *gameOfLife) willLive(i, j int) bool {
	var neighbourCount int

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

func (g *gameOfLife) step() {
	for i := range g.board {
		for j := range g.board[i] {
			if g.willLive(i, j) {
				g.board[i][j] |= 2
			}
		}
	}

	for i := range g.board {
		for j := range g.board[i] {
			g.board[i][j] >>= 1
		}
	}
}

func (g *gameOfLife) populate(h, w int, d float64) {
	rand.Seed(time.Now().UnixNano())
	threshold := int(100 * d)

	g.board = make([][]int, h)
	for i := range g.board {
		g.board[i] = make([]int, w)
		for j := range g.board[i] {
			if rand.Intn(100) < threshold {
				g.board[i][j] = 1
			}
		}
	}
}

func (g *gameOfLife) String() string {
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
