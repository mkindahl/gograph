package directed

import "testing"

func TestStorageSimple(t *testing.T) {
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

	for i := 0 ; i < 10 ; i++ {
		for j := 0 ; j < 10 ; j++ {
			CheckedAddEdge(i, j, true, "Edge (%v,%v) cannot be added");
			CheckedAddEdge(i, j, false, "Duplicate edge (%v,%v) should not be possible to add");
		}
	}
	for i := 0 ; i < 10 ; i++ {
		for j := 0 ; j < 10 ; j++ {
			CheckedHasEdge(i, j, true, "Edge (%v,%v) missing")
			CheckedHasEdge(i+10, j, false, "Edge (%v,%v) extreneous")
			CheckedHasEdge(i, j+10, false, "Edge (%v,%v) extreneous")
		}
	}
}
