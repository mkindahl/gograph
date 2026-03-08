// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import (
	"errors"
	"testing"
)

func TestAddVertex(t *testing.T) {
	graph := New[int]()

	checkedAddVertex := func(vertex int, expected bool, format string) {
		if graph.AddVertex(vertex) != expected {
			t.Errorf(format, vertex)
		}
	}

	checkedHasVertex := func(vertex int, expected bool, format string) {
		if graph.HasVertex(vertex) != expected {
			t.Errorf(format, vertex)
		}
	}

	for i := 0; i < 10; i++ {
		checkedAddVertex(i, true, "Vertex %v cannot be added")
		checkedAddVertex(i, false, "Vertex %v should not be added")
	}

	for i := 0; i < 10; i++ {
		checkedHasVertex(i, true, "Vertex %v missing")
		checkedHasVertex(i+10, false, "Extreneous vertex %v")
	}

	// Check that the iterator function process all vertices in
	// the graph.
	check := new([10]bool)
	for vertex := range graph.Vertices() {
		check[vertex] = true
	}

	for i, b := range check {
		if !b {
			t.Errorf("Vertex %d not done", i)
		}
	}

	// Check that early termination works via break.
	count := 0
	for range graph.Vertices() {
		if count > 5 {
			break
		}
		count++
	}
	if count > 6 {
		t.Errorf("Iterator did not stop early")
	}
}

func TestAddEdge(t *testing.T) {
	graph := New[int]()

	checkedAddEdge := func(x, y int, expected bool, format string) {
		if graph.AddEdge(x, y) != expected {
			t.Errorf(format, x, y)
		}
	}

	checkedAdjacent := func(x, y int, expected bool, format string) {
		if graph.Adjacent(x, y) != expected {
			t.Errorf(format, x, y)
		}
	}

	checkedHasVertex := func(vertex int, expected bool, format string) {
		if graph.HasVertex(vertex) != expected {
			t.Errorf(format, vertex)
		}
	}

	for i := 0; i < 10; i++ {
		for j := 10; j < 20; j++ {
			checkedAddEdge(i, j, true, "Edge (%v,%v) cannot be added")
			checkedAddEdge(i, j, false, "Duplicate edge (%v,%v) can be added")
		}
	}

	// Adding the edges should add the vertices as well
	for i := 0; i < 20; i++ {
		checkedHasVertex(i, true, "Vertex %v is missing")
	}

	for i := 0; i < 10; i++ {
		for j := 10; j < 20; j++ {
			checkedAdjacent(i, j, true, "Edge (%v,%v) missing")
			checkedAdjacent(i+10, j, false, "Edge (%v,%v) extreneous")
			checkedAdjacent(i, j+10, false, "Edge (%v,%v) extreneous")
		}
	}

	// Check that the edge iterator function processes all edges
	// of the graph and only those.
	check := new([10][20]bool)
	for source, target := range graph.Edges() {
		check[source][target] = true
	}
	for i, js := range check {
		for j := range js {
			if check[i][j] != (0 <= i && i < 10 && 10 <= j && j < 20) {
				t.Errorf("Edge (%d,%d) not processed", i, j)
			}
		}
	}
}

func TestRemoveEdge(t *testing.T) {
	graph := New[int]()

	checkedAdjacent := func(x, y int, expected bool, format string) {
		if graph.Adjacent(x, y) != expected {
			t.Errorf(format, x, y)
		}
	}

	checkedRemoveEdge := func(source, target int) {
		graph.RemoveEdge(source, target)
		checkedAdjacent(source, target, false, "Edge (%d,%d) should not be there")
	}

	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(3, 4)
	graph.AddEdge(2, 4)
	graph.AddEdge(4, 1)

	checkedRemoveEdge(4, 1)
	checkedAdjacent(1, 2, true, "Edge (%d, %d) missing")
	checkedAdjacent(1, 3, true, "Edge (%d, %d) missing")
	checkedAdjacent(3, 4, true, "Edge (%d, %d) missing")
	checkedAdjacent(2, 4, true, "Edge (%d, %d) missing")
}

func TestRemoveVertex(t *testing.T) {
	graph := New[int]()

	checkedHasVertex := func(vertex int, expected bool, format string) {
		if graph.HasVertex(vertex) != expected {
			t.Errorf(format, vertex)
		}
	}

	checkedAdjacent := func(x, y int, expected bool, format string) {
		if graph.Adjacent(x, y) != expected {
			t.Errorf(format, x, y)
		}
	}

	checkedRemoveVertex := func(vtx int) {
		graph.RemoveVertex(vtx)
		checkedHasVertex(vtx, false, "Vertex %d should not be there")
	}

	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(3, 4)
	graph.AddEdge(2, 4)
	graph.AddEdge(4, 1)

	checkedRemoveVertex(4)
	checkedAdjacent(1, 2, true, "Edge (%d, %d) missing")
	checkedAdjacent(1, 3, true, "Edge (%d, %d) missing")
	checkedAdjacent(3, 4, false, "Edge (%d, %d) present")
	checkedAdjacent(2, 4, false, "Edge (%d, %d) present")
	checkedAdjacent(4, 1, false, "Edge (%d, %d) present")
}

