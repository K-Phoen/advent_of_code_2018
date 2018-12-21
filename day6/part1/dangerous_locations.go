package main

// See: https://en.wikipedia.org/wiki/Voronoi_diagram and https://en.wikipedia.org/wiki/Fortune%27s_algorithm

import (
	"fmt"
	"os"
)

const WorldSize = 350

type Cell struct {
	place *Place

	closestPlace        *Place
	distanceFromClosest int

	hasSeveralClosestPlaces bool
}

type World struct {
	space  [WorldSize][WorldSize]*Cell
	places []*Place
	areas  map[byte]int
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
	w.areas[place.name] = 0
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

			if cell.hasSeveralClosestPlaces {
				fmt.Print("+")
			} else if cell.closestPlace != nil {
				fmt.Print(string(cell.closestPlace.name))
			} else {
				fmt.Print(".")
			}
		}

		fmt.Print("\n")
	}

	fmt.Print("\n")

	for name := range w.areas {
		fmt.Printf("%s: %d\n", string(name), w.areas[name])
	}
}

func (w *World) ComputeAreas() {
	var dist int
	var cell *Cell

	for _, place := range w.places {
		for i := 0; i < WorldSize; i++ {
			for j := 0; j < WorldSize; j++ {
				cell = w.space[i][j]

				dist = place.DistanceFrom(Coordinate{
					x: i,
					y: j,
				})

				if dist < cell.distanceFromClosest {
					if cell.closestPlace != nil {
						w.areas[cell.closestPlace.name] -= 1
					}

					cell.closestPlace = place
					cell.distanceFromClosest = dist
					cell.hasSeveralClosestPlaces = false
					w.areas[place.name] += 1
				} else if dist == cell.distanceFromClosest {
					if cell.closestPlace != nil {
						w.areas[cell.closestPlace.name] -= 1
					}

					cell.closestPlace = nil
					cell.hasSeveralClosestPlaces = true
				}
			}
		}
	}
}

func (w *World) FindBiggestFiniteArea() (byte, int) {
	infiniteAreas := make(map[byte]bool)
	max := 0
	name := byte('.')

	for x := 0; x < WorldSize; x++ {
		cell := w.space[x][0]

		if !cell.hasSeveralClosestPlaces {
			infiniteAreas[cell.closestPlace.name] = true
		}

		cell = w.space[x][WorldSize-1]

		if !cell.hasSeveralClosestPlaces {
			infiniteAreas[cell.closestPlace.name] = true
		}
	}

	for y := 0; y < WorldSize; y++ {
		cell := w.space[0][y]

		if !cell.hasSeveralClosestPlaces {
			infiniteAreas[cell.closestPlace.name] = true
		}

		cell = w.space[WorldSize-1][y]

		if !cell.hasSeveralClosestPlaces {
			infiniteAreas[cell.closestPlace.name] = true
		}

	}

	for a := range w.areas {
		if _, exists := infiniteAreas[a]; exists {
			continue
		}

		if w.areas[a] > max {
			max = w.areas[a]
			name = a
		}
	}

	return name, max
}

func NewWorld() *World {
	w := World{
		areas: make(map[byte]int),
	}

	for i := 0; i < WorldSize; i++ {
		for j := 0; j < WorldSize; j++ {
			w.space[i][j] = &Cell{
				distanceFromClosest:     WorldSize * 4,
				hasSeveralClosestPlaces: false,
			}
		}
	}

	return &w
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

	world.ComputeAreas()

	//world.Output()

	biggestArea, area := world.FindBiggestFiniteArea()
	fmt.Printf("Biggest area is %c: %d\n", biggestArea, area)
}
