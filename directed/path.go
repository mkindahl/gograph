// Copyright (c) 2013, Ryan Marcus. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed
import "errors"
import "container/list"

type shortestPathWalker struct { 
	childOf map[Vertex]Vertex
	targetVertex Vertex
}

func (spw *shortestPathWalker) Init() {
	spw.childOf = make(map[Vertex]Vertex)
}

func (spw *shortestPathWalker) OnDiscover(parent, vertex Vertex) error {
	return nil
}

func (spw *shortestPathWalker) OnFinish(parent, vertex Vertex) error {
	spw.childOf[vertex] = parent
	if vertex == spw.targetVertex {
		return errors.New("Found node")
	}
	return nil
}

func (spw *shortestPathWalker) OnBackEdge(parent, vertex Vertex) error {
	return nil
}

func (spw *shortestPathWalker) OnCrossEdge(parent, vertex Vertex) error {
	return nil
}

// Finds the shortest path between two vertices, if such a path exists. The list
// returned will be either nil if no path was found (in which case, error will be set) or
// a list of vertices starting with start and ending with stop. This is implemented using BFS
// and has complexity O(|E|)
func (graph *Graph) FindShortestPath(start, stop Vertex) (*list.List, error) {
	w := new(shortestPathWalker)
	w.Init()
	w.targetVertex = stop

	toR := list.New()
	graph.BreadthFirstWalkFromVertex(w, start)

	last := stop
	for ; ; {
		toR.PushFront(last)

		if last == start {
			return toR, nil
		}

		next, ok := w.childOf[last]
		if !ok {
			return nil, errors.New("Could not find path")
		}
		last = next
		
	}

	
}
