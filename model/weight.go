package model

type Weight struct {
	time         int
	requirements map[*Reward]int
}

func (weight *Weight) Time() int {
	return weight.time
}

func (weight *Weight) Requirements() map[*Reward]int {
	return weight.requirements
}

func (weight *Weight) AddRequirement(reward *Reward, quantity int) {
	weight.requirements[reward] = quantity
}

func CreateWeight(time int) *Weight {
	weight := new(Weight)
	weight.time = time
	weight.requirements = make(map[*Reward]int)
	return weight
}

type ByTime []*Weight

func (a ByTime) Len() int {
	return len(a)
}

func (a ByTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByTime) Less(i, j int) bool {
	return a[i].Time() < a[j].Time()
}
