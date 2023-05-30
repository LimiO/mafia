package server

import (
	"fmt"
	"mafia/pkg/proto/game"
	"strings"
)

func (g *Game) ProcessWarmup() error {
	if len(g.users) != MinPlayers {
		return nil
	}
	err := g.StartGame()
	if err != nil {
		return fmt.Errorf("failed to start controller: %v", err)
	}
	return nil
}

func (g *Game) ProcessDay() {
	if len(g.status.Ended) != g.GetAliveCount() {
		return
	}
	g.Ban()
	g.status.VoteBanned = map[string]string{}

	g.status.Ended = map[string]bool{}
	if status := g.EndGameStatus(); status != STATUS_NOT_YET {
		g.status.EndGame(g.gameID)
	} else {
		g.SendState(game.State_NIGHT)
		g.status.SetNight(g.gameID)
	}
}

func (g *Game) ProcessNight() {
	if len(g.status.Commited) != MinCommitPlayers {
		return
	}
	var results []string

	for userID, target := range g.status.Commited {
		if _, ok := g.users[userID]; !ok || g.status.Roles[userID] != game.Role_MAFIA {
			continue
		}
		user, ok := g.users[target]
		if !ok {
			continue
		}
		user.alive = false
		results = append(results, fmt.Sprintf("player %q dead", target))
		g.SendKillNotification(target)
	}

	g.SendToChat("game", strings.Join(results, ", "))
	g.status.Commited = map[string]string{}
	if status := g.EndGameStatus(); status != STATUS_NOT_YET {
		g.status.EndGame(g.gameID)
	} else {
		g.status.SetDay(g.gameID)
		g.SendToChat("game", "GAME STATUS: DAY")
		g.SendState(game.State_DAY)
	}
}
