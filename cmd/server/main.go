package main

import (
	"fmt"

	"mafia/server"
)

func main() {
	srv, err := server.MakeServer()
	if err != nil {
		panic(fmt.Errorf("failed to make server: %v", err))
	}
	err = srv.Start()
	if err != nil {
		panic(fmt.Errorf("failed to start server listen: %v", err))
	}
}
