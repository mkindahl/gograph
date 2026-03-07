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

func (walker *sccWalker) OnDiscover(parent, vertex Vertex) error {
	walker.time++
	vinfo := &sccInfo{
		number: walker.time,
		low:    walker.time,
	}
	walker.info[vertex] = vinfo
	walker.stack.PushBack(vertex)
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
		// Pop all vertices from stack until we find the root
		scc := New()
		for {
			elem := walker.stack.Back()
			if elem == nil {
				break
			}
			svertex := elem.Value
			walker.stack.Remove(elem)
			scc.AddVertex(svertex)
			if svertex == vertex {
				break
			}
		}

		// Only report SCCs with more than one vertex
		if scc.Order() > 1 {
			// Add all edges in the original graph between
			// vertices in the SCC graph.
			addEdge := func(source, target Vertex) error {
				if scc.HasVertex(target) {
					scc.AddEdge(source, target)
				}
				return nil
			}
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
