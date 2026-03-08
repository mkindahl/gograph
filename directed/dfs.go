// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import (
	"container/list"
	"iter"
)

const (
	WHITE = iota // Undiscovered
	GREY         // Discovered, but not finalized
	BLACK        // Finalized
)

// dfsFrame represents a single stack frame for the iterative DFS.
type dfsFrame[V comparable] struct {
	parent, vertex V
	nextChild      *list.Element
}

// DepthFirst returns an iterator that yields TraversalEvent values
// for a depth-first traversal of the entire graph.
func (graph *Digraph[V]) DepthFirst() iter.Seq[TraversalEvent[V]] {
	return func(yield func(TraversalEvent[V]) bool) {
		info := make(map[V]uint8)
		var stack []dfsFrame[V]

		for vertex := range graph.edges {
			if info[vertex] != WHITE {
				continue
			}

			var zero V
			stack = append(stack, dfsFrame[V]{
				parent: zero,
				vertex: vertex,
			})

			for len(stack) > 0 {
				frame := &stack[len(stack)-1]

				if info[frame.vertex] == WHITE {
					info[frame.vertex] = GREY
					if !yield(TraversalEvent[V]{Kind: Discover, Parent: frame.parent, Vertex: frame.vertex}) {
						return
					}
					lst := graph.edges[frame.vertex]
					if lst != nil {
						frame.nextChild = lst.Front()
					}
				}

				advanced := false
				for frame.nextChild != nil {
					child := frame.nextChild.Value.(V)
					frame.nextChild = frame.nextChild.Next()

					switch info[child] {
					case WHITE:
						stack = append(stack, dfsFrame[V]{
							parent: frame.vertex,
							vertex: child,
						})
						advanced = true
					case GREY:
						if !yield(TraversalEvent[V]{Kind: BackEdge, Parent: frame.vertex, Vertex: child}) {
							return
						}
					case BLACK:
						if !yield(TraversalEvent[V]{Kind: CrossEdge, Parent: frame.vertex, Vertex: child}) {
							return
						}
					}
					if advanced {
						break
					}
				}

				if !advanced {
					info[frame.vertex] = BLACK
					if !yield(TraversalEvent[V]{Kind: Finish, Parent: frame.parent, Vertex: frame.vertex}) {
						return
					}
					stack = stack[:len(stack)-1]
				}
			}
		}
	}
}
