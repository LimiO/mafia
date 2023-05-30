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

func (g *Game) StartGame() error {
	var users []string
	for user := range g.users {
		users = append(users, user)
	}
	roles := internal.ShuffleRoles(users)

	for userID, user := range g.users {
		user.role = roles[userID]
		err := user.conn.Send(&connection.ServerResponse{
			Response: &connection.ServerResponse_Start{
				Start: &connection.StartGameResponse{
					Role:  roles[userID],
					Users: users,
					Game:  g.gameID,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to send role to player %q: %v", userID, err)
		}
	}
	for userID, role := range roles {
		g.status.Roles[userID] = role
	}
	g.SendToChat("game", "Game started!")
	g.status.SetDay(g.gameID)
	return nil
}

func (s *Server) PollGame() error {
	for {
		time.Sleep(time.Second)
		games := s.Games
		for _, g := range games {
			var err error
			switch g.status.State {
			case game.State_UNKNOWN:
				err = g.ProcessWarmup()
			case game.State_DAY:
				g.ProcessDay()
			case game.State_NIGHT:
				g.ProcessNight()
			case game.State_END:
				s.StopGame(g)
			}
			if err != nil {
				log.Println(err)
			}
		}

	}
}

func (s *Server) StopGame(g *Game) {
	delete(s.Games, g.gameID)
	var msg string
	switch g.EndGameStatus() {
	case STATUS_MAFIA:
		msg = "Mafia won!"
	case STATUS_HUMAN:
		msg = "People won!"
	}
	g.SendToChat("game", msg)
	g.SendState(game.State_END)
}
