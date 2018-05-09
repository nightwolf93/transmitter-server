package protocol

import (
	"github.com/golang/protobuf/proto"
	"github.com/nightwolf93/transmitter-server/net/protobuf"
)

const (
	HandshakeRequest           = 0x01
	HandshakeResponse          = 0x02
	SubscribeToChannelRequest  = 0x03
	SubscribeToChannelResponse = 0x04
	CustomDataEvent            = 0x05
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

// NewSubscribeToChannelResponse notify the result of the channel subscribtion
func NewSubscribeToChannelResponse(messageID string, result int32, comment string) *JSONRPCPayload {
	payload := &protobuf.SubscribeToChannelResponse{
		Result:  result,
		Comment: comment,
	}
	data, err := proto.Marshal(payload)
	if err != nil {
		return nil
	}
	return &JSONRPCPayload{
		ID:     messageID,
		OpCode: SubscribeToChannelResponse,
		Data:   data,
	}
}

// NewCustomDataEvent send a custom event to the client
func NewCustomDataEvent(messageID string, eventName string, data []byte) *JSONRPCPayload {
	payload := &protobuf.CustomDataEvent{
		EventName: eventName,
		Data:      data,
	}
	data, err := proto.Marshal(payload)
	if err != nil {
		return nil
	}
	return &JSONRPCPayload{
		ID:     messageID,
		OpCode: CustomDataEvent,
		Data:   data,
	}
}
