package jsonrpc

import (
	"encoding/json"
	"fmt"
)

// Error JSONRPC error object representation
type Error struct {
	Code    int            `json:"code"`
	Message string         `json:"message,omitempty"`
	Data    json.Marshaler `json:"data,omitempty"`
}

// NewError create new error object be code and error message
func NewError(code int, params ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprint(params...),
	}
}
