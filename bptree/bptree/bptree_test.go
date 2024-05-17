package bptree

import (
	"fmt"
	"testing"
)

func TestFindKeyResponseEmptyTree(t *testing.T) {
	tree := New()

	err := tree.Find(1)
	if err != EmptyTreeErr {
		t.Errorf("test failed")
	}
}

func TestFindKeySingleNoInternalNode(t *testing.T) {
	tree := New()
	tree.Root = makeLeaf()

	cur := tree.Root
	cur.Keys[0] = 5
	cur.Keys[1] = 15
	cur.NumKeys = 2

	err := tree.Find(5)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(15)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(0)
	if err != NodeNotFound {
		t.Errorf("test failed, expected %q got nil %q", NodeNotFound, err)
	}

	tree = New()
	tree.Root = makeLeaf()

	cur = tree.Root
	cur.Keys[0] = 5
	cur.NumKeys = 1

	err = tree.Find(5)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(0)
	if err != NodeNotFound {
		t.Errorf("test failed, expected %q got nil %q", NodeNotFound, err)
	}
}

func TestFindKeyInComplexTree(t *testing.T) {
	tree := New()
	tree.Root = makeNode()

	cur := tree.Root
	cur.Keys[0] = 15
	cur.NumKeys = 1

	leaf1 := makeLeaf()
	leaf1.Keys[0] = 5
	leaf1.NumKeys = 1

	leaf2 := makeLeaf()
	leaf2.Keys[0] = 15
	leaf2.Keys[1] = 25
	leaf2.NumKeys = 2

	leaf1.Next = leaf2

	cur.Pointers[0] = leaf1
	cur.Pointers[1] = leaf2

	err := tree.Find(5)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(15)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(25)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(35)
	if err != NodeNotFound {
		t.Errorf("test find key failed, expected %q got %q", NodeNotFound, err)
	}

	tree = New()
	tree.Root = makeNode()

	cur = tree.Root
	cur.Keys[0] = 25
	cur.NumKeys = 1

	node1 := makeNode()
	node1.Keys[0] = 15
	node1.NumKeys = 1

	node2 := makeNode()
	node2.Keys[0] = 35
	node2.NumKeys = 1

	cur.Pointers[0] = node1
	cur.Pointers[1] = node2

	leaf1 = makeLeaf()
	leaf1.Keys[0] = 5
	leaf1.NumKeys = 1

	leaf2 = makeLeaf()
	leaf2.Keys[0] = 15
	leaf2.NumKeys = 1

	node1.Pointers[0] = leaf1
	node1.Pointers[1] = leaf2

	leaf3 := makeLeaf()
	leaf3.Keys[0] = 25
	leaf3.NumKeys = 1

	leaf4 := makeLeaf()
	leaf4.Keys[0] = 35
	leaf4.Keys[1] = 45
	leaf4.NumKeys = 2

	node2.Pointers[0] = leaf3
	node2.Pointers[1] = leaf4

	leaf1.Next = leaf2
	leaf2.Next = leaf3
	leaf3.Next = leaf4

	err = tree.Find(5)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(15)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(25)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(35)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(45)
	if err != nil {
		t.Errorf("test find key failed, expected nil got %q", err)
	}

	err = tree.Find(55)
	if err != NodeNotFound {
		t.Errorf("test find key failed, expected %q got %q", NodeNotFound, err)
	}

	err = tree.Find(0)
	if err != NodeNotFound {
		t.Errorf("test find key failed, expected %q got %q", NodeNotFound, err)
	}
}

func TestInsertIntoEmptyTree(t *testing.T) {
	tree := New()
	keyVal := 1
	n, _ := tree.Insert(keyVal)
	fmt.Println(n)
	if n.Keys[0] != keyVal {
		t.Errorf("test failed, expected key %q got %q", keyVal, n.Keys[0])
	}

	if n.NumKeys != 1 {
		t.Errorf("test failed, expected nums keys %q got %q", 1, n.NumKeys)
	}
}
