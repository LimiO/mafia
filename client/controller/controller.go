package controller

import (
	pgame "mafia/pkg/proto/game"
)

const (
	chatOption    = "chat"
	votebanOption = "voteban"
	endDayOption  = "end day"
)

type Participant struct {
	ID    string
	Alive bool
}

type Controller struct {
	State pgame.State
	Role  pgame.Role

	Participants map[string]*Participant
	ID           string
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
