package gamecontroller

func (c *Controller) makeParticipantIds() []string {
	var ids []string
	for userID := range c.Participants {
		if userID == c.ID {
			continue
		}
		ids = append(ids, userID)
	}
	return ids
}

func (c *Controller) makeAliveParticipantIds() []string {
	var ids []string
	for userID, p := range c.Participants {
		if userID == c.ID {
			continue
		}
		if !p.Alive {
			continue
		}
		ids = append(ids, userID)
	}
	return ids
}

func (c *Controller) makeSelectParticipants() []string {
	var ids []string
	for userID, p := range c.Participants {
		if !p.Alive || userID == c.ID {
			continue
		}
		ids = append(ids, userID)
	}
	return ids
}
