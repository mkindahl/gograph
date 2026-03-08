// Copyright (c) 2013, Ryan Marcus. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "testing"

func TestNonCyclicBreadthFirstWalk(t *testing.T) {

	graph := New[string]()

	// this is the graph that wikipedia uses as an example for BFS
	// a visualization can be found here:
	// http://en.wikipedia.org/wiki/File:Animated_BFS.gif
	// two extra nodes (1 and 2) are added that are not connected to the rest of
	// the graph
	graph.AddEdges(
		Edge[string]{"a", "b"},
		Edge[string]{"a", "c"},
		Edge[string]{"b", "d"},
		Edge[string]{"b", "e"},
		Edge[string]{"e", "h"},
		Edge[string]{"c", "f"},
		Edge[string]{"c", "g"},
		Edge[string]{"1", "2"},
	)

	info := make(map[string]int)
	fin := 0

	for event := range graph.BreadthFirstFrom("a") {
		if event.Kind == Finish {
			info[event.Vertex] = fin
			fin++
		}
	}

	// a < b
	// a < c
	// b < d
	// b < e
	// c < f
	// c < g
	// d, e, f, g < h
	t.Logf("Ordering: %v\n", info)
	if info["a"] >= info["b"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["a"] >= info["c"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["b"] >= info["d"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["b"] >= info["e"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["c"] >= info["f"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["c"] >= info["g"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["d"] >= info["h"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["e"] >= info["h"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["f"] >= info["h"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["g"] >= info["h"] {
		t.Fatalf("Incorrect ordering")
	}

	if _, ok := info["1"]; ok {
		t.Fatalf("Found unreachable node")
	}

	if _, ok := info["2"]; ok {
		t.Fatalf("Found unreachable node")
	}

}

func TestCyclicBreadthFirstWalk(t *testing.T) {
	graph := New[string]()

	// this is the graph that wikipedia uses as an example for BFS
	// a visualization can be found here:
	// http://en.wikipedia.org/wiki/File:Animated_BFS.gif
	// with the exception that node "h" is connected to node "a"
	graph.AddEdges(
		Edge[string]{"a", "b"},
		Edge[string]{"a", "c"},
		Edge[string]{"b", "d"},
		Edge[string]{"b", "e"},
		Edge[string]{"e", "h"},
		Edge[string]{"c", "f"},
		Edge[string]{"c", "g"},
		Edge[string]{"h", "a"},
	)

	info := make(map[string]int)
	fin := 0

	for event := range graph.BreadthFirstFrom("a") {
		if event.Kind == Finish {
			info[event.Vertex] = fin
			fin++
		}
	}

	t.Logf("Ordering: %v\n", info)
	if info["a"] >= info["b"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["a"] >= info["c"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["b"] >= info["d"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["b"] >= info["e"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["c"] >= info["f"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["c"] >= info["g"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["d"] >= info["h"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["e"] >= info["h"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["f"] >= info["h"] {
		t.Fatalf("Incorrect ordering")
	}

	if info["g"] >= info["h"] {
		t.Fatalf("Incorrect ordering")
	}

}

func TestBreadthFirstWalk(t *testing.T) {
	graph := New[string]()

	graph.AddEdges(
		Edge[string]{"a", "b"},
		Edge[string]{"a", "c"},
		Edge[string]{"b", "d"},
		Edge[string]{"b", "e"},
		Edge[string]{"e", "h"},
		Edge[string]{"c", "f"},
		Edge[string]{"c", "g"},
		Edge[string]{"h", "a"},
		Edge[string]{"1", "2"},
	)

	info := make(map[string]int)
	fin := 0

	for event := range graph.BreadthFirst() {
		if event.Kind == Finish {
			info[event.Vertex] = fin
			fin++
		}
	}

	for _, v := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "1", "2"} {
		if _, ok := info[v]; !ok {
			t.Fatalf("Could not find node %s in breadth first walk", v)
		}
	}
}

func TestBreadthFirstIterator(t *testing.T) {
	graph := New[string]()
	graph.AddEdges(
		Edge[string]{"a", "b"},
		Edge[string]{"a", "c"},
		Edge[string]{"b", "d"},
		Edge[string]{"1", "2"},
	)

	discovered := make(map[string]bool)
	finished := make(map[string]bool)
	for event := range graph.BreadthFirst() {
		switch event.Kind {
		case Discover:
			discovered[event.Vertex] = true
		case Finish:
			finished[event.Vertex] = true
		}
	}
	for _, v := range []string{"a", "b", "c", "d", "1", "2"} {
		if !discovered[v] {
			t.Errorf("Vertex %s not discovered", v)
		}
		if !finished[v] {
			t.Errorf("Vertex %s not finished", v)
		}
	}
}

func TestBreadthFirstFromIterator(t *testing.T) {
	graph := New[string]()
	graph.AddEdges(
		Edge[string]{"a", "b"},
		Edge[string]{"a", "c"},
		Edge[string]{"b", "d"},
		Edge[string]{"1", "2"},
	)

	discovered := make(map[string]bool)
	for event := range graph.BreadthFirstFrom("a") {
		if event.Kind == Discover {
			discovered[event.Vertex] = true
		}
	}
	for _, v := range []string{"a", "b", "c", "d"} {
		if !discovered[v] {
			t.Errorf("Vertex %s not discovered", v)
		}
	}
	if discovered["1"] || discovered["2"] {
		t.Errorf("Unreachable vertices should not be discovered")
	}
}

func TestBreadthFirstFromIteratorBreak(t *testing.T) {
	graph := New[string]()
	graph.AddEdges(Edge[string]{"a", "b"}, Edge[string]{"b", "c"}, Edge[string]{"c", "d"})
	count := 0
	for range graph.BreadthFirstFrom("a") {
		count++
		if count >= 3 {
			break
		}
	}
	if count != 3 {
		t.Errorf("Expected 3 events, got %d", count)
	}
}
