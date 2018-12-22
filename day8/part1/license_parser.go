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
	Children      []*Node
	ChildrenCount int

	Metadata       []int
	MetadataLength int
}

type Tree struct {
	Root        *Node
	MetadataSum int
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
	tree := &Tree{
		MetadataSum: 0,
	}

	for scanner.Scan() {
		value := StringToInt(scanner.Text())

		if step == ReadChildNodesCount {
			nodes.Push(&Node{
				ChildrenCount: value,
			})

			step = ReadMetadataCount
		} else if step == ReadMetadataCount {
			node := nodes.Peek().(*Node)
			node.MetadataLength = value

			// this node is a leaf, we can read its metdata
			if node.ChildrenCount == 0 {
				step = ReadMetadata

				parent := nodes.Peek().(*Node)
				parent.Children = append(parent.Children, node)
			} else {
				// this node is NOT a leaf, we have to read its children first
				step = ReadChildNodesCount
			}
		} else if step == ReadMetadata {
			node := nodes.Peek().(*Node)

			tree.MetadataSum += value
			node.Metadata = append(node.Metadata, value)

			// we finished reading the metadata for the current node
			if len(node.Metadata) == node.MetadataLength {
				nodes.Pop()

				// if there are no nodes in the stack, it means that we found the root
				if nodes.Len() == 0 {
					tree.Root = node
					break
				}

				parent := nodes.Peek().(*Node)
				parent.Children = append(parent.Children, node)

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

	fmt.Printf("Sum of all metadata entries: %d\n", licenseTree.MetadataSum)
}
