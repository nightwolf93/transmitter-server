package net

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/nightwolf93/transmitter-server/common"
	"github.com/nightwolf93/transmitter-server/net/protobuf"
	"github.com/nightwolf93/transmitter-server/net/protocol"
	log "github.com/sirupsen/logrus"
)

const (
	// ClientGatewayTypeWebsocket websocket transport type
	ClientGatewayTypeWebsocket = 1
)

// Client is the struct that handle a client
type Client struct {
	UID           string
	TransportType int
	Alive         bool
	Handlers      map[int]func(string, []byte)
	ClientDriver  *protobuf.ClientDriver
	Peer          *protobuf.PeerItem

	SendMessage chan *protocol.JSONRPCPayload

	// Channel
	SubscribedChannels map[string]*ChannelSubscriber

	// Specific connection
	WSConn               *websocket.Conn
	WSReceiveMessageChan chan []byte
}

// NewClient create a new client with the given transport type
func NewClient(transportType int) *Client {
	client := &Client{
		Alive:              true,
		TransportType:      transportType,
		UID:                common.GenerateLongUniqueID(),
		SubscribedChannels: make(map[string]*ChannelSubscriber),
	}
	client.Peer = &protobuf.PeerItem{
		Uid: client.UID,
	}
	client.InitHandlers()
	return client
}

// SendMessage send a message to the client, set the id from the request message for respond to the client
func (client *Client) sendMessage(payload *protocol.JSONRPCPayload) {
	// Send message with transport of the client
	switch client.TransportType {
	case ClientGatewayTypeWebsocket:
		client.sendWebsocketMessage(payload)
		break
	}
}

func (client *Client) sendWebsocketMessage(payload *protocol.JSONRPCPayload) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Can't serialize message to send to the client: %s", err)
		return
	}
	err = client.WSConn.WriteMessage(1, jsonBytes)
	if err != nil {
		log.Errorf("Can't send message to the client: %s", err)
		return
	}
	log.Debugf("Send a message to the client using a websocket transport (len: %d)", len(jsonBytes))
}

// UnRegisterFromAllChannels unregister the client from all subscribed channels
func (client *Client) UnRegisterFromAllChannels() {
	for _, sub := range client.SubscribedChannels {
		sub.Channel.UnRegisterClient(client)
	}
}
