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
		reward := m.CreateReward(r.Id)
		reward.SetUnique(r.Unique)
		rewards[reward.Id()] = reward
	}
	nodes := make(map[string]*m.Node)
	for _, n := range g.Nodes {
		node := m.CreateNode(n.Id)
		node.SetRevisitable(n.Revisitable)
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
		graph.AddEdge(edge)
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

/*
func addRewardFromJSON(rewards map[string]*m.Reward, dec *json.Decoder) {
	var r reward
	dec.Decode(&r)
	reward := m.CreateReward(r.id)
	reward.SetUnique(r.unique)
	rewards[reward.Id()] = reward
}

func addNodeFromJSON(nodes map[string]*m.Node, rewards map[string]*m.Reward, dec *json.Decoder) {
	var n node
	dec.Decode(&n)
	node := m.CreateNode(n.id)
	node.SetRevisitable(n.revisitable)
	for _, reward := range n.rewards {
		node.AddReward(rewards[reward.rewardId], getPointerValueOrOne(reward.quantity))
	}
	nodes[node.Id()] = node
}

func addEdgeFromJSON(graph *m.Graph, nodes map[string]*m.Node, rewards map[string]*m.Reward, dec *json.Decoder) {
	var e edge
	dec.Decode(&e)
	edge := m.CreateEdge(nodes[e.from], nodes[e.to])
	for _, w := range e.weights {
		weight := m.CreateWeight(getPointerValueOrOne(w.time))
		for _, reward := range w.requirements {
			weight.AddRequirement(rewards[reward.rewardId], getPointerValueOrOne(reward.quantity))
		}
		edge.AddWeight(weight)
	}
	graph.AddEdge(edge)
}

func addGraphFromJSON(g *m.Graph, nodes map[string]*m.Node, dec *json.Decoder) {
	var graphModel graph
	dec.Decode(&graphModel)
	g.AddStartNode(nodes[graphModel.startId])
	g.AddEndNode(nodes[graphModel.endId])
}

func CreateGraphFromFile(path string) *m.Graph {
	file, _ := os.Open(path)
	defer file.Close()
	graph := m.CreateGraph()
	rewards := make(map[string]*m.Reward)
	nodes := make(map[string]*m.Node)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Reward") {
			addRewardFromJSON(rewards, getDecoder(line[6:]))
		} else if strings.HasPrefix(line, "Node") {
			addNodeFromJSON(nodes, rewards, getDecoder(line[4:]))
		} else if strings.HasPrefix(line, "Edge") {
			addEdgeFromJSON(graph, nodes, rewards, getDecoder(line[4:]))
		} else if strings.HasPrefix(line, "Graph") {
			addGraphFromJSON(graph, nodes, getDecoder(line[5:]))
		}
	}
	return graph
}

func getPointerValueOrOne(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 1
}

func getDecoder(data string) *json.Decoder {
	return json.NewDecoder(strings.NewReader(data))
}
*/
