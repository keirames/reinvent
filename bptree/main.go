package main

import "main/bptree"

func main() {
	tree := bptree.New()
	tree.Insert(1)
	tree.Traversal()
	tree.Insert(3)
	tree.Traversal()
	tree.Insert(2)
	tree.Traversal()
	tree.Insert(4)
	tree.Traversal()
	tree.Insert(5)
	tree.Traversal()
	// tree.Insert(6)
	// tree.Traversal()
}
