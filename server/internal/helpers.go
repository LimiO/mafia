package internal

import (
	game "mafia/pkg/proto/game"
	"math/rand"
)

func ShuffleRoles(users []string) map[string]game.Role {
	for i := range users {
		j := rand.Intn(i + 1)
		users[i], users[j] = users[j], users[i]
	}
	result := map[string]game.Role{
		//users[0]: game.Role_HUMAN,
		//users[1]: game.Role_HUMAN,
		users[0]: game.Role_MAFIA,
		users[1]: game.Role_POLICE,
	}
	return result
}
