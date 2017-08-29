package jsonrpc

import (
	"sync"
)

// NewDispatcher creates default dispatcher
func NewDispatcher() *DefaultDispatcher {
	return &DefaultDispatcher{
		callbacks: make(map[string]Callback),
		handlers:  make(map[string]Handler),
		methods:   make(map[string]Method),
	}
}

// DefaultDispatcher provides default dispatcher implementation
type DefaultDispatcher struct {
	guard     sync.RWMutex
	callbacks map[string]Callback
	handlers  map[string]Handler
	methods   map[string]Method
}

// RegisterMethod register new RPC method
func (dispatcher *DefaultDispatcher) RegisterMethod(name string, method Method) {
	dispatcher.guard.Lock()
	defer dispatcher.guard.Unlock()

	dispatcher.methods[name] = method
}

// RemoveMethod removes RPC method
func (dispatcher *DefaultDispatcher) RemoveMethod(name string) {
	dispatcher.guard.Lock()
	defer dispatcher.guard.Unlock()

	delete(dispatcher.methods, name)
}

// RegisterCallback register new RPC method
func (dispatcher *DefaultDispatcher) RegisterCallback(id string, callback Callback) {
	dispatcher.guard.Lock()
	defer dispatcher.guard.Unlock()

	dispatcher.callbacks[id] = callback
}

// RegisterHandler register new RPC method
func (dispatcher *DefaultDispatcher) RegisterHandler(name string, handler Handler) {
	dispatcher.guard.Lock()
	defer dispatcher.guard.Unlock()

	dispatcher.handlers[name] = handler
}

// RemoveHandler removes RPC method
func (dispatcher *DefaultDispatcher) RemoveHandler(name string) {
	dispatcher.guard.Lock()
	defer dispatcher.guard.Unlock()

	delete(dispatcher.handlers, name)
}

// DispatchNotification provided request params to RPC method
func (dispatcher *DefaultDispatcher) DispatchNotification(request Request) {
	dispatcher.guard.RLock()
	defer dispatcher.guard.RUnlock()

	handler, exists := dispatcher.handlers[request.Method]
	if exists {
		handler.Process(request.Params.Data)
	}
}

// DispatchResponse provided request params to RPC method
func (dispatcher *DefaultDispatcher) DispatchResponse(response Response) {
	dispatcher.guard.RLock()
	defer dispatcher.guard.RUnlock()

	id := string(response.ID)

	callback, exists := dispatcher.callbacks[id]
	if exists {
		callback.Process(response)
		delete(dispatcher.callbacks, id)
	}
}

// DispatchRequest provided request params to RPC method
func (dispatcher *DefaultDispatcher) DispatchRequest(request Request, responder Responder) {
	dispatcher.guard.RLock()
	defer dispatcher.guard.RUnlock()

	method, exists := dispatcher.methods[request.Method]
	if exists {
		responder.Respond(method.Process(request.Params.Data))
	}

	responder.Respond(nil, NewError(MethodNotFound, "Method ", request.Method, " Not Found"))
}
