package server

import (
	"fmt"
	"log"
	"mafia/pkg/proto/game"
	"strings"
)

func (s *Server) ProcessWarmup() error {
	if len(s.users) != MinPlayers {
		return nil
	}
	err := s.StartGame()
	if err != nil {
		return fmt.Errorf("failed to start controller: %v", err)
	}
	return nil
}

func (s *Server) Ban() {
	var maxBannedID string
	var maxVoted uint32

	bannedCounter := map[string]uint32{}

	for _, target := range s.status.VoteBanned {
		bannedCounter[target]++
		if bannedCounter[target] > maxVoted {
			maxBannedID = target
			maxVoted = bannedCounter[target]
		}
	}
	s.users[maxBannedID].alive = false
	s.SendStateTo(maxBannedID, game.State_SPIRIT)

}

func (s *Server) ProcessDay() {
	log.Println("status: DAY")
	if len(s.status.Ended) != s.GetAliveCount() {
		return
	}
	//s.Ban()
	s.status.VoteBanned = map[string]string{}

	s.status.Ended = map[string]bool{}
	s.SendToChat("game", "GAME STATUS: NIGHT")
	s.SendState(game.State_NIGHT)
	s.status.SetNight()
}

func (s *Server) ProcessNight() {
	log.Println("status: NIGHT")
	if len(s.status.Commited) != MinCommitPlayers {
		return
	}
	var results []string

	for userID, target := range s.status.Commited {
		if user, ok := s.users[userID]; !ok || user.role != game.Role_MAFIA {
			continue
		}
		user, ok := s.users[target]
		if !ok {
			continue
		}
		user.alive = false
		results = append(results, fmt.Sprintf("player %q dead", target))
		s.SendStateTo(target, game.State_SPIRIT)
		s.SendKillNotification(target)
	}

	s.SendToChat("game", strings.Join(results, ", "))
	s.status.Commited = map[string]string{}
	if s.NeedEndGame() {
		log.Println("status: END")
		s.status.EndGame()
	} else {
		s.status.SetDay()
		s.SendToChat("game", "GAME STATUS: DAY")
		s.SendState(game.State_DAY)
	}
}

func (s *Server) NeedEndGame() bool {
	var countHuman int
	var countMafia int
	for _, user := range s.users {
		if !user.alive {
			continue
		}
		if user.role == game.Role_HUMAN {
			countHuman++
		}
		if user.role == game.Role_MAFIA {
			countMafia++
		}
	}
	return countMafia >= countHuman || countMafia == 0
}

func (s *Server) GetAliveCount() int {
	var count int
	for _, user := range s.users {
		if user.alive {
			count++
		}
	}
	return count
}
