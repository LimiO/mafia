package gamecontroller

import (
	pgame "mafia/pkg/proto/game"
	"mafia/roles"
)

const (
	chatOption         = "chat"
	votebanOption      = "voteban"
	endDayOption       = "end day"
	publishInfoOption  = "publish info"
	selectTargetOption = "select target"
)

type Participant struct {
	ID    string
	Alive bool
}

type Controller struct {
	State pgame.State
	Role  roles.Role

	ChatChan chan string

	IsAuto    bool
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
