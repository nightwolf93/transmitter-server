package gateway

import (
	"github.com/golang/protobuf/proto"
	"github.com/nightwolf93/transmitter-server/config"
	"github.com/nightwolf93/transmitter-server/net"
	"github.com/nightwolf93/transmitter-server/net/protobuf"
	"github.com/rs/xid"
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

// JSONRPCPayload is a rpc call struct by the client
type JSONRPCPayload struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Data   []byte `json:"data"`
}

// InitGateway initialize the gateway for client
func InitGateway() {
	cfg := config.GetConfig()
	go InitWSGatewayServer(cfg.Websocket.Host, cfg.Websocket.Port)

	// Hub for handle channels of the gateway
	go func() {
		for {
			select {
			case client := <-AddClient:
				addClient(client)
				break

			case client := <-RemoveClient:
				delete(Clients, client.ID)
				break

			case data := <-ReceiveClientData:
				log.Debugf("Received data from client (len: %d)", len(data.bytes))
				receiveDataFromClient(data.client, data.bytes)
				break
			}
		}
	}()
}

func generateClientUniqueID() string {
	return xid.New().String()
}

func addClient(client *net.Client) {
	id := generateClientUniqueID()
	client.ID = id
	Clients[id] = client
	log.Infof("Client %s registered on the server", id)
}

func receiveDataFromClient(client *net.Client, bytes []byte) {
	jsonPayload := &JSONRPCPayload{}
	json.Unmarshal(bytes, &jsonPayload)

	switch jsonPayload.Method {
	case "handshake_response":
		handshakeResponsePayload := &protobuf.HandshakeResponse{}
		proto.Unmarshal(jsonPayload.Data, handshakeResponsePayload)
		handleHandshakeResponse(client, handshakeResponsePayload)
		break
	}
}
