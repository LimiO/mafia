package server

import (
	"context"
	"fmt"
	"log"
	"mafia/pkg/proto/game"
	"net"

	"google.golang.org/grpc"

	connection "mafia/pkg/proto/connection"
)

const (
	MinPlayers = 4
)

type User struct {
	conn  connection.MafiaServer_ConnectServer
	alive bool
	role  game.Role
}

var GlobalGameID uint32 = 0

type Server struct {
	Port uint32

	Games map[uint32]*Game

	GrpcSrv *grpc.Server

	connection.UnimplementedMafiaServerServer
}

func (s *Server) Publish(_ context.Context, req *game.PublishRequest) (*game.PublishResponse, error) {
	curGame := s.Games[req.Game]

	if curGame.status.State != game.State_DAY {
		return nil, fmt.Errorf("not available action")
	}
	userID := req.GetUserId()
	info := req.GetInfo()
	curGame.SendToChat("game", fmt.Sprintf("user %q has info: %v", userID, info))
	return &game.PublishResponse{}, nil
}

func (s *Server) Commit(_ context.Context, req *game.CommitRequest) (*game.CommitResponse, error) {
	curGame := s.Games[req.Game]

	if curGame.status.State != game.State_NIGHT {
		return nil, fmt.Errorf("not available action")
	}
	userID := req.GetUserId()
	role, ok := curGame.status.Roles[userID]
	if !ok {
		return nil, fmt.Errorf("failed to find this user in the game")
	}
	result := game.CommitResponse_FAIL
	if role == game.Role_POLICE {
		targetRole, _ := curGame.status.Roles[req.GetTarget()]
		if targetRole == game.Role_MAFIA {
			result = game.CommitResponse_OK
		}
	}
	curGame.status.Commited[userID] = req.GetTarget()
	return &game.CommitResponse{
		Result: result,
	}, nil
}

func (s *Server) VoteBan(_ context.Context, req *game.VoteBanRequest) (*game.VoteBanResponse, error) {
	curGame := s.Games[req.Game]
	if curGame.status.State != game.State_DAY {
		return nil, fmt.Errorf("not available action")
	}
	userID := req.GetUserId()
	curGame.status.VoteBanned[userID] = req.GetTarget()
	return &game.VoteBanResponse{}, nil
}

func (s *Server) Chat(_ context.Context, req *game.ChatRequest) (*connection.ChatResponse, error) {
	curGame := s.Games[req.Game]
	if curGame.status.State != game.State_DAY {
		return nil, fmt.Errorf("not available action")
	}
	//userID := req.GetUserId()
	//curGame.SendToChat(userID, req.GetText())
	return &connection.ChatResponse{
		UserId: "Game",
		Text:   "message successfully send",
	}, nil
}

func (s *Server) End(_ context.Context, req *game.EndRequest) (*game.EndResponse, error) {
	curGame := s.Games[req.Game]
	if curGame.status.State != game.State_DAY {
		return nil, fmt.Errorf("not available action")
	}
	userID := req.GetUserId()
	curGame.status.Ended[userID] = true
	return &game.EndResponse{}, nil
}

func (s *Server) ListParticipants(
	_ context.Context,
	req *connection.ListParticipantsRequest,
) (*connection.ListParticipantsResponse, error) {
	curGame := s.Games[req.Game]

	var users []string
	for key := range curGame.users {
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

	var curGame *Game
	for _, g := range s.Games {
		if !g.status.Started {
			curGame = g
			break
		}
	}
	if curGame == nil {
		var err error
		curGame, err = NewGame()
		if err != nil {
			return fmt.Errorf("failed to make game: %v", err)
		}
		s.Games[curGame.gameID] = curGame
	}

	curGame.mu.Lock()
	if _, ok := curGame.users[req.GetUserId()]; ok {
		rspType = connection.UserJoinResponse_EXISTS
	} else if curGame.status.Started {
		rspType = connection.UserJoinResponse_STARTED
	} else {
		rspType = connection.UserJoinResponse_OK
		err := curGame.QueueCtl.AddProducer(fmt.Sprintf("client.%d.%s", curGame.gameID, req.GetUserId()))
		if err != nil {
			return fmt.Errorf("failed to add producer")
		}
		curGame.SendToChat(req.GetUserId(), "joined the game")
		curGame.users[req.GetUserId()] = &User{
			alive: true,
			conn:  stream,
		}
		log.Printf("user %q joined", req.GetUserId())
	}
	curGame.mu.Unlock()

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
				curGame.DeleteUser(userID)
			}
			return nil
		default:
			continue
		}
	}
}

func MakeServer() (*Server, error) {
	g, err := NewGame()
	if err != nil {
		return nil, fmt.Errorf("failed to make game: %v", err)
	}
	server := &Server{
		Port: 9000,
		Games: map[uint32]*Game{
			g.gameID: g,
		},
	}
	srv := grpc.NewServer()
	server.GrpcSrv = srv
	connection.RegisterMafiaServerServer(server.GrpcSrv, server)
	return server, nil
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
