package main

import "main/bptree"

func main() {
	tree := bptree.New()
	tree.Insert(1)
	tree.Insert(3)
	tree.Insert(2)
	tree.Insert(4)
	// tree.Insert(5)
	tree.Traversal()
}
