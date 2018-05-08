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
	client.Handlers = make(map[int]func([]byte))
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
	log.Debugf("Request handshake to the client %s", client.ID)
	client.SendMessage <- protocol.NewHandshakeRequest("", client.ID, &protobuf.ServerInformations{
		Name:            common.ApplicationName,
		ProtocolVersion: common.ApplicationProtocolVersion,
	})
}

// handleHandshakeResponse handle the handshake response of the client
func (client *Client) handleHandshakeResponse(bytes []byte) {
	payload := &protobuf.HandshakeResponse{}
	proto.Unmarshal(bytes, payload)

	// Set the client driver
	client.ClientDriver = payload.ClientDriver

	//TDOD: Check the token given by the client for recover a state

	log.Debugf("Received Handshake response from the client, Driver: %s(%d)", payload.ClientDriver.Name, payload.ClientDriver.ProtocolVersion)
}

// handleSubscribeToChannelRequest handle the request for joining a channel by the client
func (client *Client) handleSubscribeToChannelRequest(bytes []byte) {
	payload := &protobuf.SubscribeToChannelRequest{}
	proto.Unmarshal(bytes, payload)

	// Register client to the channel
	channel := GetOrNewChannel(payload.Name)
	channel.RegisterClient(client)

	// Respond to the client that he is registered to channel

}
