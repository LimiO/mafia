package controller

import (
	"fmt"
	"mafia/internal"
	"math/rand"
	"time"

	"mafia/client/cli"
	pgame "mafia/pkg/proto/game"
	"mafia/roles"
)

const (
	chatOption        = "chat"
	votebanOption     = "voteban"
	endDayOption      = "end day"
	publishInfoOption = "publish info"
)

type Participant struct {
	ID    string
	Alive bool
}

type Controller struct {
	State pgame.State
	Role  roles.Role

	IsAuto    bool
	DayNumber int
	GameID    uint32

	Participants map[string]*Participant
	ID           string
}

func (c *Controller) SelectAction(msg string, options []string) string {
	if c.IsAuto {
		rand.Seed(time.Now().UnixNano())
		selected := options[rand.Intn(100001)%len(options)]
		fmt.Printf("Selected random option to msg \"%s...\": %q\n", msg[:10], selected)
		return selected
	}
	return cli.AskSelect(msg, options)
}

func (c *Controller) AskInput(msg string) string {
	if c.IsAuto {
		rand.Seed(time.Now().UnixNano())
		result := internal.RandStringRunes(10)
		fmt.Printf("message to send: %q\n", result)
		return result
	}
	return cli.AskInput(msg)
}

func (c *Controller) SetGameID(gameID uint32) {
	c.GameID = gameID
}

func (c *Controller) SetRole(role pgame.Role) {
	switch role {
	case pgame.Role_HUMAN:
		c.Role = &roles.Human{}
	case pgame.Role_MAFIA:
		c.Role = &roles.Mafia{}
	case pgame.Role_POLICE:
		c.Role = &roles.Police{}
	}
}

func (c *Controller) MakeParticipantIds() []string {
	var ids []string
	for userID := range c.Participants {
		if userID == c.ID {
			continue
		}
		ids = append(ids, userID)
	}
	return ids
}

func (c *Controller) MakeAliveParticipantIds() []string {
	var ids []string
	for userID, p := range c.Participants {
		if userID == c.ID {
			continue
		}
		if !p.Alive {
			continue
		}
		ids = append(ids, userID)
	}
	return ids
}
