package jsonrpc

import (
	"encoding/json"
)

// Message incoming message
type Message struct {
	ID      json.RawMessage `json:"id"`
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	Result  json.RawMessage `json:"result"`
	Error   *Error          `json:"error"`
}

// IsRequest checks that incoming message is JSONRPC request
func (message Message) IsRequest() bool {
	if message.Result == nil && message.Error == nil {
		return true
	}

	return false
}

// ToRequest build JSONRPC request object from Message struct
func (message Message) ToRequest() Request {
	return Request{
		ID:      message.ID,
		Method:  message.Method,
		Params:  &Params{Data: message.Params},
		Version: message.Version,
	}
}

// ToResponse build JSONRPC response object from Message struct
func (message Message) ToResponse() Response {
	return Response{
		ID:      message.ID,
		Result:  &Params{Data: message.Result},
		Error:   message.Error,
		Version: message.Version,
	}
}
