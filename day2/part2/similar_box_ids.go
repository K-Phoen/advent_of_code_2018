package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func readBoxIds(fileName string) chan string {
	channel := make(chan string)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file '%s'\n", fileName)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)

	go func() {
		defer close(channel)

		for scanner.Scan() {
			channel <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			os.Exit(2)
		}
	}()

	return channel
}

func hasSingleCharDifference(first string, second string) bool {
	// inputs are expected to have same length
	if len(first) != len(second) {
		return false
	}

	diffsCount := 0

	for i := 0; i < len(first); i++ {
		if first[i] != second[i] {
			diffsCount += 1
		}

		if diffsCount >= 2 {
			return false
		}
	}

	return diffsCount == 1
}

func commonChars(first string, second string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(first); i++ {
		if first[i] == second[i] {
			buffer.WriteByte(first[i])
		}
	}

	return buffer.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	for firstBoxId := range readBoxIds(os.Args[1]) {
		for secondBoxId := range readBoxIds(os.Args[1]) {
			if hasSingleCharDifference(firstBoxId, secondBoxId) {
				fmt.Printf("%s and %s differ by a single letter\n", firstBoxId, secondBoxId)
				fmt.Printf("Common letters: %s\n", commonChars(firstBoxId, secondBoxId))
				os.Exit(0)
			}
		}
	}
}
