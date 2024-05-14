package bptree

import "fmt"

var (
	// default order of a tree
	order = 3
)

var (
	EmptyTreeErr = fmt.Errorf("Empty tree.")
	NodeNotFound = fmt.Errorf("Node not found.")
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

	Pointers []*Node

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

func findLeafValue(t *Tree, key int) (int, error) {
	n, err := findLeaf(t, key)
	if err != nil {
		return 0, err
	}

	for i := range n.NumKeys - 1 {
		if n.Keys[i] == key {
			return key, nil
		}
	}

	return 0, NodeNotFound
}

func findLeaf(t *Tree, key int) (*Node, error) {
	n := t.Root
	if n == nil {
		return nil, EmptyTreeErr
	}

	i := 0
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

		n, _ = n.Pointers[i].(*Node)

		if n == nil {
			break
		}
	}

	if n == nil {
		return nil, NodeNotFound
	}

	return n, nil
}

func (t *Tree) find(key int) error {
	cur := t.Root
	if cur == nil {
		return EmptyTreeErr
	}

	// loop until cur is leaf
	for !cur.IsLeaf {
		for i := range cur.NumKeys - 1 {
			if cur.Keys[i] >= key {
				// right bias
				cur = cur.Pointers[i+1]
			}
		}
	}

	for i := range cur.NumKeys - 1 {
		if cur.Keys[i] == key {
			return nil
		}
	}

	return NodeNotFound
}

func (t *Tree) Insert(key int) (*Node, error) {
	if t.Root == nil {
		t.Root = makeLeaf()
		t.Root.NumKeys++

		return t.Root, nil
	}

	return nil, fmt.Errorf("")
}
