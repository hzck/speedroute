package model

// Reward stores reward id and if it is unique.
type Reward struct {
	id     string
	unique bool
}

// ID returns reward id.
func (r *Reward) ID() string {
	return r.id
}

// Unique returns if the reward is unique.
func (r *Reward) Unique() bool {
	return r.unique
}

// CreateReward constructs a new Reward.
func CreateReward(id string, unique bool) *Reward {
	return &Reward{id, unique}
}
