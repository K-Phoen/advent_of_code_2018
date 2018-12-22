package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Requirement struct {
	task        byte
	comesBefore byte
}

func (r *Requirement) String() string {
	return fmt.Sprintf("Task %c comes before %c", r.task, r.comesBefore)
}

func extractTask(input string) byte {
	step := strings.Index(input, "Step")

	return input[step+5]
}

func extractNextTask(input string) byte {
	step := strings.Index(input, "step")

	return input[step+5]
}

func ReadRequirements(fileName string) []*Requirement {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file '%s'\n", fileName)
		os.Exit(1)
	}
	defer file.Close()

	var requirements []*Requirement

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		requirements = append(requirements, &Requirement{
			task:        extractTask(line),
			comesBefore: extractNextTask(line),
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(2)
	}

	return requirements
}
