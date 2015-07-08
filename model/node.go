package model

import "math"

type Node struct {
	fromEdges   []*Edge
	toEdges     []*Edge
	id          string
	rewards     map[*Reward]int
	minPathLeft int
	revisitable bool
}

func (node *Node) FromEdges() []*Edge {
	return node.fromEdges
}

func (node *Node) AddFromEdge(edge *Edge) {
	node.fromEdges = append(node.fromEdges, edge)
}

func (node *Node) ToEdges() []*Edge {
	return node.toEdges
}

func (node *Node) AddToEdge(edge *Edge) {
	node.toEdges = append(node.toEdges, edge)
}

func (node *Node) Id() string {
	return node.id
}

func (node *Node) Rewards() map[*Reward]int {
	return node.rewards
}

func (node *Node) AddReward(reward *Reward, quantity int) {
	node.rewards[reward] = quantity
}

func (node *Node) MinPathLeft() int {
	return node.minPathLeft
}

func (node *Node) SetMinPathLeft(mpl int) {
	node.minPathLeft = mpl
}

func (node *Node) Revisitable() bool {
	return node.revisitable
}

func (node *Node) SetRevisitable(revisit bool) {
	node.revisitable = revisit
}

func CreateNode(id string) *Node {
	node := new(Node)
	node.id = id
	node.rewards = make(map[*Reward]int)
	node.minPathLeft = math.MaxInt32
	return node
}

type DijkstraPrio []*Node

func (dp DijkstraPrio) Len() int {
	return len(dp)
}

func (dp DijkstraPrio) Less(i, j int) bool {
	return dp[i].MinPathLeft() < dp[j].MinPathLeft()
}

func (dp DijkstraPrio) Swap(i, j int) {
	dp[i], dp[j] = dp[j], dp[i]
}

func (dp *DijkstraPrio) Push(x interface{}) {
	*dp = append(*dp, x.(*Node))
}

func (dp *DijkstraPrio) Pop() interface{} {
	old := *dp
	n := len(old)
	x := old[n-1]
	*dp = old[0 : n-1]
	return x
}
