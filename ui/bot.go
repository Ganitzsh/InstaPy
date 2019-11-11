package main

import (
	"errors"
	"sync"

	"github.com/go-cmd/cmd"
	"github.com/sirupsen/logrus"
)

type runTicket struct {
	id       string
	botID    string
	settings *runSettings
}

type runSettings struct {
	account  *instagramAccount
	settings *botSettings
}

type bot struct {
	m   sync.Mutex
	id  string
	cmd *cmd.Cmd
}

func newBot() *bot {
	return &bot{}
}

func (b *bot) run(settings *runSettings) (*runTicket, error) {
	b.m.Lock()
	if b.cmd != nil {
		b.m.Unlock()
		return nil, errors.New("Already running")
	}

	b.cmd = cmd.NewCmd("python3", "../main.py")
	statusChan := b.cmd.Start()

	ticket := &runTicket{
		id:       "123",
		botID:    b.id,
		settings: settings,
	}

	go func() {
		<-statusChan
		logrus.WithField("id", b.id).Infof("[BOT] Finished with ticket %s", ticket.id)
		b.m.Lock()
		b.cmd = nil
		b.m.Unlock()
	}()

	return ticket, nil
}
