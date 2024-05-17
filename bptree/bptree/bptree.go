package bptree

import (
	"fmt"
	"sort"
)

var (
	// default order of a tree
	order = 3
)

var (
	EmptyTreeErr = fmt.Errorf("Empty tree.")
	NodeNotFound = fmt.Errorf("Node not found.")
	DupKeyErr    = fmt.Errorf("Key already exists.")
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
	n.Pointers = make([]*Node, order)
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

// func findLeafValue(t *Tree, key int) (int, error) {
// 	n, err := findLeaf(t, key)
// 	if err != nil {
// 		return 0, err
// 	}

// 	for i := range n.NumKeys - 1 {
// 		if n.Keys[i] == key {
// 			return key, nil
// 		}
// 	}

// 	return 0, NodeNotFound
// }

func (t *Tree) FindLeaf(key int) (*Node, error) {
	cur := t.Root
	if cur == nil {
		return nil, EmptyTreeErr
	}

	// loop until cur is leaf
	for !cur.IsLeaf {
		flag := false
		for i := range cur.NumKeys {
			if cur.Keys[i] == key {
				// right bias
				cur = cur.Pointers[i+1]
				flag = true
				break
			}

			if cur.Keys[i] > key {
				cur = cur.Pointers[i]
				flag = true
				break
			}
		}

		if !flag {
			cur = cur.Pointers[cur.NumKeys]
		}
	}

	return cur, nil
}

func (t *Tree) Find(key int) error {
	cur := t.Root
	if cur == nil {
		return EmptyTreeErr
	}

	// loop until cur is leaf
	for !cur.IsLeaf {
		flag := false
		for i := range cur.NumKeys {
			if cur.Keys[i] == key {
				// right bias
				cur = cur.Pointers[i+1]
				flag = true
				break
			}

			if cur.Keys[i] > key {
				cur = cur.Pointers[i]
				flag = true
				break
			}
		}

		if !flag {
			cur = cur.Pointers[cur.NumKeys]
		}
	}

	for i := range cur.NumKeys {
		if cur.Keys[i] == key {
			return nil
		}
	}

	return NodeNotFound
}

func (t *Tree) Insert(key int) (*Node, error) {
	// empty tree, create new tree
	if t.Root == nil {
		t.Root = makeLeaf()
		t.Root.Keys[0] = key
		t.Root.NumKeys++

		return t.Root, nil
	}

	// insert into leaf node
	leaf, err := t.FindLeaf(key)
	if err != nil {
		panic(err)
	}

	// leaf node still has space
	if leaf.NumKeys < order-1 {
		// TODO: is there a way to insert into arr with less code ?
		tempArr := make([]int, leaf.NumKeys)
		copy(tempArr, leaf.Keys)

		if IsDup(leaf.Keys, key) {
			return nil, DupKeyErr
		}

		leaf.NumKeys++
		leaf.Keys = InsertIntoSortedArray(leaf.Keys, key)
	} else {
		// splitting
		tempArr := make([]int, leaf.NumKeys)
		copy(tempArr, leaf.Keys)
		tempArr = InsertIntoSortedArray(tempArr, key)

		mid := len(tempArr) / 2
		dupKey := tempArr[mid]

		leaf.Keys = tempArr[0:mid]
		leaf.NumKeys = len(leaf.Keys)
		prevNext := leaf.Next

		newLeaf := makeLeaf()
		newLeaf.Keys = tempArr[mid:]
		newLeaf.NumKeys = len(newLeaf.Keys)

		// linking
		leaf.Next = newLeaf
		newLeaf.Next = prevNext

		// lift up to parent
		// * case: no parent
		if leaf.Parent == nil {
			parent := makeNode()
			leaf.Parent = parent
			newLeaf.Parent = parent

			parent.NumKeys = 1
			parent.Keys[0] = dupKey
			parent.Pointers[0] = leaf
			parent.Pointers[1] = newLeaf

			return nil, fmt.Errorf("")
		}

		// * case: has parent - parent doesn't need split

		// * case: has parent - parent need split
	}

	return nil, fmt.Errorf("")
}

func InsertIntoSortedArray(arr []int, n int) []int {
	i := sort.SearchInts(arr, n)
	arr = append(arr, 0)
	copy(arr[i+1:], arr[i:])
	arr[i] = n

	return arr
}

func IsDup(arr []int, n int) bool {
	for _, num := range arr {
		if num == n {
			return true
		}
	}

	return false
}
