package jsonrpc

import (
	"encoding/json"
)

// Response JSONRPC response object representation
type Response struct {
	ID      json.RawMessage `json:"id,omitempty"`
	Version string          `json:"jsonrpc"`
	Result  *Params         `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

// NewErrorResponse build response for provided error
func NewErrorResponse(err *Error) Response {
	return Response{
		Version: SupportedVersion,
		Error:   err,
	}
}

// Process request with method dispatcher
func (response Response) Process(dispatcher Dispatcher) {
	if response.Version != SupportedVersion || response.ID == nil {
		return
	}

	dispatcher.DispatchResponse(response)
}
