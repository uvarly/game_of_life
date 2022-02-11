package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"
)

const (
	hLimit = 1_000
	wLimit = 2_000
	uLimit = 120
	cLimit = math.MaxInt
)

func reset(height int) {
	fmt.Printf("\x1B[%dA", height+4)
}

func main() {
	var (
		h = flag.Int("height", 50, fmt.Sprintf("Board height. Acceptable range: [1, %d]", hLimit))
		w = flag.Int("width", 150, fmt.Sprintf("Board width. Acceptable range: [1, %d]", wLimit))
		d = flag.Float64("density", 0.33, "Spawn chance. 0 - no cells spawn. 1 - fills board completely.  Acceptable range: [0.0, 1.0]")
		u = flag.Int("update", 10, fmt.Sprintf("Simulation updates per second. Acceptable range: [1, %d]", uLimit))
		c = flag.Int("count", 1000, fmt.Sprintf("Simulation steps. Acceptable range: [1, %d]", cLimit))
	)

	flag.Parse()

	if *h < 1 || *h > hLimit {
		fmt.Printf("Board scale out of acceptable range [1, %d]\n", hLimit)
		return
	}

	if *w < 1 || *w > wLimit {
		fmt.Printf("Board scale out of acceptable range [1, %d]\n", wLimit)
		return
	}

	if *d < 0.0 || *d > 1.0 {
		fmt.Printf("Density value out of acceptable range [0.0, 1.0]\n")
		return
	}

	if *u < 1 || *u > uLimit {
		fmt.Printf("Simulation update value out of acceptable range [1, %d]\n", uLimit)
		return
	}

	if *c < 1 || *c > cLimit {
		fmt.Printf("Simulation steps value out of acceptable range [1, %d]\n", cLimit)
		return
	}

	var (
		g      = NewGameOfLife()
		ticker = time.NewTicker(time.Second / time.Duration(*u))
		i      = 1
	)

	g.Populate(*h, *w, *d)

	fmt.Println("Turn 0")
	fmt.Println(g)

	for range ticker.C {
		if i++; i > *c {
			break
		}

		time.Sleep(time.Second / time.Duration(*u))

		g.Step()

		reset(len(g.board))

		fmt.Printf("Turn %d\n", i)
		fmt.Println(g)
	}

	os.Exit(0)
}
