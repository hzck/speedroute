package model

type Reward struct {
	id     string
	unique bool
}

func (r *Reward) Id() string {
	return r.id
}

func (r *Reward) Unique() bool {
	return r.unique
}

func (reward *Reward) SetUnique(unique bool) {
	reward.unique = unique
}

func CreateReward(id string) *Reward {
	r := new(Reward)
	r.id = id
	r.unique = false
	return r
}
