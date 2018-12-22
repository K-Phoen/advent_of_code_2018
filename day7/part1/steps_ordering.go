package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	requirements := ReadRequirements(os.Args[1])
	graph := NewGraph()

	for _, requirement := range requirements {
		graph.AddEdge(requirement.task, requirement.comesBefore)
	}

	graph.OutputAsDot()

	stepsOrdering := graph.TopologicalSort()

	for _, node := range stepsOrdering {
		fmt.Printf("%c", node.Name)
	}
	fmt.Print("\n")
}
