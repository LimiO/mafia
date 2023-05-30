package roles

type Police struct {
	BaseRole
}

func (h *Police) NeedProcess() bool {
	return true
}
