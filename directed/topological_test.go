// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "testing"

func TestTopologicalWalk(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 4)
	time := 1
	for v := range graph.Topological() {
		if v != time {
			t.Errorf("Vertex %d processed at time %d\n", v, time)
		}
		time++
	}

	graph = New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(2, 4)
	graph.AddEdge(3, 4)
	when := new([5]int)
	time = 1
	for v := range graph.Topological() {
		when[time] = v
		time++
	}
	if !(when[1] < when[2] && when[1] < when[3] && when[2] < when[4] && when[3] < when[4]) {
		t.Errorf("Not in topological order %v\n", when)
	}
}

func TestTopologicalIterator(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 4)
	time := 1
	for v := range graph.Topological() {
		if v != time {
			t.Errorf("Vertex %d processed at time %d", v, time)
		}
		time++
	}
}

func TestTopologicalIteratorDiamond(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(2, 4)
	graph.AddEdge(3, 4)
	when := make(map[int]int)
	time := 1
	for v := range graph.Topological() {
		when[v] = time
		time++
	}
	if !(when[1] < when[2] && when[1] < when[3] && when[2] < when[4] && when[3] < when[4]) {
		t.Errorf("Not in topological order %v", when)
	}
}

func TestTopologicalIteratorBreak(t *testing.T) {
	graph := New[int]()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 4)
	count := 0
	for range graph.Topological() {
		count++
		if count >= 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("Expected 2, got %d", count)
	}
}
