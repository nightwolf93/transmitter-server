package gateway

import (
	"github.com/nightwolf93/transmitter-server/config"
	"github.com/nightwolf93/transmitter-server/net"
	"github.com/nightwolf93/transmitter-server/net/protocol"
	log "github.com/sirupsen/logrus"

	"encoding/json"
)

// Clients store all client connected to the gateway
var Clients map[string]*net.Client = make(map[string]*net.Client)

// AddClient add a client throught a channel
var AddClient chan *net.Client = make(chan *net.Client)

// RemoveClient remove a client throught a channel
var RemoveClient chan *net.Client = make(chan *net.Client)

// ReceiveClientData the channel who receive all the client data stream
var ReceiveClientData chan struct {
	client *net.Client
	bytes  []byte
} = make(chan struct {
	client *net.Client
	bytes  []byte
})

// InitGateway initialize the gateway for client
func InitGateway() {
	cfg := config.GetConfig()
	go InitWSGatewayServer(cfg.Websocket.Host, cfg.Websocket.Port, cfg.Websocket.ReadBufferSize, cfg.Websocket.WriteBufferSize)

	// Hub for handle channels of the gateway
	go func() {
		for {
			select {
			case client := <-AddClient:
				addClient(client)
				break

			case client := <-RemoveClient:
				removeClient(client)
				break

			case data := <-ReceiveClientData:
				log.Debugf("Received data from client (len: %d)", len(data.bytes))
				receiveDataFromClient(data.client, data.bytes)
				break
			}
		}
	}()
}

func addClient(client *net.Client) {
	Clients[client.UID] = client
	log.Infof("Client %s registered on the server", client.UID)

	// Do the handshake with the client
	client.RequestHandshake()
}

func removeClient(client *net.Client) {
	client.Alive = false
	client.UnRegisterFromAllChannels()

	// Delete him from the client list
	delete(Clients, client.UID)
}

func receiveDataFromClient(client *net.Client, bytes []byte) {
	jsonPayload := &protocol.JSONRPCPayload{}
	json.Unmarshal(bytes, &jsonPayload)

	// Check if the method is handled, if yes call it
	if method, found := client.Handlers[jsonPayload.OpCode]; found {
		method(jsonPayload.ID, jsonPayload.Data)
	} else {
		log.Debugf("Can't find the handler for the method '%d'", jsonPayload.OpCode)
		return
	}
}
