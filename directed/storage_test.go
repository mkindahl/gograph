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
	graph := New()

	checkedAddVertex := func(vertex Vertex, expected bool, format string) {
		if graph.AddVertex(vertex) != expected {
			t.Errorf(format, vertex)
		}
	}

	checkedHasVertex := func(vertex Vertex, expected bool, format string) {
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
	graph.DoVertices(func(vertex Vertex) error {
		check[vertex.(int)] = true
		return nil
	})

	for i, b := range check {
		if !b {
			t.Errorf("Vertex %d not done", i)
		}
	}

	// Check that the iterator function abort when an error is
	// given and that it passes back the right error.
	count := 0
	err := graph.DoVertices(func(Vertex) error {
		if count > 5 {
			return errors.New("count > 5")
		}
		count++
		return nil
	})
	if err == nil {
		t.Errorf("Error not returned")
	} else if err.Error() != "count > 5" {
		t.Errorf("Incorrect error returned: %v", err)
	}
}

func TestAddEdge(t *testing.T) {
	graph := New()

	checkedAddEdge := func(x, y Vertex, expected bool, format string) {
		if graph.AddEdge(x, y) != expected {
			t.Errorf(format, x, y)
		}
	}

	checkedHasEdge := func(x, y Vertex, expected bool, format string) {
		if graph.HasEdge(x, y) != expected {
			t.Errorf(format, x, y)
		}
	}

	checkedHasVertex := func(vertex Vertex, expected bool, format string) {
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
			checkedHasEdge(i, j, true, "Edge (%v,%v) missing")
			checkedHasEdge(i+10, j, false, "Edge (%v,%v) extreneous")
			checkedHasEdge(i, j+10, false, "Edge (%v,%v) extreneous")
		}
	}

	// Check that the edge iterator function processes all edges
	// of the graph and only those.
	check := new([10][20]bool)
	graph.DoEdges(func(source, target Vertex) error {
		check[source.(int)][target.(int)] = true
		return nil
	})
	for i, js := range check {
		for j := range js {
			if check[i][j] != (0 <= i && i < 10 && 10 <= j && j < 20) {
				t.Errorf("Edge (%d,%d) not processed", i, j)
			}
		}
	}
}

func TestRemoveEdge(t *testing.T) {
	graph := New()

	checkedHasEdge := func(x, y Vertex, expected bool, format string) {
		if graph.HasEdge(x, y) != expected {
			t.Errorf(format, x, y)
		}
	}

	checkedRemoveEdge := func(source, target Vertex) {
		graph.RemoveEdge(source, target)
		checkedHasEdge(source, target, false, "Edge (%d,%d) should not be there")
	}

	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(3, 4)
	graph.AddEdge(2, 4)
	graph.AddEdge(4, 1)

	checkedRemoveEdge(4, 1)
	checkedHasEdge(1, 2, true, "Edge (%d, %d) missing")
	checkedHasEdge(1, 3, true, "Edge (%d, %d) missing")
	checkedHasEdge(3, 4, true, "Edge (%d, %d) missing")
	checkedHasEdge(2, 4, true, "Edge (%d, %d) missing")
}

func TestRemoveVertex(t *testing.T) {
	graph := New()

	checkedHasVertex := func(vertex Vertex, expected bool, format string) {
		if graph.HasVertex(vertex) != expected {
			t.Errorf(format, vertex)
		}
	}

	checkedHasEdge := func(x, y Vertex, expected bool, format string) {
		if graph.HasEdge(x, y) != expected {
			t.Errorf(format, x, y)
		}
	}

	checkedRemoveVertex := func(vtx Vertex) {
		graph.RemoveVertex(vtx)
		checkedHasVertex(vtx, false, "Vertex %d should not be there")
	}

	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(3, 4)
	graph.AddEdge(2, 4)
	graph.AddEdge(4, 1)

	checkedRemoveVertex(4)
	checkedHasEdge(1, 2, true, "Edge (%d, %d) missing")
	checkedHasEdge(1, 3, true, "Edge (%d, %d) missing")
	checkedHasEdge(3, 4, false, "Edge (%d, %d) present")
	checkedHasEdge(2, 4, false, "Edge (%d, %d) present")
	checkedHasEdge(4, 1, false, "Edge (%d, %d) present")
}

func checkGraphCount(t *testing.T, graph *Graph, vertices_expected, edges_expected int) {
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
	graph := New()
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
