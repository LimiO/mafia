package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"mafia/pkg/proto/game"
	"mafia/server/status"

	"google.golang.org/grpc"

	connection "mafia/pkg/proto/connection"
)

const (
	MinPlayers       = 4
	MinCommitPlayers = 2
)

type User struct {
	conn  connection.MafiaServer_ConnectServer
	alive bool
	role  game.Role
}

type Server struct {
	Port  uint32
	Host  string
	users map[string]*User

	mu      sync.Mutex
	status  status.Status
	GrpcSrv *grpc.Server

	connection.UnimplementedMafiaServerServer
}

func (s *Server) Commit(_ context.Context, request *game.CommitRequest) (*game.CommitResponse, error) {
	if s.status.State != game.State_NIGHT {
		return nil, fmt.Errorf("not available action")
	}
	userID := request.GetUserId()
	role, ok := s.status.Roles[userID]
	if !ok {
		return nil, fmt.Errorf("failed to find this user in the game")
	}
	result := game.CommitResponse_OK
	switch role {
	case game.Role_HUMAN:
		return nil, fmt.Errorf("human can't commit")
	case game.Role_POLICE:
		targetRole, _ := s.status.Roles[request.GetTarget()]
		if targetRole != game.Role_MAFIA {
			result = game.CommitResponse_FAIL
		}
	}
	s.status.Commited[userID] = request.GetTarget()
	return &game.CommitResponse{
		Result: result,
	}, nil
}

func (s *Server) VoteBan(_ context.Context, request *game.VoteBanRequest) (*game.VoteBanResponse, error) {
	if s.status.State != game.State_DAY {
		return nil, fmt.Errorf("not available action")
	}
	userID := request.GetUserId()
	s.status.VoteBanned[userID] = request.GetTarget()
	return &game.VoteBanResponse{}, nil
}

func (s *Server) Chat(_ context.Context, request *game.ChatRequest) (*connection.ChatResponse, error) {
	if s.status.State != game.State_DAY {
		return nil, fmt.Errorf("not available action")
	}
	userID := request.GetUserId()
	s.SendToChat(userID, request.GetText())
	return &connection.ChatResponse{
		UserId: "Game",
		Text:   "message successfully send",
	}, nil
}

func (s *Server) End(_ context.Context, request *game.EndRequest) (*game.EndResponse, error) {
	if s.status.State != game.State_DAY {
		return nil, fmt.Errorf("not available action")
	}
	userID := request.GetUserId()
	s.status.Ended[userID] = true
	return &game.EndResponse{}, nil
}

func (s *Server) ListParticipants(
	_ context.Context,
	_ *connection.ListParticipantsRequest,
) (*connection.ListParticipantsResponse, error) {
	var users []string
	for key := range s.users {
		users = append(users, key)
	}
	return &connection.ListParticipantsResponse{
		Users: users,
	}, nil
}

func (s *Server) Connect(
	req *connection.UserJoinRequest,
	stream connection.MafiaServer_ConnectServer,
) error {
	var rspType connection.UserJoinResponse_Type

	s.mu.Lock()
	if _, ok := s.users[req.GetUserId()]; ok {
		rspType = connection.UserJoinResponse_EXISTS
	} else if len(s.users) == MinPlayers {
		rspType = connection.UserJoinResponse_STARTED
	} else {
		rspType = connection.UserJoinResponse_OK
		s.SendToChat(req.GetUserId(), "joined the status")
		s.users[req.GetUserId()] = &User{
			alive: true,
			conn:  stream,
		}
		log.Printf("user %q joined", req.GetUserId())
	}
	s.mu.Unlock()

	stream.Send(&connection.ServerResponse{
		Response: &connection.ServerResponse_Join{
			Join: &connection.UserJoinResponse{
				Type: rspType,
			},
		},
	})
	userID := req.GetUserId()
	for {
		select {
		case <-stream.Context().Done():
			if rspType == connection.UserJoinResponse_OK {
				delete(s.users, userID)
				if s.status.State == game.State_END {
					return nil
				}
				s.SendToChat(req.GetUserId(), "disconnected from the game")
				s.SendKillNotification(req.GetUserId())
				log.Printf("user %q disconnected", req.GetUserId())
			}
			return nil
		default:
			continue
		}
	}
}

func MakeServer() (*Server, error) {
	return &Server{
		Port:  9000,
		Host:  "",
		users: make(map[string]*User),
		status: status.Status{
			VoteBanned: map[string]string{},
			Ended:      map[string]bool{},
			Commited:   map[string]string{},
			Roles:      map[string]game.Role{},
		},
	}, nil
}

func (s *Server) StartListen() error {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return fmt.Errorf("failed to make tcp connection to port %d: %v", s.Port, err)
	}

	srv := grpc.NewServer()
	s.GrpcSrv = srv

	connection.RegisterMafiaServerServer(srv, s)
	log.Println("server started!")
	err = srv.Serve(conn)
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

func (s *Server) Stop() {
	s.GrpcSrv.Stop()
}
