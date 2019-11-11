package main

import "github.com/google/uuid"

type instagramAccount struct {
	id       string
	label    string
	username string
	password string
}

func newInstagramAccount() *instagramAccount {
	return &instagramAccount{}
}

type user struct {
	id                string
	username          string
	password          string
	instagramAccounts []*instagramAccount
	settings          *botSettings
}

func newUser() *user {
	return &user{}
}

func saveUser(u *user) (*user, error) {
	if u.id == "" {
		u.id = uuid.New().String()
	}
	if u.instagramAccounts != nil {
		for _, account := range u.instagramAccounts {
			if account.id == "" {
				account.id = uuid.New().String()
			}
		}
	}
	users[u.id] = u
	return u, nil
}
