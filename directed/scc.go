// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import (
	"container/list"
)

type sccInfo struct {
	number, low int
}

// GraphWalkFunc is a function called on subgraphs of a graph, for
// example, when iterating over the strongly connected components of a
// graph.
type GraphWalkFunc func(graph *Graph) error

// sccWalker is used to perform a discovery of the SCCs (Strongly
// Connected Components) in a graph.
type sccWalker struct {
	DefaultWalker
	onComponent GraphWalkFunc
	time        int
	info        map[Vertex]*sccInfo
	stack       *list.List
	graph       *Graph
}

// pushStack will push a vertex on the stack of unassigned vertices.
func (walker *sccWalker) pushStack(vertex Vertex) {
	walker.stack.PushBack(vertex)
}

// popStack will pop the topmost vertex from the stack of unassigned
// vertices and return it with an indication if there are more
// elements that need to be popped after this one.
func (walker *sccWalker) popStack(low int) (more bool, vertex Vertex) {
	if elem := walker.stack.Back(); elem != nil {
		vertex = elem.Value
		walker.stack.Remove(elem)
		top := walker.stack.Back()
		if top != nil && walker.info[top.Value].low == low {
			more = true
		}
	}
	return
}

func (walker *sccWalker) OnDiscover(parent, vertex Vertex) error {
	walker.time++
	vinfo := &sccInfo{
		number: walker.time,
		low:    walker.time,
	}
	walker.info[vertex] = vinfo
	walker.pushStack(vertex)
	if parent != nil {
		pinfo := walker.info[parent]
		if vinfo.low < pinfo.low {
			pinfo.low = vinfo.low
		}
	}
	return nil
}

func (walker *sccWalker) OnBackEdge(source, target Vertex) error {
	sinfo := walker.info[source]
	tinfo := walker.info[target]
	if tinfo.low < sinfo.low {
		sinfo.low = tinfo.low
	}
	return nil
}

func (walker *sccWalker) OnFinish(parent, vertex Vertex) error {
	vinfo := walker.info[vertex]
	pinfo := walker.info[parent]
	if pinfo != nil && vinfo.low < pinfo.low {
		pinfo.low = vinfo.low
	}

	// Check if this is an SCC root vertex
	if vinfo.number == vinfo.low {
		more, svertex := walker.popStack(vinfo.number)
		// Check if there is at least one more vertex to pop,
		// if there is, we have an SCC of size > 1
		if more {
			// Create a graph and add all vertices on the
			// stack that is part of the SCC.
			scc := New()
			scc.AddVertex(svertex)
			for more {
				more, svertex = walker.popStack(vinfo.number)
				scc.AddVertex(svertex)
			}
			// Function that check if the target vertex is
			// (also) in the SCC. In that case, the edge
			// is added to the SCC graph.
			addEdge := func(source, target Vertex) error {
				if scc.HasVertex(target) {
					scc.AddEdge(source, target)
				}
				return nil
			}
			// Add all edges in the original graph between
			// vertices in the SCC graph.
			scc.DoVertices(func(vtx Vertex) error {
				walker.graph.DoOutEdges(vtx, addEdge)
				return nil
			})
			if error := walker.onComponent(scc); error != nil {
				return error
			}
		}
	}
	return nil
}

// DoCycles will call the onComponent function for each SCC (strongly
// connected component) of size larger than 1 found in the graph.
//
// The reason that only components of size > 1 is picked is that
// Tarjan's algorithm considers nodes that are not part of a SCC of
// size larger than 1 an SCC in itself. This would mean that the
// function is called for all nodes of the graph if it is, for
// example, a DAG (directed acyclic graph).
//
// The intention is normally to use the algorithm to find cycles in
// the graph, and this behaviour defeats the purpose, so we only
// consider SCCs of size larger than 1.
func (graph *Graph) DoCycles(onComponent GraphWalkFunc) {
	walker := &sccWalker{
		graph:       graph,
		info:        make(map[Vertex]*sccInfo),
		onComponent: onComponent,
		stack:       list.New(),
	}
	graph.DepthFirstWalk(walker)
}
