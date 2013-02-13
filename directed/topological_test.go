// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

package directed

import "testing"

func TestTopologicalWalk(t *testing.T) {
	graph := New()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 4)
	time := 1
	graph.DoTopological(func(vertex Vertex) error {
		if vertex.(int) != time {
			t.Errorf("Vertex %d processed at time %d\n", vertex.(int), time)
		}
		time++
		return nil
	})

	graph = New()
	graph.AddEdge(1, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(2, 4)
	graph.AddEdge(3, 4)
	when := new([5]int)
	time = 1
	graph.DoTopological(func(vertex Vertex) error {
		when[time] = vertex.(int)
		time++
		return nil
	})
	if !(when[1] < when[2] && when[1] < when[3] && when[2] < when[4] && when[3] < when[4]) {
		t.Errorf("Not in topological order %v\n", when)
	}
}
