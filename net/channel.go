package net

import (
	"encoding/json"
	"fmt"

	"github.com/nightwolf93/transmitter-server/common"
	"github.com/nightwolf93/transmitter-server/net/protocol"
	"github.com/nightwolf93/transmitter-server/storage"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	// The version of the channel used by the server
	ChannelVersionNumber = 1
)

// Channel a channel is a group of client that can communicate each other and can receive broadcasted data
type Channel struct {
	ID       string
	Name     string
	Password string
	Version  int

	Subscribers map[*Client]*ChannelSubscriber
}

// ChannelSubscriber client who subscribe to this channel
type ChannelSubscriber struct {
	RoutingKeys []string
	Client      *Client
	Channel     *Channel
}

// HasRoutingKey check if the subscriber has the routing key
func (sub *ChannelSubscriber) HasRoutingKey(key string) bool {
	for _, e := range sub.RoutingKeys {
		if e == key {
			return true
		}
	}
	return false
}

var channels map[string]*Channel = make(map[string]*Channel)

// GetOrNewChannel get or create a new channel
func GetOrNewChannel(name string, password string) *Channel {
	channel, found := channels[name]
	if !found {
		channel = &Channel{
			ID:          common.GenerateLongUniqueID(),
			Name:        name,
			Password:    password,
			Subscribers: make(map[*Client]*ChannelSubscriber),
			Version:     ChannelVersionNumber,
		}

		// Save the channel because this is a new channel
		channel.Save()
		channels[name] = channel
	}

	return channel
}

// RegisterClient add a client to channel
func (channel *Channel) RegisterClient(client *Client, routingKeys []string) *ChannelSubscriber {
	sub := &ChannelSubscriber{
		Client:      client,
		RoutingKeys: routingKeys,
		Channel:     channel,
	}
	channel.Subscribers[client] = sub
	log.Debugf("Client %s registered to the channel %s (count: %d)", client.UID, channel.Name, len(channel.Subscribers))
	return sub
}

// GetSubscribersWithRoutingKeys get client who subscribe to this channel and matching routing keys
func (channel *Channel) GetSubscribersWithRoutingKeys(routingKeys []string) []*ChannelSubscriber {
	targets := make([]*ChannelSubscriber, 0)
	for _, sub := range channel.Subscribers {
		// Check if the sub has no routing keys set and the routing keys requested in empty, if yes select all the subbscribers
		if len(routingKeys) == 0 && len(sub.RoutingKeys) == 0 {
			targets = append(targets, sub)
			continue
		}

		// Check of all keys for the client
		for _, key := range routingKeys {
			if sub.HasRoutingKey(key) {
				targets = append(targets, sub)
				break
			}
		}
	}
	return targets
}

// BroadcastEvent broadcast data to matching keys
func (channel *Channel) BroadcastEvent(eventName string, data []byte, routingKeys []string) {
	// Find all targets by routing keys
	targets := channel.GetSubscribersWithRoutingKeys(routingKeys)

	message := protocol.NewCustomDataEvent("", eventName, data)

	// Broadcast data to all targets
	for _, sub := range targets {
		sub.Client.SendMessage <- message
	}

	log.Debugf("Event %s broadcasted in the channel %s to %d subscriber(s)", eventName, channel.Name, len(channel.Subscribers))
}

// UnRegisterClient unregister a client from this channel
func (channel *Channel) UnRegisterClient(client *Client) {
	delete(channel.Subscribers, client)

	//TODO: Broadcast peer disconnection to the channel
	log.Debugf("Client %s unregistered from the channel %s (count: %d)", client.UID, channel.Name, len(channel.Subscribers))
}

// ToSerialize return the channel to serializable format
func (channel *Channel) ToSerialize() map[string]interface{} {
	return map[string]interface{}{
		"ID":       channel.ID,
		"Name":     channel.Name,
		"Password": channel.Password,
		"Version":  channel.Version,
	}
}

// Save save the channel into the database
func (channel *Channel) Save() {
	// Build the data with be saved
	dat, err := json.Marshal(channel.ToSerialize())
	if err != nil {
		log.Errorf("Can't serialize the channel: %s", err)
		return
	}

	db := storage.GetDB()
	db.Put([]byte(fmt.Sprintf("/channel/%s", channel.Name)), dat, nil)
	log.Debugf("Channel %s saved into database", channel.Name)
}

// LoadChannelsFromDatabase load all channels from the database
func LoadChannelsFromDatabase() {
	db := storage.GetDB()
	iter := db.NewIterator(util.BytesPrefix([]byte("/channel/")), nil)
	for iter.Next() {
		channel := &Channel{
			Subscribers: make(map[*Client]*ChannelSubscriber),
		}
		json.Unmarshal(iter.Value(), channel)

		channels[channel.Name] = channel
		log.Debugf("Channel %s loaded from database", channel.Name)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Errorf("Can't load channels: %s", err)
	}
}
