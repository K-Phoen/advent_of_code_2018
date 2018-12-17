package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//const FabricWidth = 11
//const FabricHeight = 9

const FabricWidth = 1000
const FabricHeight = 1000

type Fabric struct {
	material [FabricHeight][FabricWidth][]int
}

type Claim struct {
	id int

	leftEdge int
	topEdge  int

	width  int
	height int
}

func stringToInt(input string) int {
	item, err := strconv.Atoi(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Incorrect input found: '%v' (could not be converted to a number: %s)\n", item, err.Error())
		os.Exit(2)
	}

	return item
}

func extractClaimId(input string) int {
	pos := strings.Index(input, " ")

	return stringToInt(input[1:pos])
}

func extractLeftEdge(input string) int {
	at := strings.Index(input, "@")
	comma := strings.Index(input, ",")

	return stringToInt(input[at+2 : comma])
}

func extractTopEdge(input string) int {
	comma := strings.Index(input, ",")
	colon := strings.Index(input, ":")

	return stringToInt(input[comma+1 : colon])
}

func extractWidth(input string) int {
	colon := strings.Index(input, ":")
	x := strings.Index(input, "x")

	return stringToInt(input[colon+2 : x])
}

func extractHeight(input string) int {
	x := strings.Index(input, "x")

	return stringToInt(input[x+1:])
}

func (fabric *Fabric) lay(claim *Claim) {
	for i := claim.topEdge; i < claim.topEdge+claim.height; i++ {
		for j := claim.leftEdge; j < claim.leftEdge+claim.width; j++ {
			fabric.material[i][j] = append(fabric.material[i][j], claim.id)
		}
	}
}

func (fabric *Fabric) overlappingArea() int {
	area := 0

	for i := 0; i < FabricHeight; i++ {
		for j := 0; j < FabricWidth; j++ {
			cellClaims := fabric.material[i][j]

			if len(cellClaims) > 1 {
				area += 1
			}
		}
	}

	return area
}

func (fabric *Fabric) output() {
	for i := 0; i < FabricHeight; i++ {
		for j := 0; j < FabricWidth; j++ {
			cellClaims := fabric.material[i][j]

			if len(cellClaims) == 0 {
				fmt.Print(".")
			}

			if len(cellClaims) == 1 {
				fmt.Print("#")
			}

			if len(cellClaims) > 1 {
				fmt.Print("X")
			}
		}
		fmt.Print("\n")
	}
}

func readClaims(fileName string) chan Claim {
	channel := make(chan Claim)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file '%s'\n", fileName)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)

	go func() {
		defer close(channel)

		for scanner.Scan() {
			line := scanner.Text()

			channel <- Claim{
				id:       extractClaimId(line),
				leftEdge: extractLeftEdge(line),
				topEdge:  extractTopEdge(line),
				width:    extractWidth(line),
				height:   extractHeight(line),
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			os.Exit(2)
		}
	}()

	return channel
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	fabric := Fabric{}
	for claim := range readClaims(os.Args[1]) {
		fabric.lay(&claim)
		fmt.Printf("%v\n", claim)
	}

	fmt.Printf("Overlapping area: %d\n", fabric.overlappingArea())
}
