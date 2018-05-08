package protocol

// JSONRPCPayload is a rpc call struct by the client
type JSONRPCPayload struct {
	ID     string `json:"id"`
	OpCode int    `json:"opcode"`
	Data   []byte `json:"data"`
}
