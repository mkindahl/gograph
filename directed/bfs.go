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

// BreadthFirstWalkFromVertex uses the passed walker to traverse
// the graph breadth-first. The search starts at the given vertex.
// Vertices with no path to the given vertex will NOT be discovered.
func (graph *Graph[V]) BreadthFirstWalkFromVertex(walker Walker[V], vertex V) {
	dispatchEvents(graph.BreadthFirstFrom(vertex), walker)
}

// BreadthFirstWalk uses the provided walker to perform a breadth-first search
// over the entire graph, starting at an arbitrary vertex, and completing
// after all nodes in the graph are traversed.
func (graph *Graph[V]) BreadthFirstWalk(walker Walker[V]) {
	dispatchEvents(graph.BreadthFirst(), walker)
}

// dispatchEvents consumes a TraversalEvent iterator and dispatches
// each event to the appropriate Walker method.
func dispatchEvents[V comparable](events iter.Seq[TraversalEvent[V]], walker Walker[V]) {
	for event := range events {
		var err error
		switch event.Kind {
		case Discover:
			err = walker.OnDiscover(event.Parent, event.Vertex)
		case Finish:
			err = walker.OnFinish(event.Parent, event.Vertex)
		case BackEdge:
			err = walker.OnBackEdge(event.Parent, event.Vertex)
		case CrossEdge:
			err = walker.OnCrossEdge(event.Parent, event.Vertex)
		}
		if err != nil {
			return
		}
	}
}

type fillableWalker[V comparable] struct {
	onDiscover, onFinish VertexWalkFunc[V]
}

func (w *fillableWalker[V]) OnDiscover(parent, vertex V) error {
	return w.onDiscover(vertex)
}

func (w *fillableWalker[V]) OnFinish(parent, vertex V) error {
	return w.onFinish(vertex)
}

func (w *fillableWalker[V]) OnBackEdge(parent, vertex V) error {
	return nil
}

func (w *fillableWalker[V]) OnCrossEdge(parent, vertex V) error {
	return nil
}

// DoBreadthFirstWalkFromVertex performs a breadth-first search starting at
// the given vertex, calling the onDiscover function when a new vertex is
// discovered and the onFinish function  when a vertex has been traversed.
func (graph *Graph[V]) DoBreadthFirstWalkFromVertex(startAt V, onDiscover, onFinish VertexWalkFunc[V]) {
	walker := &fillableWalker[V]{
		onDiscover: onDiscover,
		onFinish:   onFinish,
	}
	graph.BreadthFirstWalkFromVertex(walker, startAt)
}

// DoBreadthFirstWalk performs a breadth-first search walk over the entire
// graph, starting at an  arbitrary vertex and completing after all nodes in
// the graph are traversed.
func (graph *Graph[V]) DoBreadthFirstWalk(onDiscover, onFinish VertexWalkFunc[V]) {
	walker := &fillableWalker[V]{
		onDiscover: onDiscover,
		onFinish:   onFinish,
	}
	graph.BreadthFirstWalk(walker)
}
