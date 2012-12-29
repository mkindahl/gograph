package directed

// Basic functionality for creating and manipulating directed graphs

import "container/list"

type Vertex interface{}
type AdjecencyList map[Vertex]*list.List
type DirectedGraph struct {
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

func New() *DirectedGraph {
	self := &DirectedGraph{edges: make(AdjecencyList)}
	return self
}

func (graph *DirectedGraph) AddEdge(source, target Vertex) bool {
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

func (graph *DirectedGraph) AddVertex(vertex Vertex) bool {
	if (graph.edges[vertex] != nil) {
		graph.edges[vertex] = nil
		return true
	}
	return false
}

func (graph *DirectedGraph) HasVertex(vertex Vertex) bool {
	return graph.edges[vertex] != nil
}

func (graph *DirectedGraph) HasEdge(source, target Vertex) bool {
	lst := graph.edges[source]
	if lst != nil {
		found, _ := Find(lst, target)
		return found
	}
	return false
}


