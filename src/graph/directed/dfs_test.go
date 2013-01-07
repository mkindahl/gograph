package directed

import "testing"
import "fmt"

// Structure holding discovery-finish information for this test.
type Info struct {
	finish, discover int
}

// Test if x is nested inside y
func (x *Info) IsNestedIn(y *Info) bool {
	return y.discover < x.discover && x.discover < x.finish && x.finish < y.finish
}

// Test if x is disjoint with y
func (x *Info) IsDisjoint(y *Info) bool {
	return x.discover < x.finish && x.finish < y.discover && y.discover < y.finish ||
		y.discover < y.finish && y.finish < x.discover && x.discover < x.finish
}

// Test that the discovery-finish information satisfies the
// parantheses theorem
func IsValidNesting(x, y *Info) bool {
	return x.IsDisjoint(y) || x.IsNestedIn(y) || y.IsNestedIn(x)
}


func PrettyMap(input map[int]*Info) string {
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
	graph.AddEdge(1,2)
	graph.AddEdge(1,4)
	graph.AddEdge(2,5)
	graph.AddEdge(3,5)
	graph.AddEdge(3,6)
	graph.AddEdge(4,2)
	graph.AddEdge(5,4)
	graph.AddEdge(6,6)

	info := make(map[int]*Info)
	onDiscover := func (vertex Vertex, time int) {
		if info[vertex.(int)] == nil {
			info[vertex.(int)] = new(Info)
		}
		info[vertex.(int)].discover = time
	}
	onFinish := func (vertex Vertex, time int) {
		if info[vertex.(int)] == nil {
			info[vertex.(int)] = new(Info)
		}
		info[vertex.(int)].finish = time
	}
	graph.DoDepthFirst(onDiscover, onFinish)
	graph.DoEdges(func (source, target Vertex) {
		if source != target && !IsValidNesting(info[source.(int)], info[target.(int)]) {
			pretty := PrettyMap(info)
			t.Errorf("Edge %v -> %v has bad finish time (%s)\n", source, target, pretty)
		}
	})
			
}

func TestTopologicalWalk(t *testing.T) {
	graph := New()
	graph.AddEdge(1,2)
	graph.AddEdge(2,3)
	graph.AddEdge(3,4)
	time := 1
	graph.DoTopological(func (vertex Vertex) {
		if vertex.(int) != time {
			t.Errorf("Vertex %d processed at time %d\n", vertex.(int), time)
		}
		time++
	})

	graph = New()
	graph.AddEdge(1,2)
	graph.AddEdge(1,3)
	graph.AddEdge(2,4)
	graph.AddEdge(3,4)
	when := new([5]int)
	time = 1
	graph.DoTopological(func (vertex Vertex) {
		when[time] = vertex.(int)
		time++
	})
	if !(when[1] < when[2] && when[1] < when[3] && when[2] < when[4] && when[3] < when[4]) {
		t.Errorf("Not in topological order %v\n", when)
	}
}
