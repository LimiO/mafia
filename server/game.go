package server

import (
	"fmt"
	"sync"

	"mafia/pkg/proto/game"
	"mafia/server/status"
)

type Filter func(userID string) bool

const (
	STATUS_NOT_YET = 0
	STATUS_MAFIA   = 1
	STATUS_HUMAN   = 2
)

type Game struct {
	users  map[string]*User
	mu     sync.Mutex
	status status.Status
	gameID uint32
}

func NewGame() *Game {
	g := &Game{
		users: make(map[string]*User),
		status: status.Status{
			VoteBanned: map[string]string{},
			Ended:      map[string]bool{},
			Commited:   map[string]string{},
			Roles:      map[string]game.Role{},
		},
		gameID: GlobalGameID,
	}
	GlobalGameID++
	return g
}

func (g *Game) AliveFilter(userID string) bool {
	user, ok := g.users[userID]
	if !ok {
		return true
	}
	return !user.alive
}

func (g *Game) EndGameStatus() int {
	var countHuman int
	var countMafia int
	for _, user := range g.users {
		if !user.alive {
			continue
		}
		if user.role == game.Role_MAFIA {
			countMafia++
		} else {
			countHuman++
		}
	}
	if countMafia >= countHuman {
		return STATUS_MAFIA
	}
	if countMafia == 0 {
		return STATUS_HUMAN
	}
	return STATUS_NOT_YET
}

func (g *Game) GetAliveCount() int {
	var count int
	for _, user := range g.users {
		if user.alive {
			count++
		}
	}
	return count
}

func (g *Game) Ban() {
	if len(g.status.VoteBanned) == 0 {
		return
	}
	var maxBannedID string
	var maxVoted uint32

	bannedCounter := map[string]uint32{}

	for _, target := range g.status.VoteBanned {
		bannedCounter[target]++
		if bannedCounter[target] > maxVoted {
			maxBannedID = target
			maxVoted = bannedCounter[target]
		}
	}
	g.users[maxBannedID].alive = false
	g.SendToChat("game", fmt.Sprintf("user %q because of poll", maxBannedID))
	g.SendKillNotification(maxBannedID)
}
