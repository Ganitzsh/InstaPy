package main

import (
	"errors"

	"github.com/google/uuid"
)

type token struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

func newToken() *token {
	return &token{}
}

func authenticateUser(username, password string) (*token, *user, error) {
	user, err := uRepo.findUserByUsername(username)
	if err != nil {
		return nil, nil, err
	}

	if user.Password != password {
		return nil, nil, errors.New("Invalid username or password")
	}

	t := newToken()
	t.UserID = user.ID
	t.AccessToken = uuid.New().String()
	t.RefreshToken = uuid.New().String()

	return t, user, nil
}

func verifyToken(tok *token) error {
	if users[tok.AccessToken] == nil {
		return errors.New("Token not found")
	}

	return nil
}
