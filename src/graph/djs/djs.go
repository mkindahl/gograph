// Implementation of union-find data structure using the disjoint-set
// forest with union by rank and path compression. Based on
// description of disjoin-set forests in "Introduction to Algorithms"
// by Cormen et.al.
package djs

// A node in the disjoint-set forest.
type Node struct {
	rank int
	parent *Node
}

// A disjoint-set structure, keeping information about all the members
// of the sets and the sets.
type DisjointSet struct {
	nodes map[interface{}]*Node
}

// Create a new disjoint-set 
func New() *DisjointSet {
	return &DisjointSet{nodes: make(map[interface{}]*Node)}
}

// Create a set for 'value' and add it to the disjoint-set structure.
func (ds *DisjointSet) MakeSet(value interface{}) {
	node := Node{rank: 0}
	node.parent = &node
	ds.nodes[value] = &node
}

// Find the representative node for the set that 'value' is member
// of. If the representative of two nodes is identical, they are in
// the same set. If they are different, they are in different sets.
func (ds *DisjointSet) Find(value interface{}) *Node {
	node := ds.nodes[value]
	for node != node.parent {
		node = node.parent
	}
	return node.parent
}

// Merge the two sets that 'x' and 'y' are members of.
func (ds *DisjointSet) Union(x, y interface{}) {
	nx := ds.Find(x)
	ny := ds.Find(y)
	if nx.rank > ny.rank {
		ny.parent = nx
	} else {
		nx.parent = ny
		if nx.rank == ny.rank {
			ny.rank = ny.rank + 1
		}
	}
}
