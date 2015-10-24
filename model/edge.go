package model

import "sort"

// Edge contains pointers to from and to nodes as well as a list of weights.
type Edge struct {
	from    *Node
	to      *Node
	weights []*Weight
}

// From returns a pointer for the edge's from node.
func (edge *Edge) From() *Node {
	return edge.from
}

// To returns a pointer for the edge's to node.
func (edge *Edge) To() *Node {
	return edge.to
}

// Weights returns a list of the edge's weights or a default weight if no weights are set.
func (edge *Edge) Weights() []*Weight {
	if len(edge.weights) == 0 {
		return []*Weight{CreateWeight(1)}
	}
	return edge.weights
}

// AddWeight adds a weight to the edge.
func (edge *Edge) AddWeight(weight *Weight) {
	edge.weights = append(edge.weights, weight)
	sort.Sort(ByTime(edge.weights))
}

// CreateEdge takes from and to node pointers and returns a pointer to a new edge.
func CreateEdge(from, to *Node) *Edge {
	edge := new(Edge)
	edge.from = from
	from.AddFromEdge(edge)
	edge.to = to
	to.AddToEdge(edge)
	return edge
}

// FastestTime returns the fastest possible time this edge's weighs can provide.
func (edge *Edge) FastestTime() int {
	if len(edge.weights) == 0 {
		return 1
	}
	return edge.weights[0].Time()
}
