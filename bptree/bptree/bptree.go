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
	EmptyTreeErr      = fmt.Errorf("Empty tree.")
	TreeIsNotEmptyErr = fmt.Errorf("Tree is not empty")
	NodeNotFound      = fmt.Errorf("Node not found.")
	DupKeyErr         = fmt.Errorf("Key already exists.")
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
	n.Keys = make([]int, 0)
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

func (t *Tree) makeNewTree(key int) error {
	if t.Root != nil {
		return TreeIsNotEmptyErr
	}

	t.Root = makeLeaf()
	t.Root.Keys = append(t.Root.Keys, key)
	t.Root.NumKeys += 1

	return nil
}

func (t *Tree) Insert(key int) error {
	// empty tree, create new tree
	if t.Root == nil {
		return t.makeNewTree(key)
	}

	// TODO: not accept duplicate
	// err := t.Find(key)
	// if err != nil {
	// 	return DupKeyErr
	// }

	// insert into leaf node
	leaf, err := t.FindLeaf(key)
	if err != nil {
		return err
	}

	// leaf node still has space
	if leaf.NumKeys < order-1 {
		insertIntoLeaf(leaf, key)
		return nil
	}

	t.insertIntoLeafAfterSplitting(leaf, key)

	return nil
}

func insertIntoLeaf(leaf *Node, key int) {
	// tempArr := make([]int, leaf.NumKeys+1)
	// fmt.Println(tempArr)
	// copy(tempArr, leaf.Keys)
	// fmt.Println(tempArr)
	leaf.NumKeys++
	leaf.Keys = InsertIntoSortedArray(leaf.Keys, key)
}

func (t *Tree) insertIntoLeafAfterSplitting(leaf *Node, key int) {
	// splitting
	tempArr := make([]int, len(leaf.Keys))
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

	fmt.Println(leaf.Keys, newLeaf.Keys, dupKey)

	// lift up to parent
	t.insertIntoParent(leaf, newLeaf, dupKey)
}

type SimpleQueue struct {
	q []*Node
}

func NewSimpleQueue() *SimpleQueue {
	return &SimpleQueue{
		q: []*Node{},
	}
}

func (sq *SimpleQueue) Push(n *Node) {
	sq.q = append(sq.q, n)
}

func (sq *SimpleQueue) Pop() *Node {
	popNode := sq.q[0]
	newQ := []*Node{}

	for i, n := range sq.q {
		if i == 0 {
			continue
		}

		newQ = append(newQ, n)
	}

	sq.q = newQ

	return popNode
}

func (t *Tree) Traversal() {
	if t.Root == nil {
		fmt.Printf("empty tree!")
		return
	}

	sq := NewSimpleQueue()
	sq.Push(t.Root)
	for len(sq.q) != 0 {
		l := len(sq.q)
		for range l {
			n := sq.Pop()
			fmt.Println(n.Keys)

			for _, p := range n.Pointers {
				if p == nil {
					continue
				}

				sq.Push(p)
			}
		}
		fmt.Println("---")
	}
}

