package client

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"mafia/client/gamecontroller"
	"mafia/internal/queue"
	connection "mafia/pkg/proto/connection"
	"mafia/pkg/proto/game"
	"mafia/roles"
	"os"
)

type Client struct {
	ServerPort uint32
	GameCtl    *gamecontroller.Controller
	QueueCtl   *queue.Controller

	GrpcClient connection.MafiaServerClient
}

func MakeClient(id string, password string) (*Client, error) {
	chatChan := make(chan string, 100)
	return &Client{
		ServerPort: 9000,
		GameCtl: &gamecontroller.Controller{
			ID:           id,
			Password:     password,
			Participants: map[string]*gamecontroller.Participant{},
			ChatChan:     chatChan,
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
		os.Exit(0)
		return nil
	})

	eg.Go(func() error {
		for {
			if c.QueueCtl != nil {
				break
			}
		}
		err = c.QueueCtl.StartConsume(func(delivery amqp.Delivery) error {
			chatType := delivery.RoutingKey
			if delivery.RoutingKey == queue.MafiaKey {
				if _, ok := c.GameCtl.Role.(*roles.Mafia); !ok {
					return nil
				}
			}
			log.Printf("Got message from chat %q: %s", chatType, delivery.Body)
			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to consume: %v", err)
		}
		return nil
	})

	eg.Go(func() error {
		for {
			if c.QueueCtl != nil {
				break
			}
		}
		for msg := range c.QueueCtl.ChatChan {
			_, ok := c.GameCtl.Role.(*roles.Mafia)
			key := queue.AllKey
			if c.GameCtl.State == game.State_NIGHT && ok {
				key = queue.MafiaKey
			}
			err = c.QueueCtl.Push("server", key, msg)
			if err != nil {
				fmt.Printf("failed to push msg: %v", err)
			}
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		panic(err)
	}

	return nil
}
