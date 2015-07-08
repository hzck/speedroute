package algorithm

import (
	"container/heap"
	m "github.com/hzck/speedroute/model"
)

func Route(graph *m.Graph) []*m.Edge {
	if graph.StartNode() == nil || graph.EndNode() == nil || len(graph.Edges()) == 0 {
		return nil
	}
	addMinPathLeft(graph)
	startPath := m.CreatePath()
	startPath.AddRewards(graph.StartNode().Rewards())
	addNodeEdgesToPrioQueue(graph, graph.StartNode(), startPath)
	for path := graph.PrioPath(); path != nil; path = graph.PrioPath() {
		node := path.Edges()[len(path.Edges())-1].To()
		if node == graph.EndNode() {
			return path.Edges()
		}
		addNodeEdgesToPrioQueue(graph, node, path)
	}
	return nil
}

func addNodeEdgesToPrioQueue(graph *m.Graph, node *m.Node, path *m.Path) {
	for _, edge := range node.FromEdges() {
		ok, i := path.PossibleRoute(edge)
		if ok {
			newPath := path.Copy()
			newPath.AddEdge(edge, i)
			graph.AddPossiblePath(newPath)
		}
	}
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
