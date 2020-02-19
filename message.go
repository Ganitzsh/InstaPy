package main

type clientMessage struct {
	code string
	data interface{}
}

func newClientMessage() *clientMessage {
	return &clientMessage{}
}
