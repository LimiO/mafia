package roles

type Role interface {
	Commit(target *BaseRole)
	VoteBan(target *BaseRole)
	IsDead() bool
	Die()
}

type BaseRole struct {
	dead      bool
	wishToBan uint32
}

func (r *BaseRole) Commit(_ *BaseRole) {
}
func (r *BaseRole) IsDead() bool {
	return r.dead
}
func (r *BaseRole) Die() {
	r.dead = true
}
func (r *BaseRole) VoteBan(other *BaseRole) {
	other.wishToBan++
}
