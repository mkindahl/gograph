// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import (
	"fmt"
	"testing"
)

// Structure holding discovery-finish information for this test.
type Info struct {
	finish, discover int
}

// Test if x is nested inside y
func (x *Info) isNestedIn(y *Info) bool {
	return y.discover < x.discover && x.discover < x.finish && x.finish < y.finish
}

// Test if x is disjoint with y
func (x *Info) isDisjoint(y *Info) bool {
	return x.discover < x.finish && x.finish < y.discover && y.discover < y.finish ||
		y.discover < y.finish && y.finish < x.discover && x.discover < x.finish
}

// Test that the discovery-finish information satisfies the
// parantheses theorem
func isValidNesting(x, y *Info) bool {
	return x.isDisjoint(y) || x.isNestedIn(y) || y.isNestedIn(x)
}

func prettyMap(input map[int]*Info) string {
	result := "{ "
	for k, v := range input {
		result += fmt.Sprintf("%v: %v/%v, ", k, v.discover, v.finish)
	}
	return result + "}"
}

// Test the depth first walk function using an example from book by
// Cormen et.al.
func TestDepthFirstWalk(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 4)
	graph.AddEdge(2, 5)
	graph.AddEdge(3, 5)
	graph.AddEdge(3, 6)
	graph.AddEdge(4, 2)
	graph.AddEdge(5, 4)
	graph.AddEdge(6, 6)

	info := make(map[int]*Info)
	time := 1
	onDiscover := func(vertex int) error {
		if info[vertex] == nil {
			info[vertex] = new(Info)
		}
		info[vertex].discover = time
		time++
		return nil
	}
	onFinish := func(vertex int) error {
		if info[vertex] == nil {
			info[vertex] = new(Info)
		}
		info[vertex].finish = time
		time++
		return nil
	}
	graph.DoDepthFirst(onDiscover, onFinish)

	printEdge := func(source, target int) error {
		if source != target && !isValidNesting(info[source], info[target]) {
			pretty := prettyMap(info)
			t.Errorf("Edge %v -> %v has bad finish time (%s)\n", source, target, pretty)
		}
		return nil
	}

	graph.DoEdges(printEdge)
}

// TestCircularGraph tests that the depth first walk stops also for
// circular graphs.
func TestCircularGraph(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 1)
	info := make(map[int]int)
	graph.DoDepthFirst(func(vertex int) error {
		if info[vertex] == 0 {
			info[vertex] = 1
		} else {
			t.Errorf("Edge already visited!")
		}
		return nil
	}, nil)
}

func TestDepthFirstIterator(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)

	discovered := make(map[int]bool)
	finished := make(map[int]bool)
	for event := range graph.DepthFirst() {
		switch event.Kind {
		case Discover:
			discovered[event.Vertex] = true
		case Finish:
			finished[event.Vertex] = true
		}
	}
	for _, v := range []int{1, 2, 3} {
		if !discovered[v] {
			t.Errorf("Vertex %d not discovered", v)
		}
		if !finished[v] {
			t.Errorf("Vertex %d not finished", v)
		}
	}
}

func TestDepthFirstIteratorBackEdge(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 1)

	hasBackEdge := false
	for event := range graph.DepthFirst() {
		if event.Kind == BackEdge {
			hasBackEdge = true
		}
	}
	if !hasBackEdge {
		t.Errorf("Expected a back edge in cyclic graph")
	}
}

func TestDepthFirstIteratorBreak(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 4)
	count := 0
	for range graph.DepthFirst() {
		count++
		if count >= 3 {
			break
		}
	}
	if count != 3 {
		t.Errorf("Expected 3 events, got %d", count)
	}
}
