package jsonrpc

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)

var id uint64

type Params struct {
	Data  []byte
	Value interface{}
}

// MarshalJSON returns params as the JSON encoding of params.
func (params Params) MarshalJSON() ([]byte, error) {
	if params.Value == nil {
		return []byte("null"), nil
	}

	return json.Marshal(params.Value)
}

// UnmarshalJSON sets *params to a copy of data.
func (params *Params) UnmarshalJSON(data []byte) error {
	params.Data = append(params.Data[0:0], data...)

	return nil
}

// NewRequest creates new JSONRPC request object
func NewRequest(method string, params interface{}) Request {
	return Request{
		ID:      strconv.AppendUint(nil, atomic.AddUint64(&id, 1), 10),
		Method:  method,
		Version: SupportedVersion,
		Params:  &Params{Value: params},
	}
}

// NewNotification creates new JSONRPC request object without ID
func NewNotification(method string, params interface{}) Request {
	return Request{
		Method:  method,
		Version: SupportedVersion,
		Params:  &Params{Value: params},
	}
}

// Request JSONRPC request object representation
type Request struct {
	ID      json.RawMessage `json:"id,omitempty"`
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  *Params         `json:"params,omitempty"`
}

// Process request with method dispatcher
func (request Request) Process(dispatcher Dispatcher, responder Responder) {
	if request.Version != SupportedVersion {
		responder.Respond(nil, NewError(InvalidRequest, "Version '%s' Is Not Supported", request.Version))
	}

	if request.Method == "" {
		responder.Respond(nil, NewError(InvalidRequest, "Empty Method"))
	}

	dispatcher.DispatchRequest(request, responder)
}

// IsNotification handles JSONRPC notification
func (request Request) IsNotification() bool {
	return request.ID == nil
}
