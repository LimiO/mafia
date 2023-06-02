package client

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"mafia/client/controller"
	connection "mafia/pkg/proto/connection"
)

type Client struct {
	ServerPort uint32
	Ctl        controller.Controller

	GrpcClient connection.MafiaServerClient
}

func MakeClient(id string) (*Client, error) {
	return &Client{
		ServerPort: 9001,
		Ctl: controller.Controller{
			ID:           id,
			Participants: map[string]*controller.Participant{},
		},
	}, nil
}

func (c *Client) ProcessStreamResponse(rsp *connection.ServerResponse) error {
	switch rsp.Response.(type) {
	case *connection.ServerResponse_Join:
		return c.ProcessJoinResponse(rsp.GetJoin())
	case *connection.ServerResponse_Chat:
		return c.ProcessChatResponse(rsp.GetChat())
	case *connection.ServerResponse_State:
		return c.ProcessStateResponse(rsp.GetState())
	case *connection.ServerResponse_Start:
		return c.ProcessStartResponse(rsp.GetStart())
	case *connection.ServerResponse_Kill:
		return c.ProcessKillResponse(rsp.GetKill())
	default:
		return nil
	}
}

func (c *Client) StartSession() error {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", c.ServerPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("failed to connect to grpc: %v", err)
	}
	defer conn.Close()

	c.GrpcClient = connection.NewMafiaServerClient(conn)
	var eg errgroup.Group

	eg.Go(func() error {
		err = c.JoinGame()
		if err != nil {
			return fmt.Errorf("failed to join status: %v", err)
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		panic(err)
	}

	return nil
}
