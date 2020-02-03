package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/go-cmd/cmd"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type runTicket struct {
	ID       string
	Settings *runSettings
	Username string
	Label    string
	Done     bool
	Err      error
}

type runSettings struct {
	Account  *instagramAccount
	Settings *botSettings
}

type bot struct {
	m   sync.Mutex
	id  string
	cmd *cmd.Cmd
}

func newBot() *bot {
	return &bot{}
}

func (b *bot) run(label, username string, settings *runSettings) (*runTicket, error) {
	b.m.Lock()
	if b.cmd != nil {
		b.m.Unlock()
		return nil, errors.New("Already running")
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "*")
	defer tmpFile.Close()

	if err != nil {
		return nil, err
	}

	jsonSettings, _ := json.Marshal(settings)
	tmpFile.Write(jsonSettings)

	b.cmd = cmd.NewCmd("python3", "../main.py", tmpFile.Name())

	ticket := &runTicket{
		ID:       uuid.New().String(),
		Done:     false,
		Label:    label,
		Username: username,
		Settings: settings,
	}

	go func(t runTicket, statusChan <-chan cmd.Status) {
		for {
			s := <-statusChan

			if s.Error != nil {
				logrus.Errorf("Execution error: %v", s.Error)
			}

			if s.Exit == 1 {
				logrus.Errorf("There was an error:\n%v", strings.Join(s.Stderr, "\n"))
			}

			if s.Complete {
				break
			}
		}
		logrus.Infof("[BOT] Finished with ticket %s", t.ID)
		b.m.Lock()
		b.cmd = nil
		b.m.Unlock()
	}(*ticket, b.cmd.Start())

	b.m.Unlock()
	return ticket, nil
}
