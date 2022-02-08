package main

import (
	"bytes"
	"errors"
	"math/rand"
	"time"
)

type gameOfLife struct {
	board [][]int
}

func NewGameOfLife() *gameOfLife {
	return &gameOfLife{}
}

func (gol *gameOfLife) willLive(i, j int) bool {
	var neighbourCount int

	for _, k := range [3]int{-1, 0, 1} {
		for _, l := range [3]int{-1, 0, 1} {
			if k == 0 && l == 0 {
				continue
			}
			if i+k >= 0 && i+k < len(gol.board) && j+l >= 0 && j+l < len(gol.board[i]) {
				neighbourCount += gol.board[i+k][j+l] & 1
			}
		}
	}

	if gol.board[i][j]&1 == 1 && (neighbourCount == 2 || neighbourCount == 3) {
		return true
	}

	if gol.board[i][j]&1 == 0 && neighbourCount == 3 {
		return true
	}

	return false
}

func (gol *gameOfLife) step() {
	for i := range gol.board {
		for j := range gol.board[i] {
			if gol.willLive(i, j) {
				gol.board[i][j] |= 2
			}
		}
	}

	for i := range gol.board {
		for j := range gol.board[i] {
			gol.board[i][j] >>= 1
		}
	}
}

func (gol *gameOfLife) populate(h, w int, d float64) error {
	if d < 0 || d > 1 {
		return errors.New("density value out of range: [0.0, 1.0]")
	}

	rand.Seed(time.Now().UnixNano())
	threshold := int(100 * d)

	gol.board = make([][]int, h)
	for i := range gol.board {
		gol.board[i] = make([]int, w)
		for j := range gol.board[i] {
			if rand.Intn(100) < threshold {
				gol.board[i][j] = 1
			}
		}
	}

	return nil
}

func (gol *gameOfLife) String() string {
	var buf bytes.Buffer

	for i := -1; i <= len(gol.board[0]); i++ {
		if i == -1 || i == len(gol.board[0]) {
			buf.WriteByte(byte('+'))
		} else {
			buf.WriteByte(byte('-'))
		}
	}
	buf.WriteByte('\n')

	for i := 0; i < len(gol.board); i++ {
		for j := -1; j <= len(gol.board[i]); j++ {
			if j == -1 || j == len(gol.board[i]) {
				buf.WriteByte(byte('|'))
				continue
			}
			if gol.board[i][j]&1 == 1 {
				buf.WriteByte(byte('*'))
			} else {
				buf.WriteByte(byte(' '))
			}
		}
		buf.WriteByte('\n')
	}

	for i := 0; i < len(gol.board[0])+2; i++ {
		if i == 0 || i == len(gol.board[0])+1 {
			buf.WriteByte(byte('+'))
		} else {
			buf.WriteByte(byte('-'))
		}
	}
	buf.WriteByte('\n')

	return buf.String()
}
