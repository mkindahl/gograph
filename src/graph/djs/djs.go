package djs

import "fmt"

type Value interface{}

type Node struct {
	rank int
	value Value
	parent *Node
}

type DisjointSet map[interface{}]*Node

func New() *DisjointSet {
	self := make(DisjointSet)
	return &self
}

func (node *Node) String() string {
	return fmt.Sprint(node.value)
}

func (ds DisjointSet) MakeSet(v Value) {
	node := Node{rank: 0, value: v}
	node.parent = &node
	ds[v] = &node
}

func Link(x, y *Node) {
	if x.rank > y.rank {
		y.parent = x
	} else {
		x.parent = y
		if x.rank == y.rank {
			y.rank = y.rank + 1
		}
	}
}

func (ds DisjointSet) Find(v Value) *Node {
	node := ds[v]
	for node != node.parent {
		node = node.parent
	}
	return node.parent
}

func (ds DisjointSet) Union(x, y Value) {
	Link(ds.Find(x), ds.Find(y))
}
