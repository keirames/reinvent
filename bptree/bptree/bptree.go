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
	n.Pointers = make([]*Node, 0)
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
	fmt.Println("what", cur)
	for !cur.IsLeaf {
		flag := false
		for i := range cur.NumKeys {
			fmt.Println("range loop", cur.NumKeys, i)
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
			fmt.Println("xxx", cur, key)
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
		fmt.Println("empty tree, make new tree")
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
	fmt.Println("valid leaf for key: ", key, leaf.Keys)

	// leaf node still has space
	if leaf.NumKeys < order-1 {
		fmt.Println("still has space")
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
	var oldKeys []int
	if leaf.Next != nil {
		oldKeys = leaf.Next.Keys
	}
	fmt.Println("old keys", oldKeys)
	fmt.Println("old leaf", leaf)
	fmt.Println("next leaf", leaf.Next)
	leaf.NumKeys++
	leaf.Keys = InsertIntoSortedArray(leaf.Keys, key)
	fmt.Println("new leaf", leaf, leaf.Next)
	fmt.Println("next leaf", leaf.Next)

	fmt.Println(oldKeys)
}

func (t *Tree) insertIntoLeafAfterSplitting(leaf *Node, key int) {
	fmt.Println("split leaf")
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

	// parent
	newLeaf.Parent = leaf.Parent

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
			fmt.Println("node info", n)
			fmt.Println(n.Keys, n.Pointers, n.Parent)
			for _, p := range n.Pointers {
				if p == nil {
					// fmt.Print("nil ")
					continue
				}
				// fmt.Print(p.Keys)
			}

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

func becomingNewInternalNodeAsParent(t *Tree, l *Node, r *Node, key int) {
	n := makeNode()
	n.Keys = append(n.Keys, key)
	n.Pointers = append(n.Pointers, l)
	n.Pointers = append(n.Pointers, r)
	n.NumKeys = 1

	// update parent
	for _, p := range n.Pointers {
		p.Parent = n
		fmt.Println("pointer of parent", p)
	}

	t.Root = n
}

func mergeIntoInternalNodeAsParent(curNode *Node, l *Node, r *Node, key int) {
	curNode.NumKeys++
	curNode.Keys = InsertIntoSortedArray(curNode.Keys, key)

	idx := -1
	for i, n := range curNode.Keys {
		if n == key {
			idx = i
			break
		}
	}

	if idx == -1 {
		panic("It's not gonna happen!")
	}

	newPointers := []*Node{}

	// insert into first
	if idx == 0 {
		for i, p := range curNode.Pointers {
			if i == 0 {
				// skip first old pointer, insert into 2 new pointer
				newPointers = append(newPointers, l)
				newPointers = append(newPointers, r)
				continue
			}

			// insert rest
			newPointers = append(newPointers, p)
		}
	} else if idx == curNode.NumKeys-1 {
		// insert into last
		for i, p := range curNode.Pointers {
			if i == len(curNode.Pointers)-1 {
				// skip last old pointer, insert into 2 new pointer
				newPointers = append(newPointers, l)
				newPointers = append(newPointers, r)
				continue
			}

			newPointers = append(newPointers, p)
		}
	} else {
		// insert into middle
		for i, p := range curNode.Pointers {
			if i == idx {
				// skip this pointer, insert into 2 new pointer
				newPointers = append(newPointers, l)
				newPointers = append(newPointers, r)
				continue
			}

			newPointers = append(newPointers, p)
		}
	}

	curNode.Pointers = newPointers
	// reassign parent
	for _, p := range curNode.Pointers {
		p.Parent = curNode
	}

	fmt.Println("curnode", curNode)
}

func hasSpace(n int) bool {
	fmt.Println("has space", n, order)
	return n < order-1
}

func splitOverflowInternalNode(t *Tree, curNode *Node, l *Node, r *Node, key int) {
	curNode.NumKeys++
	curNode.Keys = InsertIntoSortedArray(curNode.Keys, key)

	idx := -1
	for i, n := range curNode.Keys {
		if n == key {
			idx = i
			break
		}
	}

	if idx == -1 {
		panic("It's not gonna happen!")
	}

	newPointers := []*Node{}

	// insert into first
	if idx == 0 {
		for i, p := range curNode.Pointers {
			if i == 0 {
				// skip first old pointer, insert into 2 new pointer
				newPointers = append(newPointers, l)
				newPointers = append(newPointers, r)
				continue
			}

			// insert rest
			newPointers = append(newPointers, p)
		}
	} else if idx == curNode.NumKeys-1 {
		// insert into last
		for i, p := range curNode.Pointers {
			if i == len(curNode.Pointers)-1 {
				// skip last old pointer, insert into 2 new pointer
				newPointers = append(newPointers, l)
				newPointers = append(newPointers, r)
				continue
			}

			newPointers = append(newPointers, p)
		}
	} else {
		// insert into middle
		for i, p := range curNode.Pointers {
			if i == idx {
				// skip this pointer, insert into 2 new pointer
				newPointers = append(newPointers, l)
				newPointers = append(newPointers, r)
				continue
			}

			newPointers = append(newPointers, p)
		}
	}

	// destroy curNode
	curNode.Pointers = []*Node{}

	// create 2 new node
	leftNode := makeNode()
	rightNode := makeNode()

	midIdx := len(curNode.Keys) / 2
	liftKey := curNode.Keys[midIdx]

	for i := range curNode.NumKeys {
		// mid key already lift up
		if i == midIdx {
			continue
		}

		// left node
		if i < midIdx {
			leftNode.Keys = append(leftNode.Keys, curNode.Keys[i])
			leftNode.NumKeys++
			continue
		}

		// right node
		rightNode.Keys = append(rightNode.Keys, curNode.Keys[i])
		rightNode.NumKeys++
	}

	for i, p := range newPointers {
		// left node
		if i <= midIdx {
			leftNode.Pointers = append(leftNode.Pointers, p)
			continue
		}

		// right node
		rightNode.Pointers = append(rightNode.Pointers, p)
	}

	// reassign child's parent
	for _, p := range leftNode.Pointers {
		p.Parent = leftNode
	}
	for _, p := range rightNode.Pointers {
		p.Parent = rightNode
	}

	// parent
	leftNode.Parent = curNode.Parent
	rightNode.Parent = curNode.Parent
	fmt.Println("parent xxx", curNode.Parent)

	t.insertIntoParent(leftNode, rightNode, liftKey)
}

func (t *Tree) insertIntoParent(leftNode *Node, rightNode *Node, key int) {
	// left node & right node 's parents are the same
	curNode := leftNode.Parent

	// * case: no parent
	if curNode == nil {
		fmt.Println("no parent, create parent", leftNode, rightNode)
		becomingNewInternalNodeAsParent(t, leftNode, rightNode, key)
		return
	}

	// * case: has parent - parent doesn't need split
	if hasSpace(curNode.NumKeys) {
		fmt.Println("has parent, parent doesn't need split", leftNode, rightNode)
		mergeIntoInternalNodeAsParent(curNode, leftNode, rightNode, key)
	} else {
		fmt.Println("has parent parent need split", curNode, key)
		// * case: has parent - parent need split
		splitOverflowInternalNode(t, curNode, leftNode, rightNode, key)
	}

	// fmt.Println("has parent parent need split")
	// // * case: has parent - parent need split
	// splitOverflowInternalNode(t, curNode, leftNode, rightNode, key)
	// TODO: Rearrange keys is easy task, rearrange pointer arr is hard
	// TODO: Refactor rearrange key in 1 time on every case
	// numKeysOverflow := InsertIntoSortedArray(parent.Keys, key)
	// mid := len(numKeysOverflow) / 2
	// newKeys := InsertIntoSortedArray(curNode.Keys, key)
	// idx := -1
	// for i, n := range curNode.Keys {
	// 	if n > key {
	// 		idx = i
	// 		break
	// 	}
	// }

	// newPointers := []*Node{}

	// //* insert into first
	// if idx == 0 {
	// 	newPointers = append(newPointers, leftNode)
	// 	newPointers = append(newPointers, rightNode)

	// 	for i, p := range curNode.Pointers {
	// 		if i == 0 {
	// 			continue
	// 		}

	// 		newPointers = append(newPointers, p)
	// 	}
	// }

	// //* insert into last
	// if idx == -1 {
	// 	fmt.Println("has parent, parent need split, insert into last")
	// 	fmt.Println("key", newKeys)
	// 	for _, p := range curNode.Pointers {
	// 		newPointers = append(newPointers, p)
	// 	}

	// 	fmt.Println("pointers", curNode.Pointers)
	// 	newPointers = append(newPointers, rightNode)
	// 	fmt.Println("new pointers", newPointers)
	// 	fmt.Println("---explain new pointers---")
	// 	for _, p := range newPointers {
	// 		fmt.Println(p.Keys)
	// 	}
	// 	fmt.Println("---end---")
	// }

	// //* insert into middle
	// if idx != 0 && idx != -1 {
	// 	for i, p := range curNode.Pointers {
	// 		newPointers = append(newPointers, p)

	// 		if i == idx {
	// 			newPointers = append(newPointers, rightNode)
	// 		}
	// 	}
	// }

	// midKeyIdx := (curNode.NumKeys + 1) / 2
	// fmt.Println("internal node", curNode.Keys, curNode.NumKeys)
	// fmt.Println("new internal node", newKeys, curNode.NumKeys+1)
	// fmt.Println("mid key idx is", midKeyIdx)
	// fmt.Println("mid key val is", newKeys[midKeyIdx])

	// //* splitting
	// // TODO: how to clean reference up ?
	// newLeftNode := makeNode()
	// newRightNode := makeNode()
	// leftNodeKeys := []int{}
	// rightNodeKeys := []int{}
	// leftNodePointers := []*Node{}
	// rightNodePointers := []*Node{}

	// for i := range len(newKeys) {
	// 	if i == midKeyIdx {
	// 		continue
	// 	}

	// 	if i < midKeyIdx {
	// 		leftNodeKeys = append(leftNodeKeys, newKeys[i])
	// 		continue
	// 	}

	// 	rightNodeKeys = append(rightNodeKeys, newKeys[i])
	// }

	// newLeftNode.Keys = leftNodeKeys
	// newLeftNode.NumKeys = len(leftNodeKeys)

	// newRightNode.Keys = rightNodeKeys
	// newRightNode.NumKeys = len(rightNodeKeys)

	// // newLeftNode.Next = newRightNode

	// // newLeftNode.Parent = curNode.Parent
	// // newRightNode.Parent = curNode.Parent

	// for i, p := range newPointers {
	// 	// include midKeyIdx
	// 	if i <= midKeyIdx {
	// 		leftNodePointers = append(leftNodePointers, p)
	// 		continue
	// 	}

	// 	rightNodePointers = append(rightNodePointers, p)
	// }
	// newLeftNode.Pointers = leftNodePointers
	// newRightNode.Pointers = rightNodePointers

	// fmt.Println("new left node", newLeftNode.Keys, newLeftNode.Pointers)
	// fmt.Println(newPointers)
	// fmt.Println("new right node", newRightNode.Keys, newRightNode.Pointers)

	// // change parent of child node
	// for _, p := range newRightNode.Pointers {
	// 	p.Parent = newRightNode
	// }
	// for _, p := range newLeftNode.Pointers {
	// 	p.Parent = newLeftNode
	// }

	// t.insertIntoParent(newLeftNode, newRightNode, newKeys[midKeyIdx])
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
	newArr := make([]int, 0)
	newArr = append(newArr, arr...)
	newArr = append(newArr, n)
	sort.Ints(newArr)

	// arr = append(arr, n)
	// sort.Ints(arr)

	// return arr
	return newArr
}

func IsDup(arr []int, n int) bool {
	for _, num := range arr {
		if num == n {
			return true
		}
	}

	return false
}
