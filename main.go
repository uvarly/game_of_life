package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"
)

const (
	H_LIMIT = 100
	W_LIMIT = 200
	U_LIMIT = 30
	C_LIMIT = math.MaxInt
)

func reset(height int) {
	fmt.Printf("\x1B[%dA", height+4)
}

func main() {
	h := flag.Int("height", 25, fmt.Sprintf("Board height. Acceptable range: [0,%d]", H_LIMIT))
	w := flag.Int("width", 100, fmt.Sprintf("Board width. Acceptable range: [0, %d]", W_LIMIT))
	d := flag.Float64("density", 0.33, "Spawn chance. 0 - no cells spawn. 1 - fills board completely.  Acceptable range: [0.0, 1.0]")
	u := flag.Int("update", 10, fmt.Sprintf("Simulation updates per second. Acceptable range: [1, %d]", U_LIMIT))
	c := flag.Int("count", 1000, fmt.Sprintf("Simulation steps. Acceptable range: [1, %d]", C_LIMIT))

	flag.Parse()

	if *h < 0 || *h > H_LIMIT {
		fmt.Printf("Board scale out of acceptable range [0, %d]\n", H_LIMIT)
		return
	}

	if *w < 0 || *w > W_LIMIT {
		fmt.Printf("Board scale out of acceptable range [0, %d]\n", W_LIMIT)
		return
	}

	if *d < 0.0 || *d > 1.0 {
		fmt.Printf("Density value out of acceptable range [0.0, 1.0]\n")
		return
	}

	if *u < 1 || *u > U_LIMIT {
		fmt.Printf("Simulation update value out of acceptable range [1, %d]\n", U_LIMIT)
		return
	}

	if *c < 1 || *c > C_LIMIT {
		fmt.Printf("Simulation steps value out of acceptable range [1, %d]\n", C_LIMIT)
		return
	}

	gol := NewGameOfLife()
	gol.populate(*h, *w, *d)

	fmt.Println("Turn 0")
	fmt.Println(gol)

	for i := 1; i <= *c; i++ {
		time.Sleep(time.Second / time.Duration(*u))

		gol.step()

		reset(len(gol.board))

		fmt.Printf("Turn %d\n", i)
		fmt.Println(gol)
	}

	os.Exit(0)
}
