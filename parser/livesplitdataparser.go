package parser

import (
	"encoding/json"
	"encoding/xml"
)

type segments struct {
	Segments []struct {
		Name            string `xml:"Name"`
		BestSegmentTime string `xml:"BestSegmentTime>RealTime"`
	} `xml:"Segments>Segment"`
}

// LivesplitXMLtoJSON takes an XML input from a LiveSplit .ltt file and creates the route in JSON.
func LivesplitXMLtoJSON(data string) (string, error) {
	var v segments
	if err := xml.Unmarshal([]byte(data), &v); err != nil {
		return "", err
	}
	var g graph
	currentNode := "START"
	g.Rewards = make([]reward, 0)
	g.Nodes = append(g.Nodes, node{ID: currentNode, Rewards: make([]rewardRef, 0)})
	g.StartID = currentNode
	for _, seg := range v.Segments {
		g.Nodes = append(g.Nodes, node{ID: seg.Name, Rewards: make([]rewardRef, 0)})
		g.Edges = append(g.Edges, edge{From: currentNode, To: seg.Name, Weights: []weight{weight{Time: seg.BestSegmentTime, Requirements: make([]rewardRef, 0)}}})
		currentNode = seg.Name
	}
	g.EndID = currentNode
	result, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return "", err
	}
	return string(result), nil
}
