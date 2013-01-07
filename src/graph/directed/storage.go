package directed

// Basic functionality for creating and manipulating directed graphs

import "container/list"

type Vertex interface{}
type AdjecencyList map[Vertex]*list.List
type Graph struct {
	edges AdjecencyList
}

type Edge struct {
	source, target Vertex
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

func New() *Graph {
	self := &Graph{edges: make(AdjecencyList)}
	return self
}

func (graph *Graph) AddEdge(source, target Vertex) bool {
	lst := graph.edges[source]
	if (lst == nil) {
		lst = list.New()
		graph.edges[source] = lst
	}
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

func (graph *Graph) AddVertex(vertex Vertex) bool {
	if (graph.edges[vertex] != nil) {
		graph.edges[vertex] = list.New()
		return true
	}
	return false
}

func (graph *Graph) HasVertex(vertex Vertex) bool {
	return graph.edges[vertex] != nil
}

func (graph *Graph) DoEdges(visitor (func (source, target Vertex))) {
	for vertex := range graph.edges {
		graph.DoOutEdges(vertex, visitor)
	}
}

func (graph *Graph) DoOutEdges(vertex Vertex, visitor (func (source, target Vertex))) {
	lst := graph.edges[vertex]
	if lst == nil {
		return
	}
	for elem := lst.Front() ; elem != nil ; elem = elem.Next() {
		visitor(vertex, elem.Value)
	}
}

func (graph *Graph) HasEdge(source, target Vertex) bool {
	lst := graph.edges[source]
	if lst != nil {
		found, _ := Find(lst, target)
		return found
	}
	return false
}


