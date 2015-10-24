// Package model provides model structs for a directed weighted graph.
package model

// Graph includes pointers to startNode, endNode and the internal type prioQueue.
type Graph struct {
	startNode *Node
	endNode   *Node
}

// StartNode returns startNode pointer.
func (graph *Graph) StartNode() *Node {
	return graph.startNode
}

// EndNode returns endNode pointer.
func (graph *Graph) EndNode() *Node {
	return graph.endNode
}

// CreateGraph is the constructor for graph taking pointers to startNode and endNode as parameters.
func CreateGraph(start, end *Node) *Graph {
	graph := new(Graph)
	graph.startNode = start
	graph.endNode = end
	return graph
}
