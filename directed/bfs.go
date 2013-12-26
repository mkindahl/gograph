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



func (graph *Graph) DoBreadthFirstWalkFromVertex(startAt Vertex,onDiscover, onFinish VertexWalkFunc) {
	walker := &fillableWalker {
		onDiscover: onDiscover,
		onFinish: onFinish,
	}

	graph.BreadthFirstWalkFromVertex(walker, startAt)
		
}

func (graph *Graph) DoBreadthFirstWalk(onDiscover, onFinish VertexWalkFunc) {
	walker := &fillableWalker {
		onDiscover: onDiscover,
		onFinish: onFinish,
	}

	graph.BreadthFirstWalk(walker)
		
}

func (graph *Graph) BreadthFirstWalk(walker Walker) {
	seen := make(map[Vertex]uint8)
	graph.DoVertices(func (vertex Vertex) error {
		if err := graph.breadthFirstVisit(walker, seen, vertex); err != nil {
			return err
		}

		return nil
	})
}
