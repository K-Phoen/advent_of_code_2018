package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Usage: %s FILE\n", os.Args[0])
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	frequency := 0
	for scanner.Scan() {
		text := scanner.Text()
		item, err := strconv.Atoi(text)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Incorrect input found: '%s' (could not be converted to a number)\n", text)
			os.Exit(2)
		}

		frequency += item
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(2)
	}

	fmt.Printf("Computed frequency: %d\n", frequency)
}