func (t *Tree) insertIntoParent(leaf *Node, newLeaf *Node, key int) {
	// * case: no parent
	if leaf.Parent == nil {
		parent := makeNode()
		leaf.Parent = parent
		newLeaf.Parent = parent

		parent.NumKeys = 1
		parent.Keys = append(parent.Keys, key)
		parent.Pointers[0] = leaf
		parent.Pointers[1] = newLeaf

		t.Root = parent

		return
	}

	// * case: has parent - parent doesn't need split
	parent := leaf.Parent
	if parent.NumKeys < order-1 {
		idx := -1
		for i, n := range parent.Keys {
			if n > key {
				idx = i
				break
			}
		}

		// insert into first
		if idx == 0 {
			newKeys := []int{}
			newKeys = append(newKeys, key)

			for _, n := range parent.Keys {
				newKeys = append(newKeys, n)
			}

			newPointers := []*Node{}
			newPointers = append(newPointers, leaf)
			newPointers = append(newPointers, newLeaf)

			for i, p := range parent.Pointers {
				if i == 0 {
					continue
				}

				newPointers = append(newPointers, p)
			}

			parent.Keys = newKeys
			parent.Pointers = newPointers

			return
		}

		// insert to last
		if idx == -1 {
			parent.Keys = append(parent.Keys, key)
			parent.Pointers = append(parent.Pointers, newLeaf)

			return
		}

		// insert into middle
		newKeys := InsertIntoSortedArray(parent.Keys, key)
		newPointers := []*Node{}

		for i, p := range parent.Pointers {
			newPointers = append(newPointers, p)

			if i == idx {
				newPointers = append(newPointers, newLeaf)
			}
		}

		parent.Keys = newKeys
		parent.Pointers = newPointers

		return
	}

	fmt.Println("has parent parent need split")
	// * case: has parent - parent need split
	// TODO: Rearrange keys is easy task, rearrange pointer arr is hard
	// TODO: Refactor rearrange key in 1 time on every case
	// numKeysOverflow := InsertIntoSortedArray(parent.Keys, key)
	// mid := len(numKeysOverflow) / 2
	newKeys := InsertIntoSortedArray(parent.Keys, key)
	idx := -1
	for i, n := range parent.Keys {
		if n > key {
			idx = i
			break
		}
	}

	newPointers := []*Node{}

	//* insert into first
	if idx == 0 {
		newPointers = append(newPointers, leaf)
		newPointers = append(newPointers, newLeaf)

		for i, p := range parent.Pointers {
			if i == 0 {
				continue
			}

			newPointers = append(newPointers, p)
		}
	}

	//* insert into last
	if idx == -1 {
		for _, p := range parent.Pointers {
			newPointers = append(newPointers, p)
		}

		newPointers = append(newPointers, newLeaf)
	}

	//* insert into middle
	if idx != 0 || idx != -1 {
		for i, p := range parent.Pointers {
			newPointers = append(newPointers, p)

			if i == idx {
				newPointers = append(newPointers, newLeaf)
			}
		}
	}

	//* splitting
	// TODO: how to clean reference up ?
	newLeftNode := makeNode()
	newRightNode := makeNode()
	leftNodeKeys := []int{}
	rightNodeKeys := []int{}
	leftNodePointers := []*Node{}
	rightNodePointers := []*Node{}

	for i := range len(newKeys) {
		if i == idx {
			continue
		}

		if i < idx {
			leftNodeKeys = append(leftNodeKeys, newKeys[i])
			continue
		}

		rightNodeKeys = append(rightNodeKeys, newKeys[i])
	}

	newLeftNode.Keys = leftNodeKeys
	newLeftNode.NumKeys = len(leftNodeKeys)

	newRightNode.Keys = rightNodeKeys
	newRightNode.NumKeys = len(rightNodeKeys)

	newLeftNode.Next = newRightNode

	newLeftNode.Parent = parent.Parent
	newRightNode.Parent = parent.Parent

	for i, p := range newPointers {
		// include idx
		if i <= idx {
			leftNodePointers = append(leftNodePointers, p)
			continue
		}

		rightNodePointers = append(rightNodePointers, p)
	}

	// t.insertIntoParent(newLeftNode, newRightNode, newKeys[idx])
	// t.insertIntoParentFromInternalNode(newLeftNode, newRightNode, newKeys[idx])
}

func (t *Tree) insertIntoParentFromInternalNode(
	leftNode *Node,
	rightNode *Node,
	liftKey int,
) {
	//* case: no parent
	if leftNode.Parent == nil {
		newParent := makeNode()
		newParent.Keys[0] = liftKey
		newParent.NumKeys += 1

		leftNode.Parent = newParent
		rightNode.Parent = newParent

		t.Root = newParent

		return
	}

	//* case: parent has space
	if leftNode.Parent.NumKeys < order-1 {
		parent := leftNode.Parent

		idx := -1
		for i, n := range parent.Keys {
			if n > liftKey {
				idx = i
				break
			}
		}

		parent.Keys = InsertIntoSortedArray(parent.Keys, liftKey)
		parent.NumKeys += 1

		// insert into first
		if idx == 0 {
			newPointers := []*Node{}
			newPointers = append(newPointers, leftNode)
			newPointers = append(newPointers, rightNode)

			for i, p := range parent.Pointers {
				if i == 0 {
					continue
				}

				newPointers = append(newPointers, p)
			}

			parent.Pointers = newPointers

			return
		}

		// insert into last
		if idx == -1 {
			parent.Pointers[parent.NumKeys] = rightNode
			parent.Pointers[parent.NumKeys-1] = leftNode

			return
		}

		// insert into middle
		newPointers := []*Node{}
		for i, p := range parent.Pointers {
			if i == idx {
				newPointers = append(newPointers, leftNode)
				newPointers = append(newPointers, rightNode)
			}

			newPointers = append(newPointers, p)
		}

		parent.Pointers = newPointers

		return
	}

	//* case: parent need split & lift key up
	idx := -1
	parent := leftNode.Parent
	for i, n := range parent.Keys {
		if n > liftKey {
			idx = i
			break
		}
	}

	if idx == 0 {
		// ! what the fuck i'm so confusing
	}
}

func InsertIntoSortedArray(arr []int, n int) []int {
	arr = append(arr, n)
	sort.Ints(arr)

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
