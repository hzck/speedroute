package model

type Path struct {
	edges   []*Edge
	length  int
	rewards map[*Reward]int
}

func (path *Path) Edges() []*Edge {
	return path.edges
}

func (path *Path) AddEdge(edge *Edge, i int) {
	path.edges = append(path.edges, edge)
	path.length += edge.Weights()[i].Time()
	path.AddRewards(edge.To().Rewards())
}

func (path *Path) Length() int {
	return path.length
}

func (path *Path) Rewards() map[*Reward]int {
	return path.rewards
}

func (path *Path) AddRewards(rewards map[*Reward]int) {
	for key, value := range rewards {
		path.rewards[key] += value
	}
}

func CreatePath() *Path {
	path := new(Path)
	path.rewards = make(map[*Reward]int)
	return path
}

func (path *Path) IsLongerThan(other *Path) bool {
	plength := path.Length() + path.Edges()[len(path.Edges())-1].To().MinPathLeft()
	oLength := other.Length() + other.Edges()[len(other.Edges())-1].To().MinPathLeft()
	pEdges := len(path.Edges())
	oEdges := len(other.Edges())
	return plength > oLength || (plength == oLength && pEdges > oEdges)
}

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

func (path *Path) PossibleRoute(edge *Edge) (bool, int) {
	if path.visitable(edge.To()) {
		return path.requirementsMet(edge)
	}
	return false, -1
}

func (path *Path) visitable(node *Node) bool {
	for i := len(path.edges) - 1; i >= 0; i-- {
		if path.edges[i].To() == node {
			return false
		}
		if !node.Revisitable() && path.edges[i].From() == node {
			return false
		}
		if node.Revisitable() && len(path.edges[i].To().Rewards()) > 0 {
			return true
		}
	}
	for reward, quantity := range node.Rewards() {
		rewardCount := path.Rewards()[reward]
		if rewardCount+quantity < 0 || (reward.Unique() && rewardCount > 0) {
			return false
		}
	}
	return true
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
		if path.Rewards()[reward] < quantity {
			return false
		}
	}
	return true
}

type PrioQueue []*Path

func (pq PrioQueue) Len() int {
	return len(pq)
}

func (pq PrioQueue) Less(i, j int) bool {
	return pq[j].IsLongerThan(pq[i])
}

func (pq PrioQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PrioQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Path))
}

func (pq *PrioQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}
