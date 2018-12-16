package main

import (
	"bufio"
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

func hasNSimilarLetters(text string, n int) bool {
	lettersCount := make(map[byte]int)

	for i := 0; i < len(text); i++ {
		lettersCount[text[i]] += 1
	}

	for _, count := range lettersCount {
		if count == n {
			return true
		}
	}

	return false
}

func countSimilarIds(boxIds chan string) (int, int) {
	twoLetters := 0
	threeLetters := 0

	for boxId := range readBoxIds(os.Args[1]) {
		if hasNSimilarLetters(boxId, 2) {
			twoLetters += 1
		}

		if hasNSimilarLetters(boxId, 3) {
			threeLetters += 1
		}
	}

	return twoLetters, threeLetters
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	twoLetters, threeLetters := countSimilarIds(readBoxIds(os.Args[1]))

	fmt.Printf("Two letters: %d ; three letters: %d ; checksum: %d", twoLetters, threeLetters, twoLetters*threeLetters)
}
