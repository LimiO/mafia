package roles

type Role interface {
	NeedProcess() bool
	GetInfo() string
	SetInfo(string)
	IsDead() bool
	Die()
}

type BaseRole struct {
	dead bool
	info string
}

func (r *BaseRole) GetInfo() string {
	return r.info
}

func (r *BaseRole) SetInfo(info string) {
	r.info = info
}

func (r *BaseRole) IsDead() bool {
	return r.dead
}

func (r *BaseRole) Die() {
	r.dead = true
}
