package controller

import (
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

	DayNumber int
	GameID    uint32

	Participants map[string]*Participant
	ID           string
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
