package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func StringToInt(input string) int {
	i, err := strconv.Atoi(input)

	if err != nil {
		panic("Unable to convert input to string")
	}

	return i
}

func ReadLights(fileName string) []*Light {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file '%s'\n", fileName)
		os.Exit(1)
	}
	defer file.Close()

	var lights []*Light

	re := regexp.MustCompile("\\-?\\d+")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllString(line, -1)

		if len(matches) != 4 {
			panic("Expected 4 numbers per line")
		}

		lights = append(lights, &Light{
			Position: Point{
				X: StringToInt(matches[0]),
				Y: StringToInt(matches[1]),
			},
			Speed: Velocity{
				X: StringToInt(matches[2]),
				Y: StringToInt(matches[3]),
			},
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(2)
	}

	return lights
}
