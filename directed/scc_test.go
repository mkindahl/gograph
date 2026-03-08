// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "testing"

func checkCycleCount(t *testing.T, graph *Digraph[int], expected int) {
	count := 0
	for range graph.Cycles() {
		count++
	}
	if count != expected {
		t.Errorf("Wrong number of components (was %v, should be %v)", count, expected)
	}
}

func checkCycle(t *testing.T, graph *Digraph[int], vertices, edges int, check func(graph *Digraph[int]) bool) {
	for subg := range graph.Cycles() {
		if vs := subg.Order(); vs != vertices {
			t.Errorf("Wrong number of vertices (was %d, expected %d)", vs, vertices)
		}

		if es := subg.Size(); es != edges {
			t.Errorf("Wrong number of edges (was %d, expected %d)", es, edges)
		}
		if !check(graph) {
			t.Errorf("Missing edges")
		}
	}
}

func TestCycleWalkCount(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 1)
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	checkCycleCount(t, graph, 0)
	graph.AddEdge(3, 1)
	checkCycleCount(t, graph, 1)
	checkCycle(t, graph, 3, 4, func(graph *Digraph[int]) bool {
		return graph.Adjacent(1, 2) && graph.Adjacent(2, 3) && graph.Adjacent(3, 1)
	})
	graph.AddEdge(1, 4)
	graph.AddEdge(2, 5)
	graph.AddEdge(3, 6)
	graph.AddEdge(6, 7)
	checkCycleCount(t, graph, 1)
	checkCycle(t, graph, 3, 4, func(graph *Digraph[int]) bool {
		return graph.Adjacent(1, 2) && graph.Adjacent(2, 3) && graph.Adjacent(3, 1)
	})
	graph.AddEdge(4, 1)
	checkCycleCount(t, graph, 1)
	checkCycle(t, graph, 4, 6, func(graph *Digraph[int]) bool {
		return graph.Adjacent(1, 2) && graph.Adjacent(2, 3) &&
			graph.Adjacent(3, 1) && graph.Adjacent(1, 4) &&
			graph.Adjacent(4, 1)
	})
	graph.AddEdge(7, 6)
	checkCycleCount(t, graph, 2)
}

func TestCyclesIterator(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 1)
	graph.AddEdge(6, 7)
	graph.AddEdge(7, 6)
	count := 0
	for range graph.Cycles() {
		count++
	}
	if count != 2 {
		t.Errorf("Expected 2 SCCs, got %d", count)
	}
}

func TestCyclesIteratorNoCycles(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	count := 0
	for range graph.Cycles() {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 SCCs, got %d", count)
	}
}

func TestCyclesIteratorBreak(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 1)
	graph.AddEdge(6, 7)
	graph.AddEdge(7, 6)
	count := 0
	for range graph.Cycles() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("Expected 1, got %d", count)
	}
}
