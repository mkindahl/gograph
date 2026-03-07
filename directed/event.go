// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

// EventKind represents the type of event during a graph traversal.
type EventKind int

const (
	Discover  EventKind = iota // Vertex discovered (first seen)
	Finish                     // Vertex finished (all descendants processed)
	BackEdge                   // Back edge found (target is an ancestor)
	CrossEdge                  // Cross edge found (target already finished)
)

// TraversalEvent represents a single event during a DFS or BFS
// traversal of a graph.
type TraversalEvent[V comparable] struct {
	Kind           EventKind
	Parent, Vertex V
}
