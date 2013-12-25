// Copyright (c) 2013, Mats Kindahl. All rights reserved.
//
// Use of this source code is governed by a BSD license that can be
// found in the README file.

// Package directed provides support for creating and working with
// directed graphs.
package directed

import "container/list"

//import "fmt"

// Vertex is a convenience declaration for a vertex of the
// graph. There are currently no expectations on the vertices of a
// graph: any object that can be used as key in a map can be used.
type Vertex interface{}

type adjacencyList map[Vertex]*list.List

// Graph is the respresentation of a directed graph. It contain all
// the edges and vertices of the graph.
type Graph struct {
	edges                  adjacencyList
	edgeCount, vertexCount int
}

// find is used to locate an element in a list by value. It will
// return true and a pointer to the element if the element was found
// and false and a pointer to the last element of the list (or nil)
// otherwise.
func find(lst *list.List, value Vertex) (bool, *list.Element) {
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
func New() *Graph {
	return &Graph{edges: make(adjacencyList)}
}

// AddEdge add an edge to the graph. The source and target vertices
// will be added to the graph if they are not already present. The
// function return 'true' if the edge was successfully added, and
// 'false' if the edge already existed.
func (graph *Graph) AddEdge(source, target Vertex) bool {
	graph.AddVertex(source)
	graph.AddVertex(target)
	lst := graph.edges[source]
	found, elem := find(lst, target)
	if !found {
		if elem != nil {
			lst.InsertAfter(target, elem)
		} else {
			lst.PushBack(target)
		}
		graph.edgeCount++
		//		fmt.Printf("Adding %d edges, giving %d\n", 1, graph.edgeCount)
	}
	return !found
}

// RemoveEdge will remove and edge from the graph. The vertices that
// serve as endpoints for the edge will not be removed.  The method
// returns 'true' if the edge was successfully removed, 'false'
// otherwise.
func (graph *Graph) RemoveEdge(source, target Vertex) bool {
	lst := graph.edges[source]
	if found, elem := find(lst, target); found {
		lst.Remove(elem)
		graph.edgeCount--
		//		fmt.Printf("Removing %d edges, giving %d\n", 1, graph.edgeCount)
		return true
	}
	return false
}

// AddVertex will add a vertex to the graph. The vertex will have no
// in- or out-edges.  The function return 'true' if the vertex was
// successfully added, and 'false' if the vertex already existed.
func (graph *Graph) AddVertex(vertex Vertex) bool {
	if graph.edges[vertex] == nil {
		graph.edges[vertex] = list.New()
		graph.vertexCount++
		//		fmt.Printf("Adding %d vertex, giving %d\n", 1, graph.vertexCount)
		return true
	}
	return false
}

// RemoveVertex will remove the vertex from the graph. Any edges
// connecting to the graph (either in- or out-edges) will also be
// removed.
func (graph *Graph) RemoveVertex(vertex Vertex) bool {
	// It is guaranteed that the vertex is present in the map if
	// it is present in the graph.
	if graph.edges[vertex] != nil {
		// Remove the vertex from the map to remove all
		// out-edges.
		graph.edgeCount -= graph.edges[vertex].Len()
		// fmt.Printf("Removing %d edges, giving %d\n",
		//            graph.edges[vertex].Len(), graph.edgeCount)
		delete(graph.edges, vertex)

		// Iterate over all the other lists to remove all
		// in-edges.
		for _, lst := range graph.edges {
			if found, elem := find(lst, vertex); found {
				lst.Remove(elem)
				graph.edgeCount--
				//				fmt.Printf("Removing %d edges, giving %d\n", 1, graph.edgeCount)
			}
		}
		graph.vertexCount--
		//		fmt.Printf("Removing %d vertex, giving %d\n", 1, graph.vertexCount)
		return true
	}
	return false
}

// HasVertex check if a vertex exists in the graph. Will return 'true'
// if the vertex exists and 'false' otherwise.
func (graph *Graph) HasVertex(vertex Vertex) bool {
	return graph.edges[vertex] != nil
}

// HasEdge check if an edge exists in the graph. Will return 'true' if
// the edge exists, and 'false' otherwise.
func (graph *Graph) HasEdge(source, target Vertex) bool {
	if lst := graph.edges[source]; lst != nil {
		found, _ := find(lst, target)
		return found
	}
	return false
}

// Order will return the order of the graph, that is, the number of vertices
// in the graph.
func (graph *Graph) Order() int {
	return graph.vertexCount
}

// Size will return size of the graph, that is the number of edges in
// the graph.
func (graph *Graph) Size() int {
	return graph.edgeCount
}

// VertexWalkFunc is a function called when walking vertices of a
// graph. Parent can either be a vertex, or n
type VertexWalkFunc func(vertex Vertex) error

// DoVertices iterate over all the vertices of the graph calling
// 'walkFn' with each vertex. If the walk function returns an error,
// iteration will be aborted and the error returned to the caller.
func (graph *Graph) DoVertices(walkFn VertexWalkFunc) error {
	for vertex := range graph.edges {
		if err := walkFn(vertex); err != nil {
			return err
		}
	}
	return nil
}

// EdgeWalkFunc is a function called when walking edges of a graph.
type EdgeWalkFunc func(source, target Vertex) error

// DoEdges will iterate over all the edges of the graph calling
// 'walkFn' with the source and target vertex of the edge. If the walk
// function return an error, iteration will be aborted and the error
// returned.
func (graph *Graph) DoEdges(walkFn EdgeWalkFunc) error {
	for vertex, edges := range graph.edges {
		for elem := edges.Front(); elem != nil; elem = elem.Next() {
			if err := walkFn(vertex, elem.Value); err != nil {
				return err
			}
		}
	}
	return nil
}

// DoOutEdges iterate over the out-edges of a vertex, calling 'walkFn'
// with the source and the target vertex of the edge.  The source
// target will be 'vertex' in each case, but edge walk functions use
// this common format. If the walk function return an error, iteration
// will be aborted and the error returned.
func (graph *Graph) DoOutEdges(vertex Vertex, walkFn EdgeWalkFunc) error {
	lst := graph.edges[vertex]
	if lst == nil {
		return nil
	}
	for elem := lst.Front(); elem != nil; elem = elem.Next() {
		if err := walkFn(vertex, elem.Value); err != nil {
			return err
		}
	}
	return nil
}
