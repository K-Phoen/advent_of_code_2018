package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Coordinate struct {
	x, y int
}

type Place struct {
	name   byte
	coords Coordinate
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%d, %d)", c.x, c.y)
}

func (p *Place) String() string {
	return fmt.Sprintf("%s %s", string(p.name), p.coords)
}

func stringToInt(input string) int {
	item, err := strconv.Atoi(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Incorrect input found: '%v' (could not be converted to a number: %s)\n", item, err.Error())
		os.Exit(2)
	}

	return item
}

func extractX(input string) int {
	comma := strings.Index(input, ",")

	return stringToInt(input[:comma])
}

func extractY(input string) int {
	space := strings.Index(input, " ")

	return stringToInt(input[space+1:])
}

func readPlaces(fileName string) []*Place {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file '%s'\n", fileName)
		os.Exit(1)
	}
	defer file.Close()

	var places []*Place

	scanner := bufio.NewScanner(file)

	nextName := byte('A')
	for scanner.Scan() {
		line := scanner.Text()

		places = append(places, &Place{
			name: nextName,
			coords: Coordinate{
				x: extractX(line),
				y: extractY(line),
			},
		})

		if nextName == 'Z' {
			nextName = 'a'
		} else {
			nextName++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(2)
	}

	return places
}
