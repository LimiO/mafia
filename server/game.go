package server

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"mafia/internal/queue"
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

	QueueCtl *queue.Controller

	gameID uint32
}

func NewGame() (*Game, error) {
	g := &Game{
		users: make(map[string]*User),
		status: status.Status{
			VoteBanned: map[string]string{},
			Ended:      map[string]bool{},
			Commited:   map[string]string{},
			Roles:      map[string]game.Role{},
			Started:    false,
		},
		gameID: GlobalGameID,
	}
	cfg := &queue.Config{
		Addr:                 queue.Addr,
		RoutingKeys:          []string{queue.AllKey, queue.MafiaKey},
		ProducerExchangeName: []string{},
		ConsumerExchangeName: "server",
	}
	queueCtl, err := queue.NewController(cfg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make queue controller: %v", err)
	}
	g.QueueCtl = queueCtl

	go func() error {
		return queueCtl.StartConsume(func(delivery amqp.Delivery) error {
			for userID := range g.users {
				exchangeName := fmt.Sprintf("client.%d.%s", g.gameID, userID)
				err = queueCtl.Push(exchangeName, delivery.RoutingKey, string(delivery.Body))
				if err != nil {
					return fmt.Errorf("failed to push to key %s: %v", delivery.RoutingKey, err)
				}
			}
			return nil
		})
	}()
	GlobalGameID++
	return g, nil
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

func (g *Game) DeleteUser(userID string) {
	delete(g.users, userID)
	delete(g.status.Ended, userID)
	delete(g.status.Roles, userID)
	delete(g.status.VoteBanned, userID)
	delete(g.status.Commited, userID)
	if g.status.State == game.State_END {
		return
	}
	g.SendToChat(userID, "disconnected from the game")
	g.SendKillNotification(userID)
	log.Printf("user %q disconnected", userID)
	if len(g.users) == 0 && g.status.Started {
		g.status.EndGame(g.gameID)
	}
}
