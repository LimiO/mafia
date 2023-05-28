package roles

type Human struct {
	BaseRole
}

func (h *Human) Commit(target *BaseRole) {
	h.VoteBan(target)
}
