package net

import (
	"github.com/golang/protobuf/proto"
	"github.com/nightwolf93/transmitter-server/common"
	"github.com/nightwolf93/transmitter-server/net/protobuf"
	"github.com/nightwolf93/transmitter-server/net/protocol"
	log "github.com/sirupsen/logrus"
)

// InitHandlers init client handlers
func (client *Client) InitHandlers() {
	client.Handlers = make(map[int]func(string, []byte))
	client.Handlers[protocol.HandshakeResponse] = client.handleHandshakeResponse
	client.Handlers[protocol.SubscribeToChannelRequest] = client.handleSubscribeToChannelRequest

	// Init channels too
	client.SendMessage = make(chan *protocol.JSONRPCPayload)

	go func() {
		for client.Alive {
			select {
			case jsonRPCPayload := <-client.SendMessage:
				client.sendMessage(jsonRPCPayload)
				break
			}
		}
	}()
}

// RequestHandshake request handshake to the client
func (client *Client) RequestHandshake() {
	log.Debugf("Request handshake to the client %s", client.UID)
	client.SendMessage <- protocol.NewHandshakeRequest("", client.UID, &protobuf.ServerInformations{
		Name:            common.ApplicationName,
		ProtocolVersion: common.ApplicationProtocolVersion,
	})
}

// handleHandshakeResponse handle the handshake response of the client
func (client *Client) handleHandshakeResponse(messageID string, bytes []byte) {
	payload := &protobuf.HandshakeResponse{}
	proto.Unmarshal(bytes, payload)

	// Set the client driver
	client.ClientDriver = payload.ClientDriver

	//TDOD: Check the token given by the client for recover a state

	log.Debugf("Received Handshake response from the client, Driver: %s(%d)", payload.ClientDriver.Name, payload.ClientDriver.ProtocolVersion)
}

// handleSubscribeToChannelRequest handle the request for joining a channel by the client
func (client *Client) handleSubscribeToChannelRequest(messageID string, bytes []byte) {
	payload := &protobuf.SubscribeToChannelRequest{}
	proto.Unmarshal(bytes, payload)

	// Register client to the channel
	channel := GetOrNewChannel(payload.Channel, payload.Password)
	if _, alreadyInChannel := client.SubscribedChannels[channel.Name]; alreadyInChannel {
		// The client is already registered on this channel
		client.SendMessage <- protocol.NewSubscribeToChannelResponse(messageID, -1, "Already subscribed to this channel")
		log.Debugf("The client %s is already registered on the channel %s", client.UID, channel.Name)
		return
	}

	// Check if the password is correct
	if channel.Password != payload.Password {
		client.SendMessage <- protocol.NewSubscribeToChannelResponse(messageID, -2, "Wrong channel password")
		log.Debugf("The client %s given a wrong password for the channel %s", client.UID, channel.Name)
		return
	}

	sub := channel.RegisterClient(client, payload.RoutingKeys)
	client.SubscribedChannels[channel.Name] = sub

	// Respond to the client that he is registered to channel
	client.SendMessage <- protocol.NewSubscribeToChannelResponse(messageID, 1, "Subscribed to the channel with success")
}
