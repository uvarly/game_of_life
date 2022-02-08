package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

var clear func()

func init() {
	switch runtime.GOOS {
	case "linux":
		clear = func() {
			// cmd := exec.Command("clear")
			// cmd.Stdout = os.Stdout
			// cmd.Run()
			fmt.Print("\x1B[2J")
		}
	case "windows":
		clear = func() {
			// cmd := exec.Command("cmd", "/c", "cls")
			// cmd.Stdout = os.Stdout
			// cmd.Run()
			fmt.Print("\x1B[2J")
		}
	}
}

func main() {
	h := flag.Int("height", 25, "Board height. Acceptable range: [0, 100]")
	w := flag.Int("width", 50, "Board width. Acceptable range: [0, 100]")
	d := flag.Float64("density", 0.25, "Spawn chance. 0 - no cells spawn. 1 - fills board completely.  Acceptable range: [0.0, 1.0]")
	u := flag.Int("update", 5, "Simulation updates per second. Acceptable range: [1, 20]")
	c := flag.Int("count", 200, "Simulation steps. Acceptable range: [1, 1000]")

	flag.Parse()

	if *h < 0 || *h > 100 || *w < 0 || *w > 100 {
		fmt.Println("Board scale out of acceptable range [0, 100]")
		return
	}

	if *d < 0.0 || *d > 1.0 {
		fmt.Println("Density value out of acceptable range [0.0, 1.0]")
		return
	}

	if *u < 1 || *u > 20 {
		fmt.Println("Simulation update value out of acceptable range [1, 20]")
		return
	}

	if *c < 1 || *c > 1000 {
		fmt.Println("Simulation steps value out of acceptable range [1, 1000]")
		return
	}

	gol := NewGameOfLife()
	gol.populate(*h, *w, *d)

	clear()

	fmt.Println("Turn 0")
	fmt.Println(gol)

	for i := 1; i <= *c; i++ {
		time.Sleep(time.Second / time.Duration(*u))

		gol.step()
		clear()
		fmt.Printf("Turn %d\n", i)
		fmt.Println(gol)
	}

	os.Exit(0)
}
