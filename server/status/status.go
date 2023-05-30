package status

import (
	"log"
	pgame "mafia/pkg/proto/game"
)

type Status struct {
	State      pgame.State
	Roles      map[string]pgame.Role
	VoteBanned map[string]string
	Commited   map[string]string
	Ended      map[string]bool
}

func (s *Status) SetDay(gameID uint32) {
	log.Printf("game №%d - status: DAY", gameID)
	s.State = pgame.State_DAY
}

func (s *Status) SetNight(gameID uint32) {
	log.Printf("game №%d - status: NIGHT", gameID)
	s.State = pgame.State_NIGHT
}

func (s *Status) EndGame(gameID uint32) {
	log.Printf("game №%d - status: END", gameID)
	s.State = pgame.State_END
}

func (s *Status) GameStarted() bool {
	return s.State != pgame.State_UNKNOWN
}

func (s *Status) GameEnded() bool {
	return s.State == pgame.State_END
}
