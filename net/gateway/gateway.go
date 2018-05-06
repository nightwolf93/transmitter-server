package gateway

import (
	"github.com/nightwolf93/transmitter-server/config"
	"github.com/nightwolf93/transmitter-server/net"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

// Clients store all client connected to the gateway
var Clients map[string]*net.Client = make(map[string]*net.Client)

// AddClient add a client throught a channel
var AddClient chan *net.Client = make(chan *net.Client)
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
	go InitWSGatewayServer(cfg.Websocket.Host, cfg.Websocket.Port)

	// Hub for handle channels of the gateway
	go func() {
		for {
			select {
			case client := <-AddClient:
				addClient(client)
				break

			case data := <-ReceiveClientData:
				log.Debugf("Received data from client (len: %d)", len(data.bytes))
				_ = data
				//client := data.client
				//bytes := data.bytes
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
	Clients[id] = client
	log.Infof("Client %s registered on the server", id)
}
