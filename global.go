package main

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	globMut sync.Mutex
	clients = make(map[string]*client)
	tickets = make(map[string]*runTicket)
	users   = make(map[string]*user)
	globBot *bot

	mongoDBClient *mongo.Client
	uRepo         *userRepository
)
