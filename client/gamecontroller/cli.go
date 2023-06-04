package gamecontroller

import (
	"fmt"
	"math/rand"
	"time"

	"mafia/client/cli"
	"mafia/internal"
)

func (c *Controller) SelectAction(msg string, options []string) string {
	if c.IsAuto {
		rand.Seed(time.Now().UnixNano())
		selected := options[rand.Intn(100001)%len(options)]
		fmt.Printf("Selected random option to msg \"%s...\": %q\n", msg[:10], selected)
		return selected
	}
	return cli.AskSelect(msg, options)
}

func (c *Controller) AskInput(msg string) string {
	if c.IsAuto {
		rand.Seed(time.Now().UnixNano())
		result := internal.RandStringRunes(10)
		fmt.Printf("message to send: %q\n", result)
		return result
	}
	return cli.AskInput(msg)
}
