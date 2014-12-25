// Copyright (c) 2013, Ryan Marcus. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "container/list"


type vertexPair struct {
	parent Vertex
	child  Vertex
}

func (graph *Graph) breadthFirstVisit(walker Walker, seen map[Vertex]uint8, vertex Vertex) error {
	queue := list.New()

	vp := new(vertexPair)
	vp.parent = nil
	vp.child = vertex

	seen[vp.child] = GREY
	if err := walker.OnDiscover(vp.parent, vp.child); err != nil {
		return err
	}

	queue.PushFront(vp)

	for queue.Len() != 0 {
		ve := queue.Front()
		queue.Remove(ve)
		v := ve.Value.(*vertexPair)


		graph.DoOutEdges(v.child, func(from Vertex, to Vertex) error {
			if seen[to] == WHITE {
				seen[to] = GREY
				if err := walker.OnDiscover(from, to); err != nil {
					return err
				}

				nvp := new(vertexPair)
				nvp.parent = from
				nvp.child = to
				queue.PushBack(nvp)
			}
			return nil
		})

		seen[v.child] = BLACK
		if err := walker.OnFinish(v.parent, v.child); err != nil {
			return err
		}
	}

	return nil
}

// BreadthFirstWalkFromVertex uses the passed walker to traverse 
// the graph breadth-first. The search starts at the given vertex. 
// Vertices with no path to the given vertex will NOT be discovered.
func (graph *Graph) BreadthFirstWalkFromVertex(walker Walker, vertex Vertex) {
	seen := make(map[Vertex]uint8)
	graph.breadthFirstVisit(walker, seen, vertex)

}

type fillableWalker struct {
	onDiscover, onFinish VertexWalkFunc
}

func (w *fillableWalker) OnDiscover(parent, vertex Vertex) error {
	return w.onDiscover(vertex)
}

func (w *fillableWalker) OnFinish(parent, vertex Vertex) error {
	return w.onFinish(vertex)
}

func (w *fillableWalker) OnBackEdge(parent, vertex Vertex) error {
	return nil
}

func (w *fillableWalker) OnCrossEdge(parent, vertex Vertex) error {
	return nil
}


// DoBreadthFirstWalkFromVertex performs a breadth-first search starting at 
// the given vertex, calling the onDiscover function when a new vertex is
// discovered and the onFinish function  when a vertex has been traversed.
func (graph *Graph) DoBreadthFirstWalkFromVertex(startAt Vertex,onDiscover, onFinish VertexWalkFunc) {
	walker := &fillableWalker {
		onDiscover: onDiscover,
		onFinish: onFinish,
	}

	graph.BreadthFirstWalkFromVertex(walker, startAt)
		
}

// DoBreadthFirstWalk performs a breadth-first search walk over the entire 
// graph, starting at an  arbitrary vertex and completing after all nodes in 
// the graph are traversed.
func (graph *Graph) DoBreadthFirstWalk(onDiscover, onFinish VertexWalkFunc) {
	walker := &fillableWalker {
		onDiscover: onDiscover,
		onFinish: onFinish,
	}

	graph.BreadthFirstWalk(walker)
		
}

// BreadthFirstWalk uses the provided walker to perform a breadth-first search 
// over the entire graph, starting at an arbitrary vertex, and completing 
// after all nodes in the graph are traversed.
func (graph *Graph) BreadthFirstWalk(walker Walker) {
	seen := make(map[Vertex]uint8)
	graph.DoVertices(func (vertex Vertex) error {
		if err := graph.breadthFirstVisit(walker, seen, vertex); err != nil {
			return err
		}

		return nil
	})
}
