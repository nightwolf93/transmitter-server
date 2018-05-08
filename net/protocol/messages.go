package protocol

import (
	"github.com/golang/protobuf/proto"
	"github.com/nightwolf93/transmitter-server/net/protobuf"
)

const (
	HandshakeRequest          = 0x01
	HandshakeResponse         = 0x02
	SubscribeToChannelRequest = 0x03
)

// NewHandshakeRequest build a new handshake request message
func NewHandshakeRequest(messageID string, uid string, serverInformations *protobuf.ServerInformations) *JSONRPCPayload {
	payload := &protobuf.HandshakeRequest{
		Uid:                uid,
		ServerInformations: serverInformations,
	}
	data, err := proto.Marshal(payload)
	if err != nil {
		return nil
	}
	return &JSONRPCPayload{
		ID:     messageID,
		OpCode: HandshakeRequest,
		Data:   data,
	}
}
