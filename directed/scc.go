// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "iter"

type sccInfo struct {
	number, low int
}

// Cycles returns an iterator over the strongly connected components
// of size larger than 1 found in the graph. Each yielded value is a
// subgraph containing the vertices and edges of that SCC.
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
func (graph *Digraph[V]) Cycles() iter.Seq[*Digraph[V]] {
	return func(yield func(*Digraph[V]) bool) {
		info := make(map[V]*sccInfo)
		var stack []V
		time := 0

		for event := range graph.DepthFirst() {
			switch event.Kind {
			case Discover:
				time++
				vinfo := &sccInfo{number: time, low: time}
				info[event.Vertex] = vinfo
				stack = append(stack, event.Vertex)
				var zero V
				if event.Parent != zero {
					pinfo := info[event.Parent]
					if vinfo.low < pinfo.low {
						pinfo.low = vinfo.low
					}
				}

			case BackEdge:
				sinfo := info[event.Parent]
				tinfo := info[event.Vertex]
				if tinfo.low < sinfo.low {
					sinfo.low = tinfo.low
				}

			case Finish:
				vinfo := info[event.Vertex]
				var zero V
				if event.Parent != zero {
					pinfo := info[event.Parent]
					if vinfo.low < pinfo.low {
						pinfo.low = vinfo.low
					}
				}

				if vinfo.number == vinfo.low {
					scc := New[V]()
					for {
						if len(stack) == 0 {
							break
						}
						svertex := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						scc.AddVertex(svertex)
						if svertex == event.Vertex {
							break
						}
					}

					if scc.Order() > 1 {
						for vtx := range scc.Vertices() {
							for _, target := range graph.OutEdges(vtx) {
								if scc.HasVertex(target) {
									scc.AddEdges(Edge[V]{vtx, target})
								}
							}
						}
						if !yield(scc) {
							return
						}
					}
				}
			}
		}
	}
}
