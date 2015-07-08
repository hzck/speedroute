package model

import "container/heap"

type Graph struct {
	startNode *Node
	endNode   *Node
	edges     []*Edge
	prioQueue *PrioQueue
}

func (graph *Graph) StartNode() *Node {
	return graph.startNode
}

func (graph *Graph) AddStartNode(node *Node) {
	graph.startNode = node
}

func (graph *Graph) EndNode() *Node {
	return graph.endNode
}

func (graph *Graph) AddEndNode(node *Node) {
	graph.endNode = node
}

func (graph *Graph) Edges() []*Edge {
	return graph.edges
}

func (graph *Graph) AddEdge(edge *Edge) {
	graph.edges = append(graph.edges, edge)
	edge.From().AddFromEdge(edge)
	edge.To().AddToEdge(edge)
}

func (graph *Graph) PrioPath() *Path {
	if graph.prioQueue.Len() > 0 {
		return heap.Pop(graph.prioQueue).(*Path)
	}
	return nil
}

func (graph *Graph) AddPossiblePath(newPath *Path) {
	heap.Push(graph.prioQueue, newPath)
}

func CreateGraph() *Graph {
	graph := new(Graph)
	graph.prioQueue = &PrioQueue{}
	heap.Init(graph.prioQueue)
	return graph
}
