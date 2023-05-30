package roles

type Mafia struct {
	BaseRole
}

func (m *Mafia) NeedProcess() bool {
	return true
}
