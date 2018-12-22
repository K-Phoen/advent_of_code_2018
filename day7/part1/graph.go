package main

import (
	"container/heap"
	"fmt"
)

type Node struct {
	Name byte

	Children map[byte]*Node
	Degree   int
}

type Graph struct {
	Nodes map[byte]*Node
}

func (node *Node) String() string {
	return string(node.Name)
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[byte]*Node),
	}
}

func (g *Graph) AddEdge(from byte, to byte) {
	g.ensureNodeExists(from)
	g.ensureNodeExists(to)

	fromNode := g.Node(from)

	if _, exists := fromNode.Children[to]; exists {
		return
	}

	toNode := g.Node(to)

	fromNode.Children[to] = toNode
	toNode.Degree++
}

func (g *Graph) RemoveEdge(from byte, to byte) {
	fromNode := g.Node(from)
	toNode := g.Node(to)

	// one of the nodes does not exist
	if fromNode == nil || toNode == nil {
		return
	}

	// edge does not exist
	if _, exists := fromNode.Children[to]; !exists {
		return
	}

	delete(fromNode.Children, to)
	toNode.Degree--
}

func (g *Graph) Node(name byte) *Node {
	return g.Nodes[name]
}

func (g *Graph) ensureNodeExists(name byte) {
	if _, exists := g.Nodes[name]; !exists {
		g.Nodes[name] = &Node{
			Name:     name,
			Children: make(map[byte]*Node),
			Degree:   0,
		}
	}
}

// https://en.wikipedia.org/wiki/Topological_sorting#Kahn's_algorithm
func (g *Graph) TopologicalSort() []*Node {
	var sortedNodes []*Node
	entryNodes := make(PriorityQueue, 0)
	heap.Init(&entryNodes)

	// first, find the nodes with no incoming edges
	for _, node := range g.Nodes {
		if node.Degree == 0 {
			heap.Push(&entryNodes, &Item{value: node.Name})
		}
	}

	/*
		while S is non-empty do
		    remove a node n from S
		    add n to tail of L
		    for each node m with an edge e from n to m do
		        remove edge e from the graph
		        if m has no other incoming edges then
		            insert m into S
	*/

	for entryNodes.Len() > 0 {
		item := heap.Pop(&entryNodes).(*Item)
		node := g.Node(item.value)

		sortedNodes = append(sortedNodes, node)

		for _, child := range node.Children {
			g.RemoveEdge(node.Name, child.Name)

			if child.Degree == 0 {
				heap.Push(&entryNodes, &Item{value: child.Name})
			}
		}
	}

	return sortedNodes
}

func (g *Graph) OutputAsDot() {
	fmt.Println("digraph {")

	for name := range g.Nodes {
		node := g.Node(name)

		for _, child := range node.Children {
			fmt.Printf("\t%c -> %c\n", name, child.Name)
		}
	}

	fmt.Println("}")
}
