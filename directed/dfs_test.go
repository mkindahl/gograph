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
	graph := New()
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
	onDiscover := func(vertex Vertex) error {
		if info[vertex.(int)] == nil {
			info[vertex.(int)] = new(Info)
		}
		info[vertex.(int)].discover = time
		time++
		return nil
	}
	onFinish := func(vertex Vertex) error {
		if info[vertex.(int)] == nil {
			info[vertex.(int)] = new(Info)
		}
		info[vertex.(int)].finish = time
		time++
		return nil
	}
	graph.DoDepthFirst(onDiscover, onFinish)
	graph.DoEdges(func(source, target Vertex) error {
		if source != target && !isValidNesting(info[source.(int)], info[target.(int)]) {
			pretty := prettyMap(info)
			t.Errorf("Edge %v -> %v has bad finish time (%s)\n", source, target, pretty)
		}
		return nil
	})
}

// TestCircularGraph tests that the depth first walk stops also for
// circular graphs.
func TestCircularGraph(t *testing.T) {
	graph := New()
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 1)
	info := make(map[int]int)
	graph.DoDepthFirst(func(vertex Vertex) error {
		if info[vertex.(int)] == 0 {
			info[vertex.(int)] = 1
		} else {
			t.Errorf("Edge already visited!")
		}
		return nil
	}, nil)
}
