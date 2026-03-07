// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import (
	"container/list"
	"fmt"
	"iter"
)

// Walker interface is used by the depth-first visit function. All the
// methods have to be implemented. To help with implementing default
// methods (that do nothing) please embed the DefaultWalker.
type Walker[V comparable] interface {
	OnDiscover(parent, vertex V) error
	OnFinish(parent, vertex V) error
	OnBackEdge(source, target V) error
	OnCrossEdge(source, target V) error
}

// DefaultWalker implement default methods for use when implementing a
// walker. The default methods do nothing.
type DefaultWalker[V comparable] struct{}

// OnDiscover implement the default callback for node discovery
func (walker *DefaultWalker[V]) OnDiscover(parent, vertex V) error {
	return nil
}

// OnFinish implement the default callback for completing nodes
func (walker *DefaultWalker[V]) OnFinish(parent, vertex V) error {
	return nil
}

// OnBackEdge implement the default callback for discovering back edges
func (walker *DefaultWalker[V]) OnBackEdge(source, target V) error {
	return nil
}

// OnCrossEdge implement the default callback for discovering cross edges
func (walker *DefaultWalker[V]) OnCrossEdge(source, target V) error {
	return nil
}

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
func (graph *Graph[V]) DepthFirst() iter.Seq[TraversalEvent[V]] {
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

// DepthFirstWalk will perform a depth-first walk of the graph using
// the provided Walker. It consumes the DepthFirst iterator and
// dispatches events to the walker's methods.
func (graph *Graph[V]) DepthFirstWalk(walker Walker[V]) {
	dispatchEvents(graph.DepthFirst(), walker)
}

// Structure holding information about the walk of an individual
// vertex. The 'discover' field is set to the time when the vertex was
// discovered (first seen in the depth-first walk), and the 'finish'
// field is set to the time when the processing of the subtree rooted
// at the vertex was finished.
type basicInfo struct {
	discover, finish int
}

// Return a string for the walk information of a vertex. Mainly used
// for debugging.
func (info *basicInfo) String() string {
	return fmt.Sprintf("%d/%d", info.discover, info.finish)
}

// A walker structure containing information about a depth-first walk
// of the graph.
type basicWalker[V comparable] struct {
	DefaultWalker[V]
	time                 int
	info                 map[V]*basicInfo
	onDiscover, onFinish VertexWalkFunc[V]
}

func (walker *basicWalker[V]) OnDiscover(parent, vertex V) error {
	walker.time++
	walker.info[vertex] = &basicInfo{
		discover: walker.time,
	}
	if walker.onDiscover != nil {
		if err := walker.onDiscover(vertex); err != nil {
			return err
		}
	}
	return nil
}

func (walker *basicWalker[V]) OnFinish(parent, vertex V) error {
	walker.time++
	walker.info[vertex].finish = walker.time
	if walker.onFinish != nil {
		if err := walker.onFinish(vertex); err != nil {
			return err
		}
	}
	return nil
}

// Return a string for the contents of the walker. Mainly used for
// debugging.
func (walker *basicWalker[V]) String() string {
	result := ""
	for k, v := range walker.info {
		result += fmt.Sprintf("%v: %v\n", k, v)
	}
	return result
}

// DoDepthFirst will process the graph in depth-first order. When a
// vertex is discovered (seen for the first time), onDiscover will be
// called with the vertex and the time it was discovered. When all the
// descendants of the vertex have been processed, onFinish will be
// called with the vertex and the time it was finished. Time in this
// case is a logical clock that is stepped each time a new node is
// discovered or finished.
func (graph *Graph[V]) DoDepthFirst(onDiscover, onFinish VertexWalkFunc[V]) {
	walker := &basicWalker[V]{
		onDiscover: onDiscover,
		onFinish:   onFinish,
		info:       make(map[V]*basicInfo),
	}
	graph.DepthFirstWalk(walker)
}
