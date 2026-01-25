// Package jsonrpc provides JSON-RPC 2.0 protocol types for plugin communication.
//
// This package implements the JSON-RPC 2.0 specification as defined in
// https://www.jsonrpc.org/specification
//
// Usage:
//
//	req := jsonrpc.NewRequest(1, "subtract", map[string]any{"minuend": 42, "subtrahend": 23})
//	resp := jsonrpc.Response{ID: 1, Result: json.RawMessage(`19`), JSONRPC: "2.0"}
//	errResp := jsonrpc.Response{ID: 1, Error: &jsonrpc.RPCError{Code: -32601, Message: "Method not found"}, JSONRPC: "2.0"}
package jsonrpc

import (
	"encoding/json"
)

// Request represents a JSON-RPC 2.0 request.
type Request struct {
	Params  any    `json:"params,omitempty"`
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int64  `json:"id"`
}

// NewRequest creates a new JSON-RPC request.
func NewRequest(id int64, method string, params any) *Request {
	return &Request{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

// Response represents a JSON-RPC 2.0 response.
type Response struct {
	Error   *RPCError       `json:"error,omitempty"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	ID      int64           `json:"id"`
}

// RPCError represents a JSON-RPC 2.0 error.
type RPCError struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *RPCError) Error() string {
	return e.Message
}

// Standard JSON-RPC error codes defined in the specification.
const (
	// ErrCodeParseError indicates invalid JSON was received by the server.
	// An error occurred on the server while parsing the JSON text.
	ErrCodeParseError = -32700
	// ErrCodeInvalidRequest indicates the received JSON is not a valid Request object.
	ErrCodeInvalidRequest = -32600
	// ErrCodeMethodNotFound indicates the method does not exist / is not available.
	ErrCodeMethodNotFound = -32601
	// ErrCodeInvalidParams indicates invalid method parameter(s).
	ErrCodeInvalidParams = -32602
	// ErrCodeInternalError indicates internal JSON-RPC error.
	ErrCodeInternalError = -32603
)

// Notification represents a JSON-RPC 2.0 notification.
// A notification is a Request object without an "id" member.
// A Notification signifies that the Client is not interested in the corresponding Response.
type Notification struct {
	Params  any    `json:"params,omitempty"`
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
}

// StreamEvent represents a streaming event in JSON-RPC communication.
// This is commonly used for long-running operations to send incremental updates.
type StreamEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// Common stream event types.
const (
	// StreamEventText indicates a text message event.
	StreamEventText = "text"
	// StreamEventComplete indicates the operation is complete.
	StreamEventComplete = "complete"
	// StreamEventError indicates an error occurred.
	StreamEventError = "error"
	// StreamEventProgress indicates a progress update.
	StreamEventProgress = "progress"
)
