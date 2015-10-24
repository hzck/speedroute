package model

import "math"

// Node stores all information about nodes like neighboring edges, rewards and if it's revisitable.
type Node struct {
	fromEdges   []*Edge
	toEdges     []*Edge
	id          string
	rewards     map[*Reward]int
	minPathLeft int
	revisitable bool
}

// FromEdges returns all edges which has this node as a from node.
func (node *Node) FromEdges() []*Edge {
	return node.fromEdges
}

//AddFromEdge adds an edge which has this node as a from edge.
func (node *Node) AddFromEdge(edge *Edge) {
	node.fromEdges = append(node.fromEdges, edge)
}

// ToEdges returns all edges which has this node as a to node.
func (node *Node) ToEdges() []*Edge {
	return node.toEdges
}

//AddToEdge adds an edge which has this node as a to edge.
func (node *Node) AddToEdge(edge *Edge) {
	node.toEdges = append(node.toEdges, edge)
}

// ID returns node id.
func (node *Node) ID() string {
	return node.id
}

// Rewards returns the node rewards as a map with quantity.
func (node *Node) Rewards() map[*Reward]int {
	return node.rewards
}

// AddReward sets a node reward with quantity.
func (node *Node) AddReward(reward *Reward, quantity int) {
	node.rewards[reward] = quantity
}

//MinPathLeft returns minimum path left to end node or math.MaxInt32 if not known.
func (node *Node) MinPathLeft() int {
	return node.minPathLeft
}

// SetMinPathLeft sets the minimum path left to end node.
func (node *Node) SetMinPathLeft(mpl int) {
	node.minPathLeft = mpl
}

// Revisitable returns true if the node is revisitable
func (node *Node) Revisitable() bool {
	return node.revisitable
}

//CreateNode takes id and revisitable as paramteres, returning a pointer to the new node
func CreateNode(id string, revisit bool) *Node {
	node := new(Node)
	node.id = id
	node.revisitable = revisit
	node.rewards = make(map[*Reward]int)
	node.minPathLeft = math.MaxInt32
	return node
}

// DijkstraPrio is a list of node pointers and implements heap interface
type DijkstraPrio []*Node

// Len returns length of DijkstraPrio
func (dp DijkstraPrio) Len() int {
	return len(dp)
}

// Less checks which node has shorter minimum path left to end node
func (dp DijkstraPrio) Less(i, j int) bool {
	return dp[i].MinPathLeft() < dp[j].MinPathLeft()
}

// Swap switches places in the DijkstraPrio list
func (dp DijkstraPrio) Swap(i, j int) {
	dp[i], dp[j] = dp[j], dp[i]
}

// Push adds a node into DijkstraPrio list to the correct place
func (dp *DijkstraPrio) Push(x interface{}) {
	*dp = append(*dp, x.(*Node))
}

// Pop returns the node with the shortest minimum path left to end node
func (dp *DijkstraPrio) Pop() interface{} {
	old := *dp
	n := len(old)
	x := old[n-1]
	*dp = old[0 : n-1]
	return x
}
