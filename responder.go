package jsonrpc

import (
	"encoding/json"
	"sync"
)

// Responder build Response message
type Responder struct {
	request Request
	proxy   Proxy
}

// Respond build Response message with associated Request ID if provided
func (responder *Responder) Respond(result interface{}, err *Error) {
	if err != nil {
		result = nil
	}

	responder.proxy.Encode(Response{
		ID:      responder.request.ID,
		Result:  &Params{Value: result},
		Error:   err,
		Version: SupportedVersion,
	})
}

// SingleProxy encode one Response with JSON encoder
type SingleProxy struct {
	encoder *json.Encoder
}

// Encode implements Proxy interface
func (proxy SingleProxy) Encode(response Response) {
	proxy.encoder.Encode(response)
}

// NewBatchProxy creates new BatchProxy
func NewBatchProxy(encoder *json.Encoder) *BatchProxy {
	return &BatchProxy{
		guard:     new(sync.Mutex),
		processes: new(sync.WaitGroup),
		encoder:   encoder,
	}
}

// BatchProxy encode all Response objects with JSON encoder
type BatchProxy struct {
	guard     *sync.Mutex
	processes *sync.WaitGroup
	responses []Response
	encoder   *json.Encoder
}

// Wait will wait all Response object to process
func (proxy *BatchProxy) Wait() {
	proxy.processes.Wait()
	if len(proxy.responses) > 0 {
		proxy.encoder.Encode(proxy.responses)
	}
}

// Encode implements Proxy interface
func (proxy *BatchProxy) Encode(response Response) {
	proxy.guard.Lock()
	defer proxy.guard.Unlock()

	proxy.responses = append(proxy.responses, response)
	proxy.processes.Done()
}
