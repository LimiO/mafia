package main

import (
	"fmt"
	"os"

	"mafia/client"
	"mafia/client/cli"
	"mafia/internal"
)

func main() {
	isAuto := len(os.Args) > 1 && os.Args[1] == "auto"
	var name string
	if isAuto {
		name = internal.RandStringRunes(5)
	} else {
		name = cli.AskInput("Enter your name")
	}
	mafiaClient, err := client.MakeClient(name)
	if err != nil {
		panic(fmt.Errorf("failed to make client: %v", err))
	}
	mafiaClient.GameCtl.IsAuto = isAuto
	err = mafiaClient.StartSession()
	if err != nil {
		panic(fmt.Errorf("failed to join status: %v", err))
	}
}
