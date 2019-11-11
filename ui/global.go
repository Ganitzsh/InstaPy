package main

var clients = make(map[string]*client)
var tickets = make(map[string]*runTicket)
var bots = make(map[string]*bot)
var users = make(map[string]*user)
