package model

import (
	"fmt"
	"testing"
)

func TestNodeFromEdges(t *testing.T) {
	graph := new(Graph)
	startNode := new(Node)
	endNode := new(Node)
	CreateEdge(startNode, endNode)
	graph.AddStartNode(startNode)
	graph.AddEndNode(endNode)
	if len(startNode.FromEdges()) != 1 {
		t.Fail()
		fmt.Println("No start node from edge")
	}
	if len(endNode.ToEdges()) != 1 {
		t.Fail()
		fmt.Println("No end node to edge")
	}
}
