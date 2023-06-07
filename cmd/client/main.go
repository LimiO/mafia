package main

import (
	"fmt"
	"os"

	"mafia/client"
	"mafia/client/cli"
	"mafia/internal/helpers"
)

func main() {
	isAuto := len(os.Args) > 1 && os.Args[1] == "auto"
	var name string
	var password string
	if isAuto {
		name = helpers.RandStringRunes(5)
		password = helpers.RandStringRunes(5)
	} else {
		name = cli.AskInput("Enter your name")
		password = cli.AskInput("Enter your password")
	}
	mafiaClient, err := client.MakeClient(name, password)
	if err != nil {
		panic(fmt.Errorf("failed to make client: %v", err))
	}
	mafiaClient.GameCtl.IsAuto = isAuto
	err = mafiaClient.StartSession()
	if err != nil {
		panic(fmt.Errorf("failed to join status: %v", err))
	}
}
