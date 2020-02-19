package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-cmd/cmd"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func now() *time.Time {
	now := time.Now()
	return &now
}

type runTicket struct {
	ID        string
	StartDate *time.Time
	EndDate   *time.Time
	Settings  *runSettings
	Username  string
	Label     string
	Done      bool
	Logs      []string
	ErrLogs   []string
	Err       error
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

func (b *bot) run(label, username, igpassword string, settings *runSettings) (*runTicket, error) {
	b.m.Lock()
	if b.cmd != nil {
		b.m.Unlock()
		return nil, errors.New("Already running")
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "*")
	defer tmpFile.Close()

	logrus.Println(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	settingsCpy := *settings
	settingsCpy.Account.Password = igpassword

	jsonSettings, _ := json.Marshal(settingsCpy)
	spew.Dump(settingsCpy)
	tmpFile.Write(jsonSettings)

	b.cmd = cmd.NewCmd("python3", "./main.py", tmpFile.Name())

	ticket := &runTicket{
		ID:        uuid.New().String(),
		StartDate: now(),
		Done:      false,
		Label:     label,
		Username:  username,
		Settings:  settings,
	}

	go func(t runTicket, botCmd *cmd.Cmd) {
		for {
			if botCmd == nil {
				break
			}

			s := botCmd.Status()

			globMut.Lock()

			if tickets[t.ID] != nil {
				tickets[t.ID].Logs = s.Stdout
				tickets[t.ID].ErrLogs = s.Stderr
			}

			globMut.Unlock()

			if s.Complete {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}(*ticket, b.cmd)

	go func(t runTicket, statusChan <-chan cmd.Status) {
		var err error
		var logs, errLogs []string

		for {
			s := <-statusChan

			if s.Error != nil {
				err = s.Error
				logrus.Errorf("Execution error: %v", s.Error)
			}

			if s.Exit == 1 {
				logrus.Errorf("There was an error:\n%v", strings.Join(s.Stderr, "\n"))
			}

			if s.Complete {
				logs = s.Stdout
				errLogs = s.Stderr
				break
			}
		}
		logrus.Infof("[BOT] Finished with ticket %s", t.ID)
		b.m.Lock()
		globMut.Lock()
		b.cmd = nil
		if tickets[t.ID] != nil {
			tickets[t.ID].Done = true
			tickets[t.ID].Err = err
			tickets[t.ID].EndDate = now()
			tickets[t.ID].Logs = logs
			tickets[t.ID].ErrLogs = errLogs
		}
		globMut.Unlock()
		b.m.Unlock()
	}(*ticket, b.cmd.Start())

	b.m.Unlock()
	return ticket, nil
}
