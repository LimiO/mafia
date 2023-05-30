package controller

import (
	"context"
	"fmt"
	"log"

	"mafia/client/cli"
	connection "mafia/pkg/proto/connection"
	pgame "mafia/pkg/proto/game"
)

func (c *Controller) GetOptions() []string {
	options := []string{chatOption, endDayOption}
	if c.DayNumber == 0 {
		return options
	}

	options = append(options, votebanOption)

	if c.Role.GetInfo() != "" {
		options = append(options, publishInfoOption)
	}
	return options
}

func (c *Controller) ProcessDay(client connection.MafiaServerClient) {
	for {
		if c.State != pgame.State_DAY {
			return
		}
		selected, err := cli.AskSelect(
			"Select action to do",
			c.GetOptions(),
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
		case publishInfoOption:
			err = c.PublishInfo(client)
		}
	}
}

func (c *Controller) ProcessSpirit() {
}

func (c *Controller) GetAndSendMessage(client connection.MafiaServerClient) error {
	msg, err := cli.AskInput("Type message to send")
	if err != nil {
		return fmt.Errorf("failed to ask input: %v", err)
	}
	_, err = client.Chat(context.Background(), &pgame.ChatRequest{
		UserId: c.ID,
		Text:   msg,
		Game:   c.GameID,
	})
	return err
}

func (c *Controller) SelectAndVoteBan(client connection.MafiaServerClient) error {
	var users []string
	for participantID, participant := range c.Participants {
		if !participant.Alive || participantID == c.ID {
			continue
		}
		users = append(users, participantID)
	}
	msg, err := cli.AskSelect("Select target to vote:", users)
	if err != nil {
		return fmt.Errorf("failed to ask input: %v", err)
	}

	_, err = client.VoteBan(context.Background(), &pgame.VoteBanRequest{
		UserId: c.ID,
		Target: msg,
		Game:   c.GameID,
	})
	return err
}

func (c *Controller) EndDay(client connection.MafiaServerClient) error {
	c.State = pgame.State_SPIRIT
	c.DayNumber++
	_, err := client.End(context.Background(), &pgame.EndRequest{
		UserId: c.ID,
		Game:   c.GameID,
	})
	return err
}

func (c *Controller) PublishInfo(client connection.MafiaServerClient) error {
	_, err := client.Publish(context.Background(), &pgame.PublishRequest{
		UserId: c.ID,
		Info:   c.Role.GetInfo(),
		Game:   c.GameID,
	})

	return err
}
