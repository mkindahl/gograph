// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "testing"

func checkCycleCount(t *testing.T, graph *Graph, expected int) {
	count := 0
	graph.DoCycles(func(graph *Graph) error {
		count++
		return nil
	})
	if count != expected {
		t.Errorf("Wrong number of components (was %v, should be %v)", count, expected)
	}
}

func checkCycle(t *testing.T, graph *Graph, vertices, edges int, check func(graph *Graph) bool) {
	graph.DoCycles(func(subg *Graph) error {
		if vs := subg.Order(); vs != vertices {
			t.Errorf("Wrong number of vertices (was %d, expected %d)", vs, vertices)
		}

		if es := subg.Size(); es != edges {
			t.Errorf("Wrong number of edges (was %d, expected %d)", es, edges)
		}
		if !check(graph) {
			t.Errorf("Missing edges")
		}
		return nil
	})
}

func TestCycleWalkCount(t *testing.T) {
	graph := New()
	graph.AddEdge(1, 1)
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	checkCycleCount(t, graph, 0)
	graph.AddEdge(3, 1)
	checkCycleCount(t, graph, 1)
	checkCycle(t, graph, 3, 4, func(graph *Graph) bool {
		return graph.HasEdge(1, 2) && graph.HasEdge(2, 3) && graph.HasEdge(3, 1)
	})
	graph.AddEdge(1, 4)
	graph.AddEdge(2, 5)
	graph.AddEdge(3, 6)
	graph.AddEdge(6, 7)
	checkCycleCount(t, graph, 1)
	checkCycle(t, graph, 3, 4, func(graph *Graph) bool {
		return graph.HasEdge(1, 2) && graph.HasEdge(2, 3) && graph.HasEdge(3, 1)
	})
	graph.AddEdge(4, 1)
	checkCycleCount(t, graph, 1)
	checkCycle(t, graph, 4, 6, func(graph *Graph) bool {
		return graph.HasEdge(1, 2) && graph.HasEdge(2, 3) &&
			graph.HasEdge(3, 1) && graph.HasEdge(1, 4) &&
			graph.HasEdge(4, 1)
	})
	graph.AddEdge(7, 6)
	checkCycleCount(t, graph, 2)
}
