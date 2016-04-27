// Package parser manages JSON/XML object and converts them into a graph.
package parser

import (
	"encoding/json"
	"fmt"
	m "github.com/hzck/speedroute/model"
	"io/ioutil"
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
}

type node struct {
	ID          string      `json:"id"`
	Rewards     []rewardRef `json:"rewards"`
	Revisitable bool        `json:"revisitable"`
}

type weight struct {
	Time         *int        `json:"time"`
	Requirements []rewardRef `json:"requirements"`
}

type edge struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Weights []weight `json:"weights"`
}

// CreateGraphFromFile takes a path as a parameter and creates rewards, nodes and edges before
// returning a pointer to a graph.
func CreateGraphFromFile(path string) *m.Graph { //throw error
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("got error: " + path)
		//add error handling here
	}
	var g graph
	err = json.Unmarshal(file, &g)

	if err != nil {
		fmt.Println("got error: ", err)
		//add error handling here
	}
	rewards := make(map[string]*m.Reward)
	for _, r := range g.Rewards {
		rewards[r.ID] = m.CreateReward(r.ID, r.Unique)
	}
	nodes := make(map[string]*m.Node)
	for _, n := range g.Nodes {
		node := m.CreateNode(n.ID, n.Revisitable)
		for _, reward := range n.Rewards {
			node.AddReward(rewards[reward.RewardID], getPointerValueOrOne(reward.Quantity)) //duplicate code
		}
		nodes[node.ID()] = node
	}
	for _, e := range g.Edges {
		edge := m.CreateEdge(nodes[e.From], nodes[e.To])
		for _, w := range e.Weights {
			weight := m.CreateWeight(getPointerValueOrOne(w.Time))
			for _, reward := range w.Requirements {
				weight.AddRequirement(rewards[reward.RewardID], getPointerValueOrOne(reward.Quantity)) //duplicate code
			}
			edge.AddWeight(weight)
		}
	}
	return m.CreateGraph(nodes[g.StartID], nodes[g.EndID])
}

func getPointerValueOrOne(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 1
}

func CreateJSONFromRoutedPath(path []*m.Edge) ([]byte, error) {
	if path == nil || len(path) == 0 {
		return json.Marshal(path)
	}
	result := make([]string, len(path)+1)
	result[0] = path[0].From().ID()
	for i := 1; i < len(result); i++ {
		result[i] = path[i-1].To().ID()
	}
	return json.Marshal(result)
}
