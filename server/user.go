package server

import (
	"fmt"

	"mafia/internal/db"
	"mafia/internal/helpers"
)

// TODO тут добавить возврат пользователя и затем сохранение его в обычный User
// TODO добавить создание stats при создании юзера.

func (s *Server) AuthorizeOrRegisterUser(ID string, pass string) error {
	user, err := s.DBManager.GetUser(ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}
	hash := helpers.Hash(pass)
	if user != nil {
		if user.PassHash == hash {
			return nil
		}
		return fmt.Errorf("user exists and password not equal")
	}
	user = &db.User{
		ID:       ID,
		PassHash: hash,
	}
	err = s.DBManager.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}
