package jsonrpc

import (
	"encoding/json"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req := NewRequest(1, "testMethod", map[string]any{"key": "value"})

	if req.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc version '2.0', got '%s'", req.JSONRPC)
	}
	if req.ID != 1 {
		t.Errorf("expected id 1, got %d", req.ID)
	}
	if req.Method != "testMethod" {
		t.Errorf("expected method 'testMethod', got '%s'", req.Method)
	}
	if req.Params == nil {
		t.Error("expected params to be set")
	}
}

func TestRequest_MarshalJSON(t *testing.T) {
	req := NewRequest(42, "subtract", map[string]any{"minuend": 42, "subtrahend": 23})

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal request: %v", err)
	}

	if unmarshaled["jsonrpc"] != "2.0" {
		t.Errorf("expected jsonrpc '2.0', got %v", unmarshaled["jsonrpc"])
	}
	id, ok := unmarshaled["id"].(float64)
	if !ok {
		t.Fatalf("expected id to be a number")
	}
	if id != 42 {
		t.Errorf("expected id 42, got %v", unmarshaled["id"])
	}
	if unmarshaled["method"] != "subtract" {
		t.Errorf("expected method 'subtract', got %v", unmarshaled["method"])
	}
}

func TestRequest_MarshalJSON_WithoutParams(t *testing.T) {
	req := NewRequest(1, "ping", nil)

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal request: %v", err)
	}

	// Params should be omitted when nil
	if _, ok := unmarshaled["params"]; ok {
		t.Error("expected params to be omitted when nil")
	}
}

func TestResponse_MarshalJSON_WithResult(t *testing.T) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      1,
		Result:  json.RawMessage(`{"value": 19}`),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if unmarshaled["jsonrpc"] != "2.0" {
		t.Errorf("expected jsonrpc '2.0', got %v", unmarshaled["jsonrpc"])
	}
	if _, ok := unmarshaled["result"]; !ok {
		t.Error("expected result to be present")
	}
	if _, ok := unmarshaled["error"]; ok {
		t.Error("expected error to be omitted when result is present")
	}
}

func TestResponse_MarshalJSON_WithError(t *testing.T) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      1,
		Error: &RPCError{
			Code:    ErrCodeMethodNotFound,
			Message: "method not found",
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if _, ok := unmarshaled["error"]; !ok {
		t.Error("expected error to be present")
	}
	if _, ok := unmarshaled["result"]; ok {
		t.Error("expected result to be omitted when error is present")
	}
}

func TestRPCError_Error(t *testing.T) {
	err := &RPCError{
		Code:    ErrCodeInvalidParams,
		Message: "invalid parameters",
	}

	if err.Error() != "invalid parameters" {
		t.Errorf("expected error message 'invalid parameters', got '%s'", err.Error())
	}
}

func TestRPCError_WithData(t *testing.T) {
	err := &RPCError{
		Code:    ErrCodeParseError,
		Message: "parse error",
		Data:    map[string]any{"line": 42, "column": 10},
	}

	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("failed to marshal error: %v", marshalErr)
	}

	var unmarshaled map[string]any
	if unmarshalErr := json.Unmarshal(data, &unmarshaled); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal error: %v", unmarshalErr)
	}

	if unmarshaled["data"] == nil {
		t.Error("expected data to be present")
	}
}

func TestNotification_MarshalJSON(t *testing.T) {
	notif := Notification{
		JSONRPC: "2.0",
		Method:  "update",
		Params:  map[string]any{"status": "running"},
	}

	data, err := json.Marshal(notif)
	if err != nil {
		t.Fatalf("failed to marshal notification: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal notification: %v", err)
	}

	if unmarshaled["method"] != "update" {
		t.Errorf("expected method 'update', got %v", unmarshaled["method"])
	}
	// Notification should not have an id field
	if _, ok := unmarshaled["id"]; ok {
		t.Error("expected id to be omitted in notification")
	}
}

func TestStreamEvent_MarshalJSON(t *testing.T) {
	event := StreamEvent{
		Type: "text",
		Data: json.RawMessage(`{"content": "hello"}`),
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal stream event: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal stream event: %v", err)
	}

	if unmarshaled["type"] != "text" {
		t.Errorf("expected type 'text', got %v", unmarshaled["type"])
	}
	if unmarshaled["data"] == nil {
		t.Error("expected data to be present")
	}
}

func TestStreamEvent_MarshalJSON_WithoutData(t *testing.T) {
	event := StreamEvent{
		Type: "complete",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal stream event: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal stream event: %v", err)
	}

	// Data should be omitted when nil
	if _, ok := unmarshaled["data"]; ok {
		t.Error("expected data to be omitted when nil")
	}
}

func TestStandardErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected int
	}{
		{"ParseError", ErrCodeParseError, -32700},
		{"InvalidRequest", ErrCodeInvalidRequest, -32600},
		{"MethodNotFound", ErrCodeMethodNotFound, -32601},
		{"InvalidParams", ErrCodeInvalidParams, -32602},
		{"InternalError", ErrCodeInternalError, -32603},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, tt.code)
			}
		})
	}
}

func TestStreamEventTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"Text", StreamEventText},
		{"Complete", StreamEventComplete},
		{"Error", StreamEventError},
		{"Progress", StreamEventProgress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("expected non-empty constant for %s", tt.name)
			}
		})
	}
}

func TestResponse_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"jsonrpc":"2.0","id":1,"result":19}`

	var resp Response
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc '2.0', got '%s'", resp.JSONRPC)
	}
	if resp.ID != 1 {
		t.Errorf("expected id 1, got %d", resp.ID)
	}
	if resp.Error != nil {
		t.Error("expected no error")
	}
	if resp.Result == nil {
		t.Error("expected result to be present")
	}
}

func TestResponse_UnmarshalJSON_WithError(t *testing.T) {
	jsonStr := `{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"Method not found"}}`

	var resp Response
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error to be present")
	}
	if resp.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("expected error code %d, got %d", ErrCodeMethodNotFound, resp.Error.Code)
	}
	if resp.Result != nil {
		t.Error("expected result to be nil when error is present")
	}
}
