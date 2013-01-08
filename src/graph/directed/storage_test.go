package directed

import "testing"

func TestAddVertex(t *testing.T) {
	graph := New()

	CheckedAddVertex := func (vertex Vertex, expected bool, format string) {
		if graph.AddVertex(vertex) != expected {
			t.Errorf(format, vertex)
		}
	}

	CheckedHasVertex := func (vertex Vertex, expected bool, format string) {
		if graph.HasVertex(vertex) != expected {
			t.Errorf(format, vertex)
		}
	}

	for i := 0 ; i < 10 ; i++ {
		CheckedAddVertex(i, true, "Vertex %v cannot be added")
		CheckedAddVertex(i, false, "Vertex %v should not be added")
	}

	for i := 0 ; i < 10 ; i++ {
		CheckedHasVertex(i, true, "Vertex %v missing")
		CheckedHasVertex(i+10, false, "Extreneous vertex %v")
	}

	// Check that the iterator function process all vertices in
	// the graph.
	check := new([10]bool)
	graph.DoVertices(func (vertex Vertex) {
		check[vertex.(int)] = true
	})

	for i, b := range check {
		if !b {
			t.Errorf("Vertex %d not done", i)
		}
	}
}

func TestAddEdge(t *testing.T) {
	graph := New()

	CheckedAddEdge := func (x, y Vertex, expected bool, format string) {
		if graph.AddEdge(x,y) != expected {
			t.Errorf(format, x, y)
		}
	}

	CheckedHasEdge := func (x, y Vertex, expected bool, format string) {
		if graph.HasEdge(x,y) != expected {
			t.Errorf(format, x, y)
		}
	}

	CheckedHasVertex := func (vertex Vertex, expected bool, format string) {
		if graph.HasVertex(vertex) != expected {
			t.Errorf(format, vertex)
		}
	}

	for i := 0 ; i < 10 ; i++ {
		for j := 10 ; j < 20 ; j++ {
			CheckedAddEdge(i, j, true, "Edge (%v,%v) cannot be added");
			CheckedAddEdge(i, j, false, "Duplicate edge (%v,%v) should not be possible to add");
		}
	}

	// Adding the edges should add the vertices as well
	for i := 0 ; i < 20 ; i++ {
		CheckedHasVertex(i, true, "Vertex %v is missing")
	}

	for i := 0 ; i < 10 ; i++ {
		for j := 10 ; j < 20 ; j++ {
			CheckedHasEdge(i, j, true, "Edge (%v,%v) missing")
			CheckedHasEdge(i+10, j, false, "Edge (%v,%v) extreneous")
			CheckedHasEdge(i, j+10, false, "Edge (%v,%v) extreneous")
		}
	}
}
