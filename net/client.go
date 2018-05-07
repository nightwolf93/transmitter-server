package net

import (
	"github.com/gorilla/websocket"
)

// Client is the struct that handle a client
type Client struct {
	ID    string
	Type  int
	Alive bool

	// Specific connection
	WSConn               *websocket.Conn
	WSReceiveMessageChan chan []byte
}
