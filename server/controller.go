package server

import (
	"fmt"
	"log"
	"mafia/internal/helpers"
	"time"

	"golang.org/x/sync/errgroup"

	connection "mafia/pkg/proto/connection"
	"mafia/pkg/proto/game"
)

func (s *Server) Start() error {
	var eg errgroup.Group

	eg.Go(s.StartListen)
	eg.Go(s.PollGame)

	err := eg.Wait()
	return err
}

func (g *Game) StartGame() error {
	g.startTime = time.Now()

	var users []string
	for user := range g.users {
		users = append(users, user)
	}
	roles := helpers.ShuffleRoles(users)

	g.SendToChat("game", "Game started!")
	time.Sleep(3 * time.Second)

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
	g.status.Started = true
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

func (s *Server) UpdateStats(g *Game) error {
	now := time.Since(g.startTime)
	for userID, user := range g.users {
		stats, err := s.DBManager.GetStats(userID)
		if err != nil {
			return fmt.Errorf("failed to get stats: %v", err)
		}
		stats.CountGames++

		isMafia := user.role == game.Role_MAFIA
		if isMafia && g.EndGameStatus() == STATUS_MAFIA {
			stats.CountWins++
		} else if !isMafia && g.EndGameStatus() == STATUS_HUMAN {
			stats.CountWins++
		}
		stats.TotalTime += int(now.Seconds())

		if err = s.DBManager.UpdateStats(stats); err != nil {
			return fmt.Errorf("failed to update stats: %v", err)
		}
	}
	return nil
}

func (s *Server) StopGame(g *Game) error {
	delete(s.Games, g.gameID)
	var msg string
	switch g.EndGameStatus() {
	case STATUS_MAFIA:
		msg = "Mafia won!"
	case STATUS_HUMAN:
		msg = "People won!"
	}

	// TODO добавить статус, что юзер уже такой есть
	err := s.UpdateStats(g)
	if err != nil {
		return fmt.Errorf("failed to update stats: %v", err)
	}
	g.SendToChat("game", msg)
	g.SendState(game.State_END)
	return nil
}
