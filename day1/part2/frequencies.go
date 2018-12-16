package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func readChanges(fileName string) chan int {
	channel := make(chan int)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Usage: %s FILE\n", os.Args[0])
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)

	go func() {
		defer close(channel)

		for scanner.Scan() {
			text := scanner.Text()
			value, err := strconv.Atoi(text)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Incorrect input found: '%s' (could not be converted to a number)\n", text)
				os.Exit(2)
			}

			channel <- value
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			os.Exit(2)
		}
	}()

	return channel
}

func repeatChanges(fileName string) chan int {
	channel := make(chan int)

	go func() {
		defer close(channel)

		for {
			for change := range readChanges(fileName) {
				channel <- change
			}
		}

	}()

	return channel
}

func findDuplicateFrequency(changes chan int) (int, error) {
	frequency := 0
	seen := make(map[int]bool)
	seen[frequency] = true

	for change := range changes {
		//fmt.Printf("Read change: %d\n", change)
		frequency += change
		//fmt.Printf("New frequency is: %d\n", frequency)

		if _, exists := seen[frequency]; exists {
			return frequency, nil
		}

		seen[frequency] = true
	}

	return frequency, errors.New("No duplicate found")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	frequency, err := findDuplicateFrequency(repeatChanges(os.Args[1]))
	if err != nil {
		fmt.Printf("%s (current frequency: %d)\n", err.Error(), frequency)
		os.Exit(1)
	}

	fmt.Printf("Found duplicated frequency: %d\n", frequency)
}
