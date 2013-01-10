// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

// Implementation of depth-first search and algorithms based on
// depth-first walks.
package directed

import (
	"fmt"
	"container/list"
)

// Structure holding information about the walk of an individual
// vertex. The 'discover' field is set to the time when the vertex was
// discovered (first seen in the depth-first walk), and the 'finish'
// field is set to the time when the processing of the subtree rooted
// at the vertex was finished.
type walkInfo struct {
	discover, finish int
}

// Return a string for the walk information of a vertex. Mainly used
// for debugging.
func (info walkInfo) String() string {
	return fmt.Sprintf("%d/%d", info.discover, info.finish)
}


// A walker structure containing information about a depth-first walk
// of the graph.
type walker struct {
	time int
	graph *Graph
	info map[Vertex]*walkInfo
	onDiscover, onFinish VertexWalkFunc
}

// Return a string for the contents of the walker. Mainly used for
// debugging.
func (walker *walker) String() string {
	result := ""
	for k, v := range walker.info {
		result += fmt.Sprintf("%v: %v\n", k, v)
	}
	return result
}

// DepthFirstVisit will perform a depth-first walk starting with a
// single vertex and store the information in the 'walker' structure.
// This will be a depth-first search forest, which can be used to
// deduce other properties of the graph.
func (walker *walker) depthFirstVisit(vertex Vertex) {
	if walker.info[vertex] == nil {
		walker.time++
		walker.info[vertex] = &walkInfo{discover: walker.time}
		if (walker.onDiscover != nil) {
			walker.onDiscover(vertex)
		}
		walker.graph.DoOutEdges(vertex, func (source, target Vertex) {
			walker.depthFirstVisit(target)
		})
		walker.time++
		walker.info[vertex].finish = walker.time
		if (walker.onFinish != nil) {
			walker.onFinish(vertex)
		}
	}
}

// DoDepthFirst will process the graph in depth-first order. When a
// vertex is discovered (seen for the first time), onDiscover will be
// called with the vertex and the time it was discovered. When all the
// descendants of the vertex have been processed, onFinish will be
// called with the vertex and the time it was finished. Time in this
// case is a logical clock that is stepped each time a new node is
// discovered or finished.
func (graph *Graph) DoDepthFirst(onDiscover, onFinish VertexWalkFunc) {
	walker := &walker{
		graph: graph,
		onDiscover: onDiscover,
		onFinish: onFinish,
		info: make(map[Vertex]*walkInfo),
	}
	for vertex := range graph.edges {
		walker.depthFirstVisit(vertex)
	}
}

// DoTopological will process the graph in topological order and call
// onDiscover with each step.
func (graph *Graph) DoTopological(onDiscover VertexWalkFunc) {
	lst := list.New()
	graph.DoDepthFirst(nil, func (vertex Vertex) {
		lst.PushFront(vertex)
	})
	// Process elements in reverse order of finishing time
	for elem := lst.Front() ; elem != nil ; elem = elem.Next() {
		onDiscover(elem.Value)
	}
}
