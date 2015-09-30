// Package json converts a JSON object into a graph.
package json

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
	Rewards []struct {
		ID     string `json:"id"`
		Unique bool   `json:"unique"`
	} `json:"rewards"`
	Nodes []struct {
		ID          string      `json:"id"`
		Rewards     []rewardRef `json:"rewards"`
		Revisitable bool        `json:"revisitable"`
	} `json:"nodes"`
	Edges []struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Weights []struct {
			Time         *int        `json:"time"`
			Requirements []rewardRef `json:"requirements"`
		} `json:"weights"`
	} `json:"edges"`
	StartID string `json:"startId"`
	EndID   string `json:"endId"`
}

// CreateGraphFrom file takes a path as a parameter and creates rewards, nodes and edges before
// returning a pointer to a graph
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
