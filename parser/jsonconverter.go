// Package parser manages JSON/XML object and converts them into a graph.
package parser

import (
	"encoding/json"
	"io/ioutil"

	m "github.com/hzck/speedroute/model"
)

type rewardRef struct {
	RewardID string `json:"rewardId"`
	Quantity *int   `json:"quantity"`
}

type graph struct {
	Rewards []reward `json:"rewards"`
	Nodes   []node   `json:"nodes"`
	Edges   []edge   `json:"edges"`
	StartID string   `json:"startId"`
	EndID   string   `json:"endId"`
}

type reward struct {
	ID     string `json:"id"`
	Unique bool   `json:"unique"`
	IsA    string `json:"isA"`
}

type node struct {
	ID          string      `json:"id"`
	Rewards     []rewardRef `json:"rewards"`
	Revisitable bool        `json:"revisitable"`
}

type weight struct {
	Time         string      `json:"time"`
	Requirements []rewardRef `json:"requirements"`
}

type edge struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Weights []weight `json:"weights"`
}

// CreateGraphFromFile takes a path as a parameter and creates rewards, nodes and edges before
// returning a pointer to a graph.
func CreateGraphFromFile(path string) (*m.Graph, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var g graph
	err = json.Unmarshal(file, &g)
	if err != nil {
		return nil, err
	}

	rewards := make(map[string]*m.Reward)
	// ugly loop to make sure to handle different ordered rewards
	for rewardAdded := true; rewardAdded; {
		rewardAdded = false
		for _, r := range g.Rewards {
			if rewards[r.ID] == nil && (r.IsA == "" || rewards[r.IsA] != nil) {
				rewards[r.ID] = m.CreateReward(r.ID, r.Unique, rewards[r.IsA])
				rewardAdded = true
			}
		}
	}

	nodes := make(map[string]*m.Node)
	for _, n := range g.Nodes {
		node := m.CreateNode(n.ID, n.Revisitable)
		for _, rewardRef := range n.Rewards {
			// duplicate code
			node.AddReward(rewards[rewardRef.RewardID], getPointerValueOrOne(rewardRef.Quantity))
		}
		nodes[node.ID()] = node
	}

	for _, e := range g.Edges {
		edge := m.CreateEdge(nodes[e.From], nodes[e.To])
		for _, w := range e.Weights {
			time, err := parseTime(w.Time)
			if err != nil {
				return nil, err
			}
			weight := m.CreateWeight(time)
			for _, rewardRef := range w.Requirements {
				// duplicate code
				weight.AddRequirement(rewards[rewardRef.RewardID], getPointerValueOrOne(rewardRef.Quantity))
			}
			edge.AddWeight(weight)
		}
	}
	return m.CreateGraph(nodes[g.StartID], nodes[g.EndID]), nil
}

func getPointerValueOrOne(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 1
}

// CreateJSONFromRoutedPath takes an array of edges and creates an array of the included nodes,
// and marshals it as json data in a byte array.
func CreateJSONFromRoutedPath(path []*m.Edge) ([]byte, error) {
	if len(path) == 0 {
		return json.Marshal(path)
	}
	result := make([]string, len(path)+1)
	result[0] = path[0].From().ID()
	for i := 1; i < len(result); i++ {
		result[i] = path[i-1].To().ID()
	}
	return json.Marshal(result)
}
