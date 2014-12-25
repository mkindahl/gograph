// Copyright (c) 2013, Ryan Marcus. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "testing"

func TestShortestPath(t *testing.T) {

	graph := New()

	graph.AddEdge("a", "b")
	graph.AddEdge("a", "c")
	graph.AddEdge("b", "d")
	graph.AddEdge("b", "e")
	graph.AddEdge("e", "h")
	graph.AddEdge("c", "f")
	graph.AddEdge("c", "g")

	graph.AddEdge("1", "2")

	// the shortest path from a to h should be: a b e h
	path, err := graph.FindShortestPath("a", "h")
	
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}
	

	curr := path.Front()
	if curr.Value != "a" {
		t.Fatalf("Invalid path\n")
	}

	curr = curr.Next()
	if curr.Value != "b" {
		t.Fatalf("Invalid path\n")
	}
	
	curr = curr.Next()
	if curr.Value != "e" {
		t.Fatalf("Invalid path\n")
	}

	curr = curr.Next()
	if curr.Value != "h" {
		t.Fatalf("Invalid path\n")
	}

	curr = curr.Next()
	if curr != nil {
		t.Fatalf("Path extends beyond endpoint\n")
	}
	
	
}
