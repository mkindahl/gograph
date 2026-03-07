// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "iter"

// Topological returns an iterator over the vertices of the graph in
// topological order.
func (graph *Graph[V]) Topological() iter.Seq[V] {
	return func(yield func(V) bool) {
		var order []V
		for event := range graph.DepthFirst() {
			if event.Kind == Finish {
				order = append(order, event.Vertex)
			}
		}
		for i := len(order) - 1; i >= 0; i-- {
			if !yield(order[i]) {
				return
			}
		}
	}
}

// DoTopological will process the graph in topological order and call
// onDiscover with each step.
func (graph *Graph[V]) DoTopological(onDiscover VertexWalkFunc[V]) error {
	for v := range graph.Topological() {
		if err := onDiscover(v); err != nil {
			return err
		}
	}
	return nil
}
