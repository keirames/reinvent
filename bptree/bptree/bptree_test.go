package bptree

import (
	"testing"
)

func TestFindLeaf(t *testing.T) {
	tree := New()
	_, got := findLeaf(tree, 1)
	want := EmptyTreeErr
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestFindLeaf1(t *testing.T) {
	tree := New()
	tree.Root = makeNode()
	tree.Root.Keys[0] = 1
	tree.Root.NumKeys = 1
	// tree.Root.Pointers[0] =

	got, _ := findLeaf(tree, 1)
	want := tree.Root
	if got != want {
		t.Errorf("got wrong value")
	}
}
