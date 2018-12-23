package main

import (
	"fmt"
	"math"
)

const GridSerialNumber = 8868
const GridSize = 300

type Grid [GridSize][GridSize]int

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

func ComputeCellsFuelLevel(grid *Grid) {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			grid[x][y] = ComputeCellFuelLevel(x, y)
		}
	}
}

func (grid Grid) Print() {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			fmt.Printf("%d ", grid[x][y])
		}

		fmt.Println()
	}
}

func ComputePowerSquare(grid *Grid, topLeftX int, topLeftY int, squareSize int) int {
	power := 0

	for y := topLeftY; y < topLeftY+squareSize; y++ {
		for x := topLeftX; x < topLeftX+squareSize; x++ {
			power += grid[x][y]
		}
	}

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
	var fuelGrid Grid

	ComputeCellsFuelLevel(&fuelGrid)
	//fuelGrid.Print()

	maxPower, topLeftX, topLeftY := FindHighestPowerSquare(&fuelGrid, 3)

	fmt.Printf("Highest power square starts at %d,%d (total power of %d)\n", topLeftX, topLeftY, maxPower)
}
