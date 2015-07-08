package model

import "testing"

func TestNodeFromEdges(t *testing.T) {
	graph := new(Graph)
	startNode := new(Node)
	endNode := new(Node)
	graph.AddStartNode(startNode)
	graph.AddEndNode(endNode)
	graph.AddEdge(CreateEdge(startNode, endNode))
	if len(startNode.FromEdges()) != 1 {
		t.Fail()
	}
	if len(endNode.ToEdges()) != 1 {
		t.Fail()
	}
}
