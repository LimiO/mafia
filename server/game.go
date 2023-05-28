package server

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"

	connection "mafia/pkg/proto/connection"
	"mafia/pkg/proto/game"
	"mafia/server/internal"
)

func (s *Server) Start() error {
	var eg errgroup.Group

	eg.Go(s.StartListen)
	eg.Go(s.PollGame)

	err := eg.Wait()
	return err
}

func (s *Server) StartGame() error {
	s.status.SetDay()

	var users []string
	for user := range s.users {
		users = append(users, user)
	}
	roles := internal.ShuffleRoles(users)
	for userID, user := range s.users {
		user.role = roles[userID]
		err := user.conn.Send(&connection.ServerResponse{
			Response: &connection.ServerResponse_Start{
				Start: &connection.StartGameResponse{
					Role:  roles[userID],
					Users: users,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to send role to player %q: %v", user, err)
		}
	}
	for userID, role := range roles {
		s.status.Roles[userID] = role
	}
	s.SendToChat("controller", "Game started!")
	return nil
}

func (s *Server) PollGame() error {
	for {
		time.Sleep(time.Second)

		var err error
		switch s.status.State {
		case game.State_UNKNOWN:
			err = s.ProcessWarmup()
		case game.State_DAY:
			s.ProcessDay()
		case game.State_NIGHT:
			s.ProcessNight()
		case game.State_END:
			s.StopGame()
			return nil
		}
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) StopGame() {
	s.SendToChat("game", "GAME STATUS: ENDED!")
	s.SendState(game.State_END)
	s.Stop()
}
