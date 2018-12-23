package main

import (
	"fmt"
	"math"
	"os"
)

type Point struct {
	X, Y int
}

type Velocity struct {
	X, Y int
}

type Light struct {
	Position Point
	Speed    Velocity
}

type World struct {
	Lights []*Light
}

func Min(a, b int) int {
	if a > b {
		return b
	}

	return a
}

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func (l Light) String() string {
	return fmt.Sprintf("position=<%d, %d> velocity=<%d, %d>", l.Position.X, l.Position.Y, l.Speed.X, l.Speed.Y)
}

func (l *Light) ForwardTime(seconds int) {
	l.Position.X += seconds * l.Speed.X
	l.Position.Y += seconds * l.Speed.Y
}

func (w *World) ForwardTime(seconds int) (int, Point) {
	minX := math.MaxInt64
	maxX := math.MinInt64
	minY := math.MaxInt64
	maxY := math.MinInt64

	for _, light := range w.Lights {
		//w.Sky[light.Position.X][light.Position.Y] = nil

		light.ForwardTime(seconds)

		//w.Sky[light.Position.X][light.Position.Y] = light

		maxX = Max(maxX, light.Position.X)
		minX = Min(minX, light.Position.X)
		maxY = Max(maxY, light.Position.Y)
		minY = Min(minY, light.Position.Y)
	}

	return (maxX - minX) * (maxY - minY), Point{X: maxX, Y: maxY}
}

func (w *World) AddLight(light *Light) {
	//w.Sky[light.Position.X][light.Position.Y] = light
	w.Lights = append(w.Lights, light)
}

func (w *World) Print(fartherAwayPoint Point) {
	grid := make([][]*Light, fartherAwayPoint.X+1)
	for i := range grid {
		grid[i] = make([]*Light, fartherAwayPoint.Y+1)
	}

	for _, light := range w.Lights {
		grid[light.Position.X][light.Position.Y] = light
	}

	for y := 0; y < fartherAwayPoint.Y+1; y++ {
		for x := 0; x < fartherAwayPoint.X+1; x++ {
			cell := grid[x][y]

			if cell == nil {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}

		fmt.Println("")
	}

	fmt.Println()
}

func AdvanceToMessage(world *World) (int, Point) {
	minArea := math.MaxInt64
	var point Point

	timeElapsed := 0
	for {
		area, p := world.ForwardTime(1)

		if area > minArea {
			world.ForwardTime(-1)
			break
		}

		minArea = area
		point = p
		timeElapsed++
	}

	return timeElapsed, point
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	lights := ReadLights(os.Args[1])
	world := World{}

	for _, light := range lights {
		world.AddLight(light)
	}

	timeElapsed, fartherAwayPoint := AdvanceToMessage(&world)

	world.Print(fartherAwayPoint)

	fmt.Printf("Time elapsed: %d\n", timeElapsed)
}
