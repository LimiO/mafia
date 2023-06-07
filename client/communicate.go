package client

import (
	"context"
	"fmt"
	"io"
	"log"
	pgame "mafia/pkg/proto/game"
	"time"

	connection "mafia/pkg/proto/connection"
)

func (c *Client) JoinGame() error {
	req := &connection.UserJoinRequest{
		UserId:   c.GameCtl.ID,
		Password: c.GameCtl.Password,
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(20*time.Minute))
	defer cancel()
	rsp, err := c.GrpcClient.Connect(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	for {
		if c.GameCtl.State == pgame.State_END {
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
