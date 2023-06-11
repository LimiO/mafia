package gamecontroller

import (
	"context"
	"fmt"
	"mafia/roles"

	connection "mafia/pkg/proto/connection"
	pgame "mafia/pkg/proto/game"
)

func (c *Controller) GetNightOptions() []string {
	options := []string{selectTargetOption}
	if _, ok := c.Role.(*roles.Mafia); ok {
		options = append(options, chatOption)
	}
	return options
}

func (c *Controller) ProcessNight(client connection.MafiaServerClient) {
	if !c.Role.NeedProcess() {
		c.State = pgame.State_SPIRIT
		return
	}

	for {
		if c.State != pgame.State_NIGHT {
			return
		}
		selected := c.SelectAction("Select option to do", c.GetNightOptions())
		switch selected {
		case chatOption:
			c.GetAndSendMessage()
		case selectTargetOption:
			ids := c.makeAliveParticipantIds()
			target := c.SelectAction("Select target to do", ids)
			rsp, err := client.Commit(context.Background(), &pgame.CommitRequest{
				UserId: c.ID,
				Target: target,
				Game:   c.GameID,
			})
			if rsp.Result == pgame.CommitResponse_FAIL {
				return
			}
			c.Role.SetInfo(target)
			if err != nil {
				fmt.Printf("failed to commit: %v", err)
			}
			return
		}
	}
}
