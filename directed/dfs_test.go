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
	graph.AddEdges(
		Edge[int]{1, 2},
		Edge[int]{1, 4},
		Edge[int]{2, 5},
		Edge[int]{3, 5},
		Edge[int]{3, 6},
		Edge[int]{4, 2},
		Edge[int]{5, 4},
		Edge[int]{6, 6},
	)

	info := make(map[int]*Info)
	time := 1
	for event := range graph.DepthFirst() {
		switch event.Kind {
		case Discover:
			if info[event.Vertex] == nil {
				info[event.Vertex] = new(Info)
			}
			info[event.Vertex].discover = time
			time++
		case Finish:
			if info[event.Vertex] == nil {
				info[event.Vertex] = new(Info)
			}
			info[event.Vertex].finish = time
			time++
		}
	}

	for source, target := range graph.Edges() {
		if source != target && !isValidNesting(info[source], info[target]) {
			pretty := prettyMap(info)
			t.Errorf("Edge %v -> %v has bad finish time (%s)\n", source, target, pretty)
		}
	}
}

// TestCircularGraph tests that the depth first walk stops also for
// circular graphs.
func TestCircularGraph(t *testing.T) {
	graph := New[int]()
	graph.AddEdges(Edge[int]{1, 2}, Edge[int]{2, 1})
	visited := make(map[int]int)
	for event := range graph.DepthFirst() {
		if event.Kind == Discover {
			if visited[event.Vertex] == 0 {
				visited[event.Vertex] = 1
			} else {
				t.Errorf("Edge already visited!")
			}
		}
	}
}

func TestDepthFirstIterator(t *testing.T) {
	graph := New[int]()
	graph.AddEdges(Edge[int]{1, 2}, Edge[int]{2, 3})

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
	graph.AddEdges(Edge[int]{1, 2}, Edge[int]{2, 3}, Edge[int]{3, 1})

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
	graph.AddEdges(Edge[int]{1, 2}, Edge[int]{2, 3}, Edge[int]{3, 4})
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
