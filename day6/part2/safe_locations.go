package main

// See: https://en.wikipedia.org/wiki/Voronoi_diagram and https://en.wikipedia.org/wiki/Fortune%27s_algorithm

import (
	"fmt"
	"os"
)

const WorldSize = 500
const SafeDistance = 10000

type Cell struct {
	place *Place

	isSafe bool
}

type World struct {
	space    [WorldSize][WorldSize]*Cell
	places   []*Place
	safeArea int
}

func Abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}

func (p *Place) DistanceFrom(coords Coordinate) int {
	return Abs(p.coords.x-coords.x) + Abs(p.coords.y-coords.y)
}

func (w *World) AddPlace(place *Place) {
	w.space[place.coords.x][place.coords.y].place = place
	w.places = append(w.places, place)
	w.safeArea = 0
}

func (w *World) Output() {
	var cell *Cell

	for y := 0; y < WorldSize; y++ {
		for x := 0; x < WorldSize; x++ {
			cell = w.space[x][y]

			if cell.place != nil {
				fmt.Printf("\033[1m%s\033[0m", string(cell.place.name))
				continue
			}

			if cell.isSafe {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}

		fmt.Print("\n")
	}

	fmt.Print("\n")
}

func NewWorld() *World {
	w := World{safeArea: 0}

	for i := 0; i < WorldSize; i++ {
		for j := 0; j < WorldSize; j++ {
			w.space[i][j] = &Cell{
				isSafe: false,
			}
		}
	}

	return &w
}

func (w *World) ComputeSafeArea() {
	var total int
	var coords Coordinate

	for i := 0; i < WorldSize; i++ {
		for j := 0; j < WorldSize; j++ {
			total = 0
			coords = Coordinate{x: i, y: j}

			for _, place := range w.places {
				total += place.DistanceFrom(coords)
			}

			if total < SafeDistance {
				w.space[i][j].isSafe = true
				w.safeArea += 1
			}
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	world := NewWorld()
	places := readPlaces(os.Args[1])

	for _, place := range places {
		world.AddPlace(place)
	}

	//world.Output()

	world.ComputeSafeArea()

	//world.Output()

	fmt.Printf("Safe area size: %d\n", world.safeArea)
}
