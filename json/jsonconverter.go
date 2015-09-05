package json

import (
	"encoding/json"
	"fmt"
	m "github.com/hzck/speedroute/model"
	"io/ioutil"
)

type reward struct {
	Id     string
	Unique bool
}

type rewardRef struct {
	RewardId string
	Quantity *int
}

type node struct {
	Id          string
	Rewards     []rewardRef
	Revisitable bool
}

type weight struct {
	Time         *int
	Requirements []rewardRef
}

type edge struct {
	From, To string
	Weights  []weight
}

type graph struct {
	Rewards        []reward
	Nodes          []node
	Edges          []edge
	StartId, EndId string
}

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
		reward := m.CreateReward(r.Id, r.Unique)
		rewards[reward.Id()] = reward
	}
	nodes := make(map[string]*m.Node)
	for _, n := range g.Nodes {
		node := m.CreateNode(n.Id, n.Revisitable)
		for _, reward := range n.Rewards {
			node.AddReward(rewards[reward.RewardId], getPointerValueOrOne(reward.Quantity)) //duplicate code
		}
		nodes[node.Id()] = node
	}
	graph := m.CreateGraph()
	for _, e := range g.Edges {
		edge := m.CreateEdge(nodes[e.From], nodes[e.To])
		for _, w := range e.Weights {
			weight := m.CreateWeight(getPointerValueOrOne(w.Time))
			for _, reward := range w.Requirements {
				weight.AddRequirement(rewards[reward.RewardId], getPointerValueOrOne(reward.Quantity)) //duplicate code
			}
			edge.AddWeight(weight)
		}
	}
	graph.AddStartNode(nodes[g.StartId])
	graph.AddEndNode(nodes[g.EndId])
	return graph
}

func getPointerValueOrOne(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 1
}
