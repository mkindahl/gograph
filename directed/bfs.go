// Copyright (c) 2013, Ryan Marcus. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "iter"

// bfsEntry represents a vertex in the BFS queue along with its parent.
type bfsEntry[V comparable] struct {
	parent, vertex V
}

// BreadthFirstFrom returns an iterator that yields TraversalEvent
// values for a breadth-first traversal starting from the given vertex.
// Only vertices reachable from the start vertex will be visited.
func (graph *Graph[V]) BreadthFirstFrom(vertex V) iter.Seq[TraversalEvent[V]] {
	return func(yield func(TraversalEvent[V]) bool) {
		seen := make(map[V]uint8)
		graph.bfsIter(yield, seen, vertex)
	}
}

// BreadthFirst returns an iterator that yields TraversalEvent values
// for a breadth-first traversal of the entire graph.
func (graph *Graph[V]) BreadthFirst() iter.Seq[TraversalEvent[V]] {
	return func(yield func(TraversalEvent[V]) bool) {
		seen := make(map[V]uint8)
		for vertex := range graph.edges {
			if seen[vertex] != WHITE {
				continue
			}
			if !graph.bfsIter(yield, seen, vertex) {
				return
			}
		}
	}
}

// bfsIter performs a BFS from a single vertex, yielding events via
// yield. Returns false if yield requested stop.
func (graph *Graph[V]) bfsIter(yield func(TraversalEvent[V]) bool, seen map[V]uint8, vertex V) bool {
	if seen[vertex] != WHITE {
		return true
	}

	var zero V
	seen[vertex] = GREY
	if !yield(TraversalEvent[V]{Kind: Discover, Parent: zero, Vertex: vertex}) {
		return false
	}

	queue := []bfsEntry[V]{{parent: zero, vertex: vertex}}

	for len(queue) > 0 {
		entry := queue[0]
		queue = queue[1:]

		for _, child := range graph.OutEdges(entry.vertex) {
			if seen[child] == WHITE {
				seen[child] = GREY
				if !yield(TraversalEvent[V]{Kind: Discover, Parent: entry.vertex, Vertex: child}) {
					return false
				}
				queue = append(queue, bfsEntry[V]{parent: entry.vertex, vertex: child})
			}
		}

		seen[entry.vertex] = BLACK
		if !yield(TraversalEvent[V]{Kind: Finish, Parent: entry.parent, Vertex: entry.vertex}) {
			return false
		}
	}

	return true
}
