package main

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
