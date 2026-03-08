// Copyright (c) 2013, Ryan Marcus. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "testing"

func TestShortestPathSlice(t *testing.T) {
	graph := New[string]()
	graph.AddEdge("a", "b")
	graph.AddEdge("a", "c")
	graph.AddEdge("b", "d")
	graph.AddEdge("b", "e")
	graph.AddEdge("e", "h")
	graph.AddEdge("c", "f")
	graph.AddEdge("c", "g")

	path, err := graph.ShortestPath("a", "h")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	expected := []string{"a", "b", "e", "h"}
	if len(path) != len(expected) {
		t.Fatalf("Expected path length %d, got %d", len(expected), len(path))
	}
	for i, v := range path {
		if v != expected[i] {
			t.Errorf("Path[%d] = %v, expected %s", i, v, expected[i])
		}
	}
}

func TestShortestPathSliceNoPath(t *testing.T) {
	graph := New[string]()
	graph.AddEdge("a", "b")
	graph.AddEdge("c", "d")

	_, err := graph.ShortestPath("a", "d")
	if err == nil {
		t.Fatalf("Expected error for disconnected vertices")
	}
}

func TestShortestPathSliceSameVertex(t *testing.T) {
	graph := New[string]()
	graph.AddVertex("a")

	path, err := graph.ShortestPath("a", "a")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(path) != 1 || path[0] != "a" {
		t.Errorf("Expected [a], got %v", path)
	}
}
