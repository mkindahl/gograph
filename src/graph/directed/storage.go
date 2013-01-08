package directed

// Basic functionality for creating and manipulating directed graphs

import "container/list"

type Vertex interface{}
type AdjecencyList map[Vertex]*list.List
type Graph struct {
	edges AdjecencyList
}

func Find(lst *list.List, value Vertex) (bool, *list.Element) {
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

// Create a new, empty, directed graph.
func New() *Graph {
	return &Graph{edges: make(AdjecencyList)}
}

// Add an edge to the graph. The source and target vertices will be
// added to the graph if they are not already present. The function
// return 'true' if the edge was successfully added, and 'false' if
// the edge already existed.
func (graph *Graph) AddEdge(source, target Vertex) bool {
	graph.AddVertex(source)
	graph.AddVertex(target)
	lst := graph.edges[source]
	found, elem := Find(lst, target)
	if !found {
		if elem != nil {
			lst.InsertAfter(target, elem)
		} else {	
			lst.PushBack(target)
		}
	}
	return !found
}

// Add a vertex to the graph. The vertex will have no in- or
// out-edges.  The function return 'true' if the vertex was
// successfully added, and 'false' if the vertex already existed.
func (graph *Graph) AddVertex(vertex Vertex) bool {
	if (graph.edges[vertex] == nil) {
		graph.edges[vertex] = list.New()
		return true
	}
	return false
}

// Iterate over all the vertices of the graph calling 'visitor' with
// each vertex.
func (graph *Graph) DoVertices(visitor (func (vertex Vertex))) {
	for vertex := range graph.edges {
		visitor(vertex)
	}
}

// Iterate over all the edges of the graph calling 'visitor' with the
// source and target vertex of the edge.
func (graph *Graph) DoEdges(visitor (func (source, target Vertex))) {
	for vertex, edges := range graph.edges {
		for elem := edges.Front() ; elem != nil ; elem = elem.Next() {
			visitor(vertex, elem.Value)
		}
	}
}

// Iterate over the out-edges of a vertex, calling 'visitor' with the
// source and the target vertex of the edge.  The source target will
// be 'vertex' in each case, but edge visitor functions use this
// common format.
func (graph *Graph) DoOutEdges(vertex Vertex, visitor (func (source, target Vertex))) {
	lst := graph.edges[vertex]
	if lst == nil {
		return
	}
	for elem := lst.Front() ; elem != nil ; elem = elem.Next() {
		visitor(vertex, elem.Value)
	}
}

// Check if a vertex exists in the graph. Will return 'true' if the
// vertex exists, and 'false' if it does not exist in the graph.
func (graph *Graph) HasVertex(vertex Vertex) bool {
	return graph.edges[vertex] != nil
}

// Check if an edge exists in the graph. Will return 'true' if the
// edge exists, and 'false' if it does not exist in the graph.
func (graph *Graph) HasEdge(source, target Vertex) bool {
	lst := graph.edges[source]
	if lst != nil {
		found, _ := Find(lst, target)
		return found
	}
	return false
}


