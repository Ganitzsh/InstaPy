package main

import "go.mongodb.org/mongo-driver/mongo"

var (
	clients = make(map[string]*client)
	tickets = make(map[string]*runTicket)
	bots    = make(map[string]*bot)
	users   = make(map[string]*user)

	tokens        = make(map[string][]*token)
	mongoDBClient *mongo.Client
	uRepo         *userRepository
)
