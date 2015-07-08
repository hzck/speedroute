package model

import "sort"

type Edge struct {
	from    *Node
	to      *Node
	weights []*Weight
}

func (edge *Edge) From() *Node {
	return edge.from
}

func (edge *Edge) To() *Node {
	return edge.to
}

func (edge *Edge) Weights() []*Weight {
	if len(edge.weights) == 0 {
		return []*Weight{CreateWeight(1)}
	}
	return edge.weights
}

func (edge *Edge) AddWeight(weight *Weight) {
	edge.weights = append(edge.weights, weight)
	sort.Sort(ByTime(edge.weights))
}

func CreateEdge(from, to *Node) *Edge {
	edge := new(Edge)
	edge.from = from
	edge.to = to
	return edge
}

func (edge *Edge) FastestTime() int {
	if len(edge.weights) == 0 {
		return 1
	}
	return edge.weights[0].Time()
}
