// Copyright (c) 2013, Ryan Marcus. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "testing"

func TestNonCyclicBreadthFirstWalk(t *testing.T) {

	graph := New()

	// this is the graph that wikipedia uses as an example for BFS
	// a visualization can be found here:
	// http://en.wikipedia.org/wiki/File:Animated_BFS.gif
	// two extra nodes (1 and 2) are added that are not connected to the rest of
	// the graph
	graph.AddEdge("a", "b")
	graph.AddEdge("a", "c")
	graph.AddEdge("b", "d")
	graph.AddEdge("b", "e")
	graph.AddEdge("e", "h")
	graph.AddEdge("c", "f")
	graph.AddEdge("c", "g")

	graph.AddEdge("1", "2")

	info := make(map[Vertex]int)
	fin := 0

	onDiscover := func(vertex Vertex) error {
		return nil
	}

	onFinish := func(vertex Vertex) error {
		info[vertex] = fin
		fin++
		return nil
	}

	graph.DoBreadthFirstWalkFromVertex("a", onDiscover, onFinish)

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

	if _, ok := info['1']; ok {
		t.Fatalf("Found unreachable node")
	}


	if _, ok := info['2']; ok {
		t.Fatalf("Found unreachable node")
	}

}

func TestCyclicBreadthFirstWalk(t *testing.T) {
	graph := New()

	// this is the graph that wikipedia uses as an example for BFS
	// a visualization can be found here:
	// http://en.wikipedia.org/wiki/File:Animated_BFS.gif
	// with the exception that node "h" is connected to node "a"
	graph.AddEdge("a", "b")
	graph.AddEdge("a", "c")
	graph.AddEdge("b", "d")
	graph.AddEdge("b", "e")
	graph.AddEdge("e", "h")
	graph.AddEdge("c", "f")
	graph.AddEdge("c", "g")
	graph.AddEdge("h", "a")


	info := make(map[Vertex]int)
	fin := 0

	onDiscover := func(vertex Vertex) error {
		return nil
	}

	onFinish := func(vertex Vertex) error {
		info[vertex] = fin
		fin++
		return nil
	}

	graph.DoBreadthFirstWalkFromVertex("a", onDiscover, onFinish)

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

}

func TestBreadthFirstWalk(t *testing.T) {
	graph := New()

	// this is the graph that wikipedia uses as an example for BFS
	// a visualization can be found here:
	// http://en.wikipedia.org/wiki/File:Animated_BFS.gif
	// with the exception that node "h" is connected to node "a"
	// and there is a non-connected segment containing two nodes, "1" and "2"
	graph.AddEdge("a", "b")
	graph.AddEdge("a", "c")
	graph.AddEdge("b", "d")
	graph.AddEdge("b", "e")
	graph.AddEdge("e", "h")
	graph.AddEdge("c", "f")
	graph.AddEdge("c", "g")
	graph.AddEdge("h", "a")

	graph.AddEdge("1", "2")


	info := make(map[Vertex]int)
	fin := 0

	onDiscover := func(vertex Vertex) error {
		return nil
	}

	onFinish := func(vertex Vertex) error {
		info[vertex] = fin
		fin++
		return nil
	}

	graph.DoBreadthFirstWalk(onDiscover, onFinish)

	if _, ok := info["a"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}

	if _, ok := info["b"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}

	if _, ok := info["c"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}

	if _, ok := info["d"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}

	if _, ok := info["e"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}

	if _, ok := info["f"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}

	if _, ok := info["g"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}

	if _, ok := info["h"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}

	if _, ok := info["1"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}

	if _, ok := info["2"]; !ok {
		t.Fatalf("Could not find all nodes in breadth first walk")
	}
	

}
