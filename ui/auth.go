package main

import (
	"errors"

	"github.com/google/uuid"
)

type token struct {
	userID       string
	accessToken  string
	refreshToken string
}

func newToken() *token {
	return &token{}
}

func authenticateUser(username, password string) (*token, error) {
	user, err := uRepo.findUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if user.password != password {
		return nil, errors.New("Invalid username or password")
	}

	t := newToken()
	t.userID = user.id
	t.accessToken = uuid.New().String()
	t.refreshToken = uuid.New().String()

	return t, nil
}

func verifyToken(tok *token) error {
	for _, ts := range tokens {
		for _, t := range ts {
			if t.accessToken == tok.accessToken {
				return nil
			}
		}
	}
	return errors.New("Token not found")
}
