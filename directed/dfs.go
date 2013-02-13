// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "fmt"

// Walker interface is used by the depth-first visit function. All the
// methods have to be implemented. To help with implementing default
// methods (that do nothing) please embed the DefaultWalker.
type Walker interface {
	OnDiscover(parent, vertex Vertex) error
	OnFinish(parent, vertex Vertex) error
	OnBackEdge(source, target Vertex) error
	OnCrossEdge(source, target Vertex) error
}

// DefaultWalker implement default methods for use when implementing a
// walker. The default methods do nothing.
type DefaultWalker struct{}

// OnDiscover implement the default callback for node discovery
func (walker *DefaultWalker) OnDiscover(parent, vertex Vertex) error {
	return nil
}

// OnDiscover implement the default callback for completing nodes
func (walker *DefaultWalker) OnFinish(parent, vertex Vertex) error {
	return nil
}

// OnDiscover implement the default callback for discovering back edges
func (walker *DefaultWalker) OnBackEdge(source, target Vertex) error {
	return nil
}

// OnDiscover implement the default callback for discovering cross edges
func (walker *DefaultWalker) OnCrossEdge(source, target Vertex) error {
	return nil
}

const (
	WHITE = iota // Undiscovered
	GREY         // Discovered, but not finalized
	BLACK        // Finalized
)

// DepthFirstVisit will perform a depth-first walk starting with a
// single vertex and store the information in the 'walker' structure.
// This will be a depth-first search forest, which can be used to
// deduce other properties of the graph.
func (graph *Graph) depthFirstVisit(walker Walker, info map[Vertex]uint8, parent, vertex Vertex) error {
	switch info[vertex] {
	case WHITE:
		info[vertex] = GREY
		if err := walker.OnDiscover(parent, vertex); err != nil {
			return err
		}
		graph.DoOutEdges(vertex, func(source, target Vertex) error {
			return graph.depthFirstVisit(walker, info, source, target)
		})
		info[vertex] = BLACK
		if err := walker.OnFinish(parent, vertex); err != nil {
			return err
		}
	case GREY:
		// This is part of the tree we are recursing in, so
		// this is a back edge
		if err := walker.OnBackEdge(parent, vertex); err != nil {
			return err
		}
	case BLACK:
		// The vertex was closed, so it is a cross edge.
		if err := walker.OnCrossEdge(parent, vertex); err != nil {
			return err
		}
	}
	return nil
}

func (graph *Graph) DepthFirstWalk(walker Walker) {
	seen := make(map[Vertex]uint8)
	for vertex := range graph.edges {
		graph.depthFirstVisit(walker, seen, nil, vertex)
	}
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
type basicWalker struct {
	DefaultWalker
	time                 int
	info                 map[Vertex]*basicInfo
	onDiscover, onFinish VertexWalkFunc
}

func (walker *basicWalker) OnDiscover(parent, vertex Vertex) error {
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

func (walker *basicWalker) OnFinish(parent, vertex Vertex) error {
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
func (walker *basicWalker) String() string {
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
func (graph *Graph) DoDepthFirst(onDiscover, onFinish VertexWalkFunc) {
	walker := &basicWalker{
		onDiscover: onDiscover,
		onFinish:   onFinish,
		info:       make(map[Vertex]*basicInfo),
	}
	graph.DepthFirstWalk(walker)
}
