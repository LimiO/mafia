package server

import (
	"log"
	pgame "mafia/pkg/proto/game"

	connection "mafia/pkg/proto/connection"
)

type Filter func(userID string) bool

func (s *Server) AliveFilter(userID string) bool {
	user, ok := s.users[userID]
	if !ok {
		return true
	}
	return !user.alive
}

func (s *Server) SendResponse(rsp *connection.ServerResponse, from string, filters ...Filter) {
	for userId, user := range s.users {
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

func (s *Server) SendTo(userID, text string) {
	user, ok := s.users[userID]
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

func (s *Server) SendStateTo(userID string, state pgame.State) {
	user, ok := s.users[userID]
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

func (s *Server) SendToChat(userID, text string) {
	s.SendResponse(&connection.ServerResponse{
		Response: &connection.ServerResponse_Chat{
			Chat: &connection.ChatResponse{
				Text:   text,
				UserId: userID,
			},
		},
	}, userID)
}

func (s *Server) SendState(state pgame.State) {
	var filters []Filter
	if state != pgame.State_END {
		filters = append(filters, s.AliveFilter)
	}
	s.SendResponse(&connection.ServerResponse{
		Response: &connection.ServerResponse_State{
			State: &pgame.StateResponse{
				State: state,
			},
		},
	}, "", filters...)
}

func (s *Server) SendKillNotification(userID string) {
	s.SendResponse(&connection.ServerResponse{
		Response: &connection.ServerResponse_Kill{
			Kill: &pgame.KillResponse{
				UserId: userID,
			},
		},
	}, "")
}
