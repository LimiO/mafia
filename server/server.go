package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"mafia/internal/db"
	connection "mafia/pkg/proto/connection"
	"mafia/pkg/proto/game"
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

	GrpcSrv   *grpc.Server
	DBManager *db.Manager

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
	if err := s.AuthorizeOrRegisterUser(req.GetUserId(), req.GetPassword()); err != nil {
		return fmt.Errorf("failed to authorize user: %v", err)
	}

	curGame, err := s.GetOrCreateGame()
	if err != nil {
		return fmt.Errorf("failed to get or create new game: %v", err)
	}

	err = curGame.QueueCtl.AddProducer(fmt.Sprintf("client.%d.%s", curGame.gameID, req.GetUserId()))
	if err != nil {
		return fmt.Errorf("failed to add producer")
	}
	curGame.SendToChat(req.GetUserId(), "joined the game")
	curGame.users[req.GetUserId()] = &User{
		alive: true,
		conn:  stream,
	}
	log.Printf("user %q joined", req.GetUserId())

	stream.Send(&connection.ServerResponse{
		Response: &connection.ServerResponse_Join{
			Join: &connection.UserJoinResponse{
				Type: connection.UserJoinResponse_OK,
			},
		},
	})
	userID := req.GetUserId()
	for {
		select {
		case <-stream.Context().Done():
			curGame.DeleteUser(userID)
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
	dbCfg := &db.Config{
		DBName: "mafia.db",
	}
	mgr, err := db.NewManager(dbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %v", err)
	}

	server := &Server{
		Port: 9000,
		Games: map[uint32]*Game{
			g.gameID: g,
		},
		DBManager: mgr,
	}
	srv := grpc.NewServer()
	server.GrpcSrv = srv
	connection.RegisterMafiaServerServer(server.GrpcSrv, server)
	return server, nil
}

func (s *Server) StartListen() error {
	conn, err := net.Listen("tcp", fmt.Sprintf("server:%d", s.Port))
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
