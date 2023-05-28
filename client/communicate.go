package client

import (
	"context"
	"fmt"
	"io"
	"log"
	pgame "mafia/pkg/proto/game"
	"strings"
	"time"

	connection "mafia/pkg/proto/connection"
)

// stream

func (c *Client) JoinGame() error {
	req := &connection.UserJoinRequest{
		UserId: c.Ctl.ID,
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Minute))
	defer cancel()
	rsp, err := c.GrpcClient.Connect(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	for {
		if c.Ctl.State == pgame.State_END {
			return nil
		}
		resp, err := rsp.Recv()
		if err == io.EOF {
			log.Printf("server dead: %v", err)
			return nil
		}
		if err != nil {
			log.Printf("failed to receive rsp from server: %v", err)
			continue
		}
		err = c.ProcessStreamResponse(resp)
		if err != nil {
			log.Printf("failed to process response: %v", err)
			continue
		}
		time.Sleep(time.Second)
	}
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

// other

func (c *Client) ListParticipants() error {
	req := &connection.ListParticipantsRequest{}
	rsp, err := c.GrpcClient.ListParticipants(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to send 'listParticipands' request: %v", err)
	}
	return c.ProcessListParticipants(rsp)
}

func (c *Client) ProcessListParticipants(rsp *connection.ListParticipantsResponse) error {
	log.Printf("list participants: %s", strings.Join(rsp.Users, ", "))
	return nil
}
