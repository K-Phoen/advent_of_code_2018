package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Unit byte

func (unit Unit) String() string {
	return string(unit)
}

func readPolymer(fileName string) []Unit {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file '%s'\n", fileName)
		os.Exit(1)
	}
	defer file.Close()

	var polymer []Unit

	reader := bufio.NewReader(file)

	for {
		b, err := reader.ReadByte()

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}

			break
		}

		if b == '\n' {
			continue
		}

		polymer = append(polymer, Unit(b))
	}

	return polymer
}

func (unit Unit) TypeEquals(other Unit) bool {
	diff := int(unit) - int(other)

	// in the ASCII table, the difference between the same lower case and upper
	// case character is 32
	return diff == 32 || diff == -32
}

func (unit Unit) PolarityEquals(other Unit) bool {
	return unit == other
}

func ReducePolymer(polymer []Unit) []Unit {
	length := len(polymer)
	i := 0

	for i < length-1 {
		a := polymer[i]
		b := polymer[i+1]

		if a.TypeEquals(b) && !a.PolarityEquals(b) {
			polymer = append(polymer[:i], polymer[i+2:]...)

			if i != 0 {
				i--
			}

			length -= 2

			continue
		}

		i++
	}

	return polymer
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	polymer := readPolymer(os.Args[1])

	//fmt.Printf("Original polymer: %s\n", polymer)

	//fmt.Printf("Reduced polymer: %s\n", ReducePolymer(polymer))
	fmt.Printf("Reduced polymer's size: %d\n", len(ReducePolymer(polymer)))
}
