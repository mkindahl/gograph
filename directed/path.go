// Copyright (c) 2013, Ryan Marcus. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import (
	"container/list"
	"errors"
)

// ShortestPath finds the shortest path between two vertices, returning
// a slice of vertices from start to stop. Returns an error if no path
// exists.
func (graph *Graph[V]) ShortestPath(start, stop V) ([]V, error) {
	childOf := make(map[V]V)

	for event := range graph.BreadthFirstFrom(start) {
		if event.Kind == Finish {
			childOf[event.Vertex] = event.Parent
			if event.Vertex == stop {
				break
			}
		}
	}

	var path []V
	last := stop
	for {
		path = append(path, *new(V))
		copy(path[1:], path)
		path[0] = last

		if last == start {
			return path, nil
		}

		next, ok := childOf[last]
		if !ok {
			return nil, errors.New("Could not find path")
		}
		last = next
	}
}

// FindShortestPath finds the shortest path between two vertices. The list
// returned will be either nil if no path was found (in which case, error
// will be set) or a list of vertices starting with start and ending with
// stop. This is implemented using BFS and has complexity O(|E|).
func (graph *Graph[V]) FindShortestPath(start, stop V) (*list.List, error) {
	path, err := graph.ShortestPath(start, stop)
	if err != nil {
		return nil, err
	}
	result := list.New()
	for _, v := range path {
		result.PushBack(v)
	}
	return result, nil
}
