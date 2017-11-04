package model

// Path holds information about a possible path, including edges, current length and current rewards.
type Path struct {
	edges   []*Edge
	length  int
	rewards map[*Reward]int
}

// Edges returns a list of path edges.
func (path *Path) Edges() []*Edge {
	return path.edges
}

// AddEdge adds an edge with length and its to node's rewards to the path.
func (path *Path) AddEdge(edge *Edge, i int) {
	path.edges = append(path.edges, edge)
	path.length += edge.Weights()[i].Time()
	path.AddRewards(edge.To().Rewards())
}

// Length returns the path length.
func (path *Path) Length() int {
	return path.length
}

// Rewards returns the path rewards.
func (path *Path) Rewards() map[*Reward]int {
	return path.rewards
}

//AddRewards adds a map of rewards to the path.
func (path *Path) AddRewards(rewards map[*Reward]int) {
	for key, value := range rewards {
		path.rewards[key] += value
	}
}

// CreatePath is the constructor for path pointer.
func CreatePath() *Path {
	path := new(Path)
	path.rewards = make(map[*Reward]int)
	return path
}

// IsLongerThan compares two paths to see which has best potential to be the shortest path.
func (path *Path) IsLongerThan(other *Path) bool {
	pLength := path.Length() + path.Edges()[len(path.Edges())-1].To().MinPathLeft()
	oLength := other.Length() + other.Edges()[len(other.Edges())-1].To().MinPathLeft()
	if pLength != oLength {
		return pLength > oLength
	}
	return len(path.Edges()) > len(other.Edges())
}

// Copy copies a path values into another path.
func (path *Path) Copy() *Path {
	pathCopy := CreatePath()
	pathCopy.edges = make([]*Edge, len(path.Edges()))
	copy(pathCopy.edges, path.Edges())
	pathCopy.length = path.Length()
	for k, v := range path.Rewards() {
		pathCopy.rewards[k] = v
	}
	return pathCopy
}

// PossibleRoute checks if an edge makes an eligible route to take for the current path.
func (path *Path) PossibleRoute(edge *Edge) (bool, int) {

	if path.rewardsAreNotNegativeOrUnique(edge.To()) {
		return false, -1
	}

	nodeVisited := path.hasBeenVisited(edge)

	if !edge.To().Revisitable() && nodeVisited {
		return false, -1
	}

	rewardsChangedSinceLastVisit := path.rewardsChangedSinceLastVisit(edge)

	if edge.To().Revisitable() && nodeVisited && !rewardsChangedSinceLastVisit {
		return false, -1
	}

	return path.requirementsMet(edge)
}

func (path *Path) hasBeenVisited(edge *Edge) bool {
	node := edge.To()
	for i := 0; i < len(path.edges); i++ {
		if path.edges[i].From() == node {
			return true
		}
	}
	return edge.From() == node
}

func (path *Path) rewardsChangedSinceLastVisit(edge *Edge) bool {
	rewards := make(map[*Reward]int)
	node := edge.To()

	for k, v := range edge.From().Rewards() {
		rewards[k] = rewards[k] + v
	}

	if edge.From() == edge.To() {
		return checkRewardsNotEmpty(rewards)
	}

	for i := len(path.edges) - 1; i >= 0; i-- {
		for k, v := range path.edges[i].From().Rewards() {
			rewards[k] = rewards[k] + v
		}
		if path.edges[i].From() == node {
			break
		}
	}

	return checkRewardsNotEmpty(rewards)
}

func checkRewardsNotEmpty(rewards map[*Reward]int) bool {
	for _, v := range rewards {
		if v != 0 {
			return true
		}
	}

	return false
}

func (path *Path) rewardsAreNotNegativeOrUnique(node *Node) bool {
	for reward, quantity := range node.Rewards() {
		rewardCount := path.Rewards()[reward]
		if rewardCount+quantity < 0 || (reward.Unique() && rewardCount > 0) {
			return true
		}
	}
	return false
}

func (path *Path) requirementsMet(edge *Edge) (bool, int) {
	for i, weight := range edge.Weights() {
		if path.weightRequirementsMet(weight) {
			return true, i
		}
	}
	return false, -1
}

func (path *Path) weightRequirementsMet(weight *Weight) bool {
	for reward, quantity := range weight.Requirements() {
		if path.countRewardsIncludingIsA(reward) < quantity {
			return false
		}
	}
	return true
}

func (path *Path) countRewardsIncludingIsA(reward *Reward) int {
	count := path.Rewards()[reward]
	for _, r := range reward.CanBe() {
		count += path.countRewardsIncludingIsA(r)
	}
	return count
}

// PrioQueue is a list of Path pointers and implements heap interface.
type PrioQueue []*Path

// Len returns the length of the PriorityQueue.
func (pq PrioQueue) Len() int {
	return len(pq)
}

// Less checks if one path is longer than another path in the PriorityQueue.
func (pq PrioQueue) Less(i, j int) bool {
	return pq[j].IsLongerThan(pq[i])
}

// Swap switches places of two paths in the PriorityQueue.
func (pq PrioQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push adds a Path to the correct place in the PriorityQueue.
func (pq *PrioQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Path))
}

// Pop returns the node with the best potential path for the shortest path.
func (pq *PrioQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}
