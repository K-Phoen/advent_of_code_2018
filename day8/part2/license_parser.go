package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Step uint

const ReadChildNodesCount Step = 0
const ReadMetadataCount Step = 1
const ReadMetadata Step = 2

type Node struct {
	Name byte

	Children      []*Node
	ChildrenCount int

	Metadata       []int
	MetadataLength int
}

type Tree struct {
	Root *Node
}

func (node *Node) AddChild(child *Node) {
	node.Children = append(node.Children, child)
}

func (node *Node) AddMetadata(metadata int) {
	node.Metadata = append(node.Metadata, metadata)
}

func (node *Node) MetadataSum() int {
	sum := 0

	for _, i := range node.Metadata {
		sum += i
	}

	return sum
}

func (node *Node) Value() int {
	if node.ChildrenCount == 0 {
		return node.MetadataSum()
	}

	value := 0

	for _, i := range node.Metadata {
		if i == 0 {
			continue
		}

		// entry won't exist in Children array
		if i > node.ChildrenCount {
			continue
		}

		value += node.Children[i-1].Value()
	}

	return value
}

func (tree *Tree) OutputAsDot() {
	fmt.Println("digraph {")

	var walk func(node *Node)
	walk = func(node *Node) {
		for _, child := range node.Children {
			fmt.Printf("\t%c -> %c\n", node.Name, child.Name)

			walk(child)
		}
	}

	walk(tree.Root)

	fmt.Println("}")
}

func StringToInt(input string) int {
	i, err := strconv.Atoi(input)

	if err != nil {
		panic("Unable to convert input to string")
	}

	return i
}

func ReadLicenseTree(fileName string) *Tree {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file '%s'\n", fileName)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	nodes := NewStack()
	step := ReadChildNodesCount
	tree := &Tree{}
	nextName := byte('A')

	for scanner.Scan() {
		value := StringToInt(scanner.Text())

		if step == ReadChildNodesCount {
			nodes.Push(&Node{
				Name:          nextName,
				ChildrenCount: value,
			})

			nextName++

			step = ReadMetadataCount
		} else if step == ReadMetadataCount {
			node := nodes.Peek().(*Node)
			node.MetadataLength = value

			// this node is a leaf, we can read its metdata
			if node.ChildrenCount == 0 {
				step = ReadMetadata
			} else {
				// this node is NOT a leaf, we have to read its children first
				step = ReadChildNodesCount
			}
		} else if step == ReadMetadata {
			node := nodes.Peek().(*Node)

			node.AddMetadata(value)

			// we finished reading the metadata for the current node
			if len(node.Metadata) == node.MetadataLength {
				nodes.Pop()

				// if there are no nodes in the stack, it means that we found the root
				if nodes.Len() == 0 {
					tree.Root = node
					break
				}

				parent := nodes.Peek().(*Node)
				parent.AddChild(node)

				// if the parent node had several children, we start reading the metadata for the remaining children
				if len(parent.Children) == parent.ChildrenCount {
					step = ReadMetadata
				} else {
					// otherwise we read new nodes
					step = ReadChildNodesCount
				}
			}
		} else {
			panic("Unknown step")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(2)
	}

	return tree
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	licenseTree := ReadLicenseTree(os.Args[1])

	fmt.Printf("Tree value: %d\n", licenseTree.Root.Value())
}
