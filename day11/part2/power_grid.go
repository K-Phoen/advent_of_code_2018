package main

import (
	"fmt"
	"math"
)

const GridSerialNumber = 8868
const GridSize = 300

type Grid struct {
	grid [GridSize][GridSize]int

	powerSquares map[string]int
}

func HundredsDigit(a int) int {
	return (a % 1000) / 100
}

func ComputeCellFuelLevel(x int, y int) int {
	rackId := x + 10 // rackId

	level := rackId
	level *= y
	level += GridSerialNumber
	level *= rackId

	return HundredsDigit(level) - 5
}

func ComputeCellsFuelLevel(g *Grid) {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			g.grid[x][y] = ComputeCellFuelLevel(x, y)
		}
	}
}

func (g Grid) Print() {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			fmt.Printf("%d ", g.grid[x][y])
		}

		fmt.Println()
	}
}

func ComputePowerSquare(g *Grid, topLeftX int, topLeftY int, squareSize int) int {
	if squareSize <= 0 {
		panic("invalid square size given")
	}

	if squareSize == 1 {
		return g.grid[topLeftX][topLeftY]
	}

	cacheKey := fmt.Sprintf("%d,%d,%d", topLeftX, topLeftY, squareSize)

	if val, exists := g.powerSquares[cacheKey]; exists {
		return val
	}

	power := ComputePowerSquare(g, topLeftX, topLeftY, squareSize-1)

	/*
		squareSize(4)
		→ squareSize(3) + .
		topLeft := 0, 0
		x x x .
		x x x .
		x x x .
		. . . +
	*/

	// right column
	for y := topLeftY; y < topLeftY+squareSize; y++ {
		power += g.grid[topLeftX+(squareSize-1)][y]
	}

	// bottom
	for y := topLeftY + (squareSize - 1); y < topLeftY+squareSize; y++ {
		for x := topLeftX; x < topLeftX+squareSize-1; x++ {
			power += g.grid[x][topLeftY+(squareSize-1)]
		}
	}

	g.powerSquares[cacheKey] = power

	return power
}

func FindHighestPowerSquare(grid *Grid, squareSize int) (int, int, int) {
	maxPower := math.MinInt64
	topLeftX := 0
	topLeftY := 0

	for y := 0; y < GridSize-squareSize; y++ {
		for x := 0; x < GridSize-squareSize; x++ {
			power := ComputePowerSquare(grid, x, y, squareSize)

			if power > maxPower {
				maxPower = power
				topLeftX = x
				topLeftY = y
			}
		}
	}

	return maxPower, topLeftX, topLeftY
}

func main() {
	fuelGrid := Grid{
		powerSquares: make(map[string]int),
	}

	ComputeCellsFuelLevel(&fuelGrid)
	//fuelGrid.Print()

	maxPower := math.MinInt64
	topLeftX, topLeftY, maxSize := 0, 0, 0
	for size := 1; size < GridSize; size++ {
		fmt.Printf("Testing for size %d\n", size)
		power, x, y := FindHighestPowerSquare(&fuelGrid, size)

		if power > maxPower {
			maxPower = power
			maxSize = size
			topLeftX = x
			topLeftY = y
		}
	}

	fmt.Printf("Highest power square starts at %d,%d (total power of %d, for a size of %d)\n", topLeftX, topLeftY, maxPower, maxSize)
}
