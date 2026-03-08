// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

// Package directed provides support for creating and working with
// directed graphs.
package directed

import (
	"container/list"
	"fmt"
	"iter"
)

// Edge represents a directed edge as a pair [source, target].
type Edge[V comparable] [2]V

// Digraph is the respresentation of a directed graph. It contain all
// the edges and vertices of the graph.
type Digraph[V comparable] struct {
	edges       map[V]*list.List
	edgeCount   int
	vertexCount int
}

// find is used to locate an element in a list by value. It will
// return true and a pointer to the element if the element was found
// and false and a pointer to the last element of the list (or nil)
// otherwise.
func find[V comparable](lst *list.List, value V) (bool, *list.Element) {
	elem := lst.Front()
	if elem != nil {
		for elem != lst.Back() && elem.Value != value {
			elem = elem.Next()
		}
		if elem.Value == value {
			return true, elem
		}
	}
	return false, elem
}

// New will create a new, empty, directed graph.
func New[V comparable]() *Digraph[V] {
	return &Digraph[V]{edges: make(map[V]*list.List)}
}

// AddEdges adds edges to the graph. The source and target vertices
// will be added to the graph if they are not already present. Returns
// an error if any edge already exists.
func (graph *Digraph[V]) AddEdges(edges ...Edge[V]) error {
	for _, edge := range edges {
		graph.AddVertex(edge[0])
		graph.AddVertex(edge[1])
		lst := graph.edges[edge[0]]
		found, elem := find(lst, edge[1])
		if found {
			return fmt.Errorf("edge (%v, %v) already exists", edge[0], edge[1])
		}
		if elem != nil {
			lst.InsertAfter(edge[1], elem)
		} else {
			lst.PushBack(edge[1])
		}
		graph.edgeCount++
	}
	return nil
}

// RemoveEdge will remove and edge from the graph. The vertices that
// serve as endpoints for the edge will not be removed.  The method
// Returns an error if the edge does not exist.
func (graph *Digraph[V]) RemoveEdge(source, target V) error {
	lst := graph.edges[source]
	if found, elem := find(lst, target); found {
		lst.Remove(elem)
		graph.edgeCount--
		return nil
	}
	return fmt.Errorf("edge (%v, %v) does not exist", source, target)
}

// AddVertex will add a vertex to the graph. The vertex will have no
// in- or out-edges.  The function return 'true' if the vertex was
// successfully added, and 'false' if the vertex already existed.
func (graph *Digraph[V]) AddVertex(vertex V) bool {
	if graph.edges[vertex] == nil {
		graph.edges[vertex] = list.New()
		graph.vertexCount++
		return true
	}
	return false
}

// RemoveVertex will remove the vertex from the graph. Any edges
// connecting to the graph (either in- or out-edges) will also be
// removed.
func (graph *Digraph[V]) RemoveVertex(vertex V) bool {
	// It is guaranteed that the vertex is present in the map if
	// it is present in the graph.
	if graph.edges[vertex] != nil {
		// Remove the vertex from the map to remove all
		// out-edges.
		graph.edgeCount -= graph.edges[vertex].Len()
		delete(graph.edges, vertex)

		// Iterate over all the other lists to remove all
		// in-edges.
		for _, lst := range graph.edges {
			if found, elem := find(lst, vertex); found {
				lst.Remove(elem)
				graph.edgeCount--
			}
		}
		graph.vertexCount--
		return true
	}
	return false
}

// HasVertex check if a vertex exists in the graph. Will return 'true'
// if the vertex exists and 'false' otherwise.
func (graph *Digraph[V]) HasVertex(vertex V) bool {
	return graph.edges[vertex] != nil
}

// Adjacent check if an edge exists in the graph. Will return 'true' if
// the edge exists, and 'false' otherwise.
func (graph *Digraph[V]) Adjacent(source, target V) bool {
	if lst := graph.edges[source]; lst != nil {
		found, _ := find(lst, target)
		return found
	}
	return false
}

// Order will return the order of the graph, that is, the number of vertices
// in the graph.
func (graph *Digraph[V]) Order() int {
	return graph.vertexCount
}

// Size will return size of the graph, that is the number of edges in
// the graph.
func (graph *Digraph[V]) Size() int {
	return graph.edgeCount
}

// Vertices returns an iterator over all vertices in the graph.
func (graph *Digraph[V]) Vertices() iter.Seq[V] {
	return func(yield func(V) bool) {
		for vertex := range graph.edges {
			if !yield(vertex) {
				return
			}
		}
	}
}

// Edges returns an iterator over all edges in the graph as
// (source, target) pairs.
func (graph *Digraph[V]) Edges() iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		for vertex, edges := range graph.edges {
			for elem := edges.Front(); elem != nil; elem = elem.Next() {
				if !yield(vertex, elem.Value.(V)) {
					return
				}
			}
		}
	}
}

// OutEdges returns an iterator over the out-edges of a vertex as
// (source, target) pairs.
func (graph *Digraph[V]) OutEdges(vertex V) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		lst := graph.edges[vertex]
		if lst == nil {
			return
		}
		for elem := lst.Front(); elem != nil; elem = elem.Next() {
			if !yield(vertex, elem.Value.(V)) {
				return
			}
		}
	}
}
