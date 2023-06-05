package client

import (
	"fmt"
	"log"
	"mafia/client/gamecontroller"
	"mafia/internal/queue"
	connection "mafia/pkg/proto/connection"
	"mafia/pkg/proto/game"
)

func (c *Client) ProcessStartResponse(rsp *connection.StartGameResponse) error {
	fmt.Println("Your role", rsp.GetRole())
	c.GameCtl.SetRole(rsp.GetRole())
	c.GameCtl.SetGameID(rsp.Game)
	users := rsp.GetUsers()
	for _, user := range users {
		c.GameCtl.Participants[user] = &gamecontroller.Participant{
			ID:    user,
			Alive: true,
		}
	}
	c.GameCtl.State = game.State_DAY

	cfg := &queue.Config{
		Addr:                 "amqp://guest:guest@localhost:5672/",
		RoutingKeys:          []string{queue.AllKey, queue.MafiaKey},
		ProducerExchangeName: []string{},
		ConsumerExchangeName: fmt.Sprintf("client.%d.%s", rsp.Game, c.GameCtl.ID),
	}
	queueCtl, err := queue.NewController(cfg, c.GameCtl.ChatChan)
	if err != nil {
		return fmt.Errorf("failed to make queue controller: %v", err)
	}
	err = queueCtl.AddProducer("server")
	if err != nil {
		return fmt.Errorf("failed to add producer: %v", err)
	}
	c.QueueCtl = queueCtl

	go c.GameCtl.ProcessDay(c.GrpcClient)
	return nil
}

func (c *Client) ProcessStateResponse(rsp *game.StateResponse) error {
	c.GameCtl.State = rsp.GetState()
	switch rsp.GetState() {
	case game.State_DAY:
		fmt.Println("GAME STATE: DAY")
		go c.GameCtl.ProcessDay(c.GrpcClient)
	case game.State_NIGHT:
		fmt.Println("GAME STATE: NIGHT")
		c.GameCtl.ProcessNight(c.GrpcClient)
	case game.State_SPIRIT:
		fmt.Println("GAME STATE: YOU ARE SPIRIT NOW")
		c.GameCtl.ProcessSpirit()
	case game.State_END:
		fmt.Println("GAME STATE: END")
		return nil
	}
	return nil
}

func (c *Client) ProcessKillResponse(rsp *game.KillResponse) error {
	userID := rsp.GetUserId()
	participant, ok := c.GameCtl.Participants[userID]
	if !ok {
		return nil
	}
	participant.Alive = false
	if userID == c.GameCtl.ID {
		fmt.Println("You are dead!")
		c.GameCtl.State = game.State_SPIRIT
	}
	return nil
}

func (c *Client) ProcessJoinResponse(rsp *connection.UserJoinResponse) error {
	switch rsp.Type {
	case connection.UserJoinResponse_OK:
		log.Println("User joined")
	case connection.UserJoinResponse_EXISTS:
		log.Fatalf("User already exists, join with another name")
	case connection.UserJoinResponse_STARTED:
		log.Fatalf("Game already started")
	}
	return nil
}

func (c *Client) ProcessChatResponse(rsp *connection.ChatResponse) error {
	log.Printf("user %q: %s", rsp.GetUserId(), rsp.GetText())
	return nil
}
