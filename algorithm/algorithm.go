// Package algorithm calculates the shortest path in a directed, weighted graph with a set of requirements.
package algorithm

import (
	"container/heap"

	m "github.com/hzck/speedroute/model"
)

// Route takes a created graph object and finds the shortest path from start to end, returning
// that path as a list of edges.
func Route(graph *m.Graph) []*m.Edge {
	if graph.StartNode() == nil || graph.EndNode() == nil {
		return nil
	}
	addMinPathLeft(graph)
	startPath := m.CreatePath()
	startPath.AddRewards(graph.StartNode().Rewards())
	pq := &m.PrioQueue{}
	heap.Init(pq)
	addNodeEdgesToPrioQueue(pq, graph.StartNode(), startPath)
	for path := prioPath(pq); path != nil; path = prioPath(pq) {
		node := path.Edges()[len(path.Edges())-1].To()
		if node == graph.EndNode() {
			return path.Edges()
		}
		addNodeEdgesToPrioQueue(pq, node, path)
	}
	return nil
}

func addNodeEdgesToPrioQueue(pq *m.PrioQueue, node *m.Node, path *m.Path) {
	for _, edge := range node.FromEdges() {
		ok, i := path.PossibleRoute(edge)
		if ok {
			newPath := path.Copy()
			newPath.AddEdge(edge, i)
			heap.Push(pq, newPath)
		}
	}
}

func prioPath(pq *m.PrioQueue) *m.Path {
	if pq.Len() > 0 {
		return heap.Pop(pq).(*m.Path)
	}
	return nil
}

func addMinPathLeft(graph *m.Graph) {
	dp := &m.DijkstraPrio{}
	heap.Init(dp)
	visited := make(map[*m.Node]bool)
	endNode := graph.EndNode()
	endNode.SetMinPathLeft(0)
	visited[endNode] = true
	for _, edge := range endNode.ToEdges() {
		node := edge.From()
		node.SetMinPathLeft(edge.FastestTime())
		heap.Push(dp, node)
	}
	if dp.Len() > 0 {
		for node := heap.Pop(dp).(*m.Node); dp.Len() > 0; node = heap.Pop(dp).(*m.Node) {
			visited[node] = true
			for _, edge := range node.ToEdges() {
				innerNode := edge.From()
				if !visited[innerNode] {
					innerNode.SetMinPathLeft(edge.FastestTime() + node.MinPathLeft())
					heap.Push(dp, innerNode)
				}
			}
		}
	}
}
