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

func CreateReward(id string, unique bool) *Reward {
	r := new(Reward)
	r.id = id
	r.unique = unique
	return r
}