func checkGraphCount(t *testing.T, graph *Digraph[int], vertices_expected, edges_expected int) {
	vertices := graph.Order()
	if vertices_expected != vertices {
		t.Errorf("Wrong number of vertices (was %d, expected %d)", vertices, vertices_expected)
	}
	edges := graph.Size()
	if edges_expected != edges {
		t.Errorf("Wrong number of edges (was %d, expected %d)", edges, edges_expected)
	}

}

func TestVertexEdgeCount(t *testing.T) {
	graph := New[int]()
	checkGraphCount(t, graph, 0, 0)
	graph.AddVertex(1)
	checkGraphCount(t, graph, 1, 0)
	graph.AddEdge(1, 2)
	checkGraphCount(t, graph, 2, 1)
	graph.AddEdge(1, 3)
	graph.AddEdge(1, 4)
	graph.AddEdge(1, 5)
	checkGraphCount(t, graph, 5, 4)
	graph.RemoveVertex(1)
	checkGraphCount(t, graph, 4, 0)
	graph.AddEdge(2, 3)
	checkGraphCount(t, graph, 4, 1)
	graph.RemoveVertex(3)
	checkGraphCount(t, graph, 3, 0)
}

func TestVerticesIterator(t *testing.T) {
	graph := New[int]()
	for i := 0; i < 10; i++ {
		graph.AddVertex(i)
	}
	seen := make(map[int]bool)
	for v := range graph.Vertices() {
		seen[v] = true
	}
	for i := 0; i < 10; i++ {
		if !seen[i] {
			t.Errorf("Vertex %d not seen", i)
		}
	}
	if len(seen) != 10 {
		t.Errorf("Expected 10 vertices, got %d", len(seen))
	}
}

func TestVerticesIteratorBreak(t *testing.T) {
	graph := New[int]()
	for i := 0; i < 10; i++ {
		graph.AddVertex(i)
	}
	count := 0
	for range graph.Vertices() {
		count++
		if count >= 3 {
			break
		}
	}
	if count != 3 {
		t.Errorf("Expected to visit 3 vertices, got %d", count)
	}
}

func TestEdgesIterator(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(2, 3)

	type edge struct{ src, tgt int }
	seen := make(map[edge]bool)
	for src, tgt := range graph.Edges() {
		seen[edge{src, tgt}] = true
	}
	expected := []edge{{1, 2}, {1, 3}, {2, 3}}
	for _, e := range expected {
		if !seen[e] {
			t.Errorf("Edge (%d,%d) not seen", e.src, e.tgt)
		}
	}
	if len(seen) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(seen))
	}
}

func TestEdgesIteratorBreak(t *testing.T) {
	graph := New[int]()
	for i := 0; i < 5; i++ {
		graph.AddEdge(i, i+1)
	}
	count := 0
	for range graph.Edges() {
		count++
		if count >= 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("Expected to visit 2 edges, got %d", count)
	}
}

func TestOutEdgesIterator(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(1, 4)
	graph.AddEdge(2, 5)

	targets := make(map[int]bool)
	for _, tgt := range graph.OutEdges(1) {
		targets[tgt] = true
	}
	for _, expected := range []int{2, 3, 4} {
		if !targets[expected] {
			t.Errorf("Target %d not seen in out-edges of 1", expected)
		}
	}
	if targets[5] {
		t.Errorf("Target 5 should not be in out-edges of 1")
	}
}

func TestOutEdgesIteratorMissingVertex(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	count := 0
	for range graph.OutEdges(99) {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 out-edges for missing vertex, got %d", count)
	}
}

func TestAddVertexAbortOnError(t *testing.T) {
	graph := New[int]()
	for i := 0; i < 10; i++ {
		graph.AddVertex(i)
	}
	count := 0
	var firstErr error
	for v := range graph.Vertices() {
		_ = v
		if count > 5 {
			firstErr = errors.New("count > 5")
			break
		}
		count++
	}
	if firstErr == nil {
		t.Errorf("Error not returned")
	} else if firstErr.Error() != "count > 5" {
		t.Errorf("Incorrect error returned: %v", firstErr)
	}
}
