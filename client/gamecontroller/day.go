package gamecontroller

import (
	"context"
	"fmt"
	connection "mafia/pkg/proto/connection"
	pgame "mafia/pkg/proto/game"
)

func (c *Controller) GetDayOptions() []string {
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
		selected := c.SelectAction("Select target to do", c.GetDayOptions())
		switch selected {
		case chatOption:
			c.GetAndSendMessage()
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

func (c *Controller) GetAndSendMessage() {
	msg := c.AskInput("Type message to send")
	c.ChatChan <- fmt.Sprintf("[%s]: %s", c.ID, msg)
}

func (c *Controller) SelectAndVoteBan(client connection.MafiaServerClient) error {
	users := c.makeSelectParticipants()
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
