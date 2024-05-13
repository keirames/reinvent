package bptree

import "fmt"

var (
	// default order of a tree
	order = 4
)

type Tree struct {
	Root *Node
}

type Record struct {
	Value []byte
}

type Node struct {
	// keys inside single node, capacity calc by order - 1
	Keys []int

	// num of keys have value inside node
	NumKeys int

	// is this node act as a leaf
	IsLeaf bool

	Pointers []any

	Parent *Node

	Next *Node
}

func New() *Tree {
	tree := new(Tree)
	return tree
}

func calcKeysInNode(treeOrder int) int {
	// b+ tree specific characteristic
	return order - 1
}

func makeNode() *Node {
	n := new(Node)
	n.Keys = make([]int, calcKeysInNode(order))
	n.IsLeaf = false
	n.Parent = nil
	n.Next = nil

	return n
}

func makeLeaf() *Node {
	leaf := makeNode()
	leaf.IsLeaf = true

	return leaf
}

func (t *Tree) Insert(key int, verbose bool) (*Node, error) {
	i := 0
	n := t.Root
	if n == nil {
		if verbose {
			fmt.Println("Empty tree.")
		}

		return nil, fmt.Errorf("Empty tree.")
	}

	for !n.IsLeaf {
		i = 0
		for i < n.NumKeys {
			// right bias
			if key >= n.Keys[i] {
				i++
			} else {
				break
			}
		}

		n, ok := n.Pointers[i].(*Node)
		if !ok {
			return nil, fmt.Errorf("Pointer is nil")
		} else {
			return n, nil
		}
	}

	return nil, fmt.Errorf("Not Found.")
}
