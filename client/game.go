package client

import (
	"fmt"
	"mafia/client/controller"
	connection "mafia/pkg/proto/connection"
	"mafia/pkg/proto/game"
)

func (c *Client) ProcessStartResponse(rsp *connection.StartGameResponse) error {
	fmt.Println("Your role", rsp.GetRole())
	c.Ctl.SetRole(rsp.GetRole())
	c.Ctl.SetGameID(rsp.Game)
	users := rsp.GetUsers()
	for _, user := range users {
		c.Ctl.Participants[user] = &controller.Participant{
			ID:    user,
			Alive: true,
		}
	}
	c.Ctl.State = game.State_DAY
	go c.Ctl.ProcessDay(c.GrpcClient)
	return nil
}

func (c *Client) ProcessStateResponse(rsp *game.StateResponse) error {
	c.Ctl.State = rsp.GetState()
	switch rsp.GetState() {
	case game.State_DAY:
		fmt.Println("GAME STATE: DAY")
		go c.Ctl.ProcessDay(c.GrpcClient)
	case game.State_NIGHT:
		fmt.Println("GAME STATE: NIGHT")
		go c.Ctl.ProcessNight(c.GrpcClient)
	case game.State_SPIRIT:
		fmt.Println("GAME STATE: YOU ARE SPIRIT NOW")
		go c.Ctl.ProcessSpirit()
	case game.State_END:
		fmt.Println("GAME STATE: END")
		// TODO отправить сообщение что ты выиграл или проиграл
		return nil
	}
	return nil
}

func (c *Client) ProcessKillResponse(rsp *game.KillResponse) error {
	userID := rsp.GetUserId()
	participant, ok := c.Ctl.Participants[userID]
	if !ok {
		return fmt.Errorf("cannot find user %q to kill", userID)
	}
	participant.Alive = false
	if userID == c.Ctl.ID {
		fmt.Println("You are dead!")
		c.Ctl.State = game.State_SPIRIT
	}
	return nil
}
