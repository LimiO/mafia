package roles

type Human struct {
	BaseRole
}

func (h *Human) NeedProcess() bool {
	return false
}
