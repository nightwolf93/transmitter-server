package net

import "github.com/nightwolf93/transmitter-server/common"

// Channel a channel is a group of client that can communicate each other and can receive broadcasted data
type Channel struct {
	ID       string
	Name     string
	Password string

	Subscribers map[*Client]bool
}

var channels map[string]*Channel = make(map[string]*Channel)

// GetOrNewChannel get or create a new channel
func GetOrNewChannel(name string) *Channel {
	channel, found := channels[name]
	if !found {
		channel = &Channel{
			ID:          common.GenerateLongUniqueID(),
			Name:        name,
			Subscribers: make(map[*Client]bool),
		}
	}

	return channel
}

// RegisterClient add a client to channel
func (channel *Channel) RegisterClient(client *Client) {
	channel.Subscribers[client] = true
}
