package roles

type Mafia struct {
	BaseRole
}

func (m *Mafia) Commit(target *BaseRole) error {
	target.Die()
	return nil
}
