package controller

import (
	"context"
	"log"

	"mafia/client/cli"
	connection "mafia/pkg/proto/connection"
	pgame "mafia/pkg/proto/game"
)

func (c *Controller) ProcessNight(client connection.MafiaServerClient) {
	if c.Role == pgame.Role_HUMAN {
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
	_, err = client.Commit(context.Background(), &pgame.CommitRequest{UserId: c.ID, Target: selected})
	if err != nil {
		log.Printf("failed to commit: %v", err)
		return
	}
	// TODO обработать ответ от сервера для коммисара. Сохранять его куда-то?
	// Надо добавить для комиссара опцию, чтобы он мог объявлять кто мафия, а кто нет.
}
