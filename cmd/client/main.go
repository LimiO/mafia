package main

import (
	"fmt"
	"os"

	"mafia/client"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Errorf("name param must be given"))
	}

	grpcClient, err := client.MakeClient(os.Args[1])
	if err != nil {
		panic(fmt.Errorf("failed to make client: %v", err))
	}
	err = grpcClient.StartSession()
	if err != nil {
		panic(fmt.Errorf("failed to join status: %v", err))
	}
}
