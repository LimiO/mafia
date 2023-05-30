package controller

import (
	"context"
	"log"

	"mafia/client/cli"
	connection "mafia/pkg/proto/connection"
	pgame "mafia/pkg/proto/game"
)

func (c *Controller) ProcessNight(client connection.MafiaServerClient) {
	if !c.Role.NeedProcess() {
		c.State = pgame.State_SPIRIT
		return
	}

	ids := c.MakeAliveParticipantIds()
	selected, err := cli.AskSelect("Select target to do", ids)
	if err != nil {
		log.Printf("failed to ask select: %v", err)
		return
	}
	c.State = pgame.State_SPIRIT

	rsp, err := client.Commit(context.Background(), &pgame.CommitRequest{
		UserId: c.ID,
		Target: selected,
		Game:   c.GameID,
	})
	if err != nil {
		log.Printf("failed to commit: %v", err)
		return
	}
	if rsp.Result == pgame.CommitResponse_FAIL {
		return
	}
	c.Role.SetInfo(selected)
}
