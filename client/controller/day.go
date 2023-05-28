package controller

import (
	"context"
	"fmt"
	"log"

	"mafia/client/cli"
	connection "mafia/pkg/proto/connection"
	pgame "mafia/pkg/proto/game"
)

func (c *Controller) ProcessDay(client connection.MafiaServerClient) {
	for {
		if c.State != pgame.State_DAY {
			return
		}
		selected, err := cli.AskSelect(
			"Выберите действие, которое хотите совершить",
			[]string{chatOption, votebanOption, endDayOption},
		)
		if err != nil {
			log.Printf("failed to process day: %v", err)
			return
		}
		switch selected {
		case chatOption:
			err = c.GetAndSendMessage(client)
		case votebanOption:
			err = c.SelectAndVoteBan(client)
		case endDayOption:
			err = c.EndDay(client)
			return
		}
	}
}

func (c *Controller) ProcessSpirit() {
}

func (c *Controller) GetAndSendMessage(client connection.MafiaServerClient) error {
	msg, err := cli.AskInput("Введите сообщения для продолжения")
	if err != nil {
		return fmt.Errorf("failed to ask input: %v", err)
	}
	_, err = client.Chat(context.Background(), &pgame.ChatRequest{UserId: c.ID, Text: msg})
	return err
}

func (c *Controller) SelectAndVoteBan(client connection.MafiaServerClient) error {
	var users []string
	for participantID := range c.Participants {
		users = append(users, participantID)
	}
	msg, err := cli.AskSelect("Выберите цель для продолжения:", users)
	if err != nil {
		return fmt.Errorf("failed to ask input: %v", err)
	}

	_, err = client.VoteBan(context.Background(), &pgame.VoteBanRequest{UserId: c.ID, Target: msg})
	return err
}

func (c *Controller) EndDay(client connection.MafiaServerClient) error {
	c.State = pgame.State_SPIRIT
	_, err := client.End(context.Background(), &pgame.EndRequest{UserId: c.ID})
	return err
}
