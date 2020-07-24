package model

// Weight holds the time and possible requirements to traverse an edge.
type Weight struct {
	time         int
	requirements map[*Reward]int
}

// Time returns a time.
func (weight *Weight) Time() int {
	return weight.time
}

// Requirements returns a map of requirements.
func (weight *Weight) Requirements() map[*Reward]int {
	return weight.requirements
}

// AddRequirement adds a requirement with a quantity.
func (weight *Weight) AddRequirement(reward *Reward, quantity int) {
	weight.requirements[reward] = quantity
}

// CreateWeight constructs a weight and returns a pointer to it.
func CreateWeight(time int) *Weight {
	weight := new(Weight)
	weight.time = time
	weight.requirements = make(map[*Reward]int)
	return weight
}

// ByTime is a list of Weights which implements the sort interface.
type ByTime []*Weight

// Len returns the ByTime length.
func (a ByTime) Len() int {
	return len(a)
}

// Swap changes position of two weights in the ByTime list.
func (a ByTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less compares two weight times to each other.
func (a ByTime) Less(i, j int) bool {
	return a[i].Time() < a[j].Time()
}
