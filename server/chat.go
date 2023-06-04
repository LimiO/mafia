package server

import (
	"fmt"
	"log"
	connection "mafia/pkg/proto/connection"
	pgame "mafia/pkg/proto/game"
)

func (g *Game) SendResponse(rsp *connection.ServerResponse, from string, filters ...Filter) {
	for userId, user := range g.users {
		if from == userId {
			continue
		}
		needProcess := true
		for _, filter := range filters {
			if filter(userId) {
				needProcess = false
				break
			}
		}
		if !needProcess {
			continue
		}
		err := user.conn.Send(rsp)
		if err != nil {
			log.Printf("failed to send message to user %q: %v", userId, err)
		}
	}
}

func (g *Game) SendTo(userID, text string) {
	user, ok := g.users[userID]
	if !ok {
		return
	}
	_ = user.conn.Send(&connection.ServerResponse{
		Response: &connection.ServerResponse_Chat{
			Chat: &connection.ChatResponse{
				Text:   text,
				UserId: userID,
			},
		},
	})
}

func (g *Game) SendStateTo(userID string, state pgame.State) {
	user, ok := g.users[userID]
	if !ok {
		return
	}
	_ = user.conn.Send(&connection.ServerResponse{
		Response: &connection.ServerResponse_State{
			State: &pgame.StateResponse{
				State: state,
			},
		},
	})
}

func (g *Game) SendToChat(userID, text string) {
	msg := fmt.Sprintf("[%s]: %s", userID, text)
	for playerID := range g.users {
		key := fmt.Sprintf("client.%d.%s", g.gameID, playerID)
		err := g.QueueCtl.Push(key, "all", msg)
		if err != nil {
			fmt.Printf("failed to push message to queue: %v", err)
		}
	}
	g.SendResponse(&connection.ServerResponse{
		Response: &connection.ServerResponse_Chat{
			Chat: &connection.ChatResponse{
				Text:   text,
				UserId: userID,
			},
		},
	}, userID)
}

func (g *Game) SendState(state pgame.State) {
	var filters []Filter
	if state != pgame.State_END {
		filters = append(filters, g.AliveFilter)
	}
	g.SendResponse(&connection.ServerResponse{
		Response: &connection.ServerResponse_State{
			State: &pgame.StateResponse{
				State: state,
			},
		},
	}, "", filters...)
}

func (g *Game) SendKillNotification(userID string) {
	g.SendResponse(&connection.ServerResponse{
		Response: &connection.ServerResponse_Kill{
			Kill: &pgame.KillResponse{
				UserId: userID,
			},
		},
	}, "")
}
