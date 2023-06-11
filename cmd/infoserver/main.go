package main

import (
	"fmt"
	"mafia/internal/db"
	"mafia/users"
)

func main() {
	cfg := &users.Config{
		Port: 8090,
		DBConfig: &db.Config{
			DBName: "mafia.db",
		},
	}
	ctl, err := users.MakeController(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to make controller: %v", err))
	}
	err = ctl.StartServer()
	if err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}
