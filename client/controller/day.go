package controller

import (
	"context"
	"fmt"
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
		var err error
		if c.State != pgame.State_DAY {
			return
		}
		selected := c.SelectAction("Select target to do", c.GetOptions())
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
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (c *Controller) ProcessSpirit() {
}

func (c *Controller) GetAndSendMessage(client connection.MafiaServerClient) error {
	msg := c.AskInput("Type message to send")
	_, err := client.Chat(context.Background(), &pgame.ChatRequest{
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
	target := c.SelectAction("Select target to vote:", users)

	_, err := client.VoteBan(context.Background(), &pgame.VoteBanRequest{
		UserId: c.ID,
		Target: target,
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
