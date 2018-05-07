package gateway

import (
	"github.com/nightwolf93/transmitter-server/net"
	"github.com/nightwolf93/transmitter-server/net/protobuf"
	log "github.com/sirupsen/logrus"
)

// handleHandshakeResponse handle the handshake response of the client
func handleHandshakeResponse(client *net.Client, message *protobuf.HandshakeResponse) {
	log.Debugf("Received Handshake response from the client, Driver: %s(%d)", message.ClientDriver.Name, message.ClientDriver.ProtocolVersion)
}
