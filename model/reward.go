package model

// Reward stores reward id, if it is unique and if can be counted as another reward via inheritance.
type Reward struct {
	id     string
	unique bool
	isA    *Reward
	canBe  []*Reward
}

// ID returns reward id.
func (r *Reward) ID() string {
	return r.id
}

// Unique returns if the reward is unique.
func (r *Reward) Unique() bool {
	return r.unique
}

// IsA returns reference to inherited reward.
func (r *Reward) IsA() *Reward {
	return r.isA
}

// CanBe returns all references which inherits this reward.
func (r *Reward) CanBe() []*Reward {
	return r.canBe
}

// CreateReward constructs a new Reward.
func CreateReward(id string, unique bool, isa *Reward) *Reward {
	reward := &Reward{id, unique, isa, nil}
	if isa != nil {
		isa.canBe = append(isa.canBe, reward)
	}
	return reward
}
