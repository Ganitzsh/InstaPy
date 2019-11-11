package main

import "time"

type client struct {
	id       string
	userID   string
	lastPing *time.Time
}

func (c *client) newClient() *client {
	return &client{}
}

func (c *client) sendMessage(msg clientMessage) error {
	return nil
}
