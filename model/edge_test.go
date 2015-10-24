package model

import (
	"fmt"
	"testing"
)

// TestNodeFromEdges tests that when creating an edge, from and to node gets correctly populated.
func TestNodeFromEdges(t *testing.T) {
	startNode := new(Node)
	endNode := new(Node)
	CreateEdge(startNode, endNode)
	if len(startNode.FromEdges()) != 1 {
		t.Fail()
		fmt.Println("No start node from edge")
	}
	if len(endNode.ToEdges()) != 1 {
		t.Fail()
		fmt.Println("No end node to edge")
	}
}
