package status

import (
	pgame "mafia/pkg/proto/game"
)

type Status struct {
	State      pgame.State
	Roles      map[string]pgame.Role
	VoteBanned map[string]string
	Commited   map[string]string
	Ended      map[string]bool
}

func (s *Status) SetDay() {
	s.State = pgame.State_DAY
}

func (s *Status) SetNight() {
	s.State = pgame.State_NIGHT
}

func (s *Status) EndGame() {
	s.State = pgame.State_END
}

func (s *Status) GameStarted() bool {
	return s.State != pgame.State_UNKNOWN
}

func (s *Status) GameEnded() bool {
	return s.State == pgame.State_END
}
