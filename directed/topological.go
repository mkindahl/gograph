// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "container/list"

type topologicalWalker struct {
	DefaultWalker
	vertices *list.List
}

func (walker *topologicalWalker) OnFinish(parent, vertex Vertex) error {
	walker.vertices.PushFront(vertex)
	return nil
}

// DoTopological will process the graph in topological order and call
// onDiscover with each step.
func (graph *Graph) DoTopological(onDiscover VertexWalkFunc) error {
	walker := &topologicalWalker{
		vertices: list.New(),
	}
	graph.DepthFirstWalk(walker)
	// Process elements in reverse order of finishing time
	for elem := walker.vertices.Front(); elem != nil; elem = elem.Next() {
		if err := onDiscover(elem.Value); err != nil {
			return err
		}
	}
	return nil
}
