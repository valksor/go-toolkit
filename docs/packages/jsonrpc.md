# jsonrpc

JSON-RPC 2.0 protocol types for plugin communication.

## Overview

The `jsonrpc` package provides types that implement the JSON-RPC 2.0 specification as defined in [https://www.jsonrpc.org/specification](https://www.jsonrpc.org/specification).

It includes:

- Request and Response types
- Error types with standard error codes
- Notification type (requests without responses)
- StreamEvent for long-running operations

## Installation

```bash
go get github.com/valksor/go-toolkit/jsonrpc
```

## Usage

### Creating a Request

```go
req := jsonrpc.NewRequest(1, "subtract", map[string]any{
    "minuend":    42,
    "subtrahend": 23,
})

// Marshal to JSON
data, _ := json.Marshal(req)
// Output: {"jsonrpc":"2.0","id":1,"method":"subtract","params":{"minuend":42,"subtrahend":23}}
```

### Creating a Response

```go
// Success response
resp := jsonrpc.Response{
    JSONRPC: "2.0",
    ID:      1,
    Result:  json.RawMessage(`19`),
}

// Error response
errResp := jsonrpc.Response{
    JSONRPC: "2.0",
    ID:      1,
    Error: &jsonrpc.RPCError{
        Code:    jsonrpc.ErrCodeMethodNotFound,
        Message: "Method not found",
    },
}
```

### Notifications

Notifications are requests without an ID, used when no response is expected:

```go
notif := jsonrpc.Notification{
    JSONRPC: "2.0",
    Method:  "update",
    Params:  map[string]any{"status": "running"},
}
```

### Standard Error Codes

| Code | Constant | Description |
|------|----------|-------------|
| -32700 | `ErrCodeParseError` | Invalid JSON was received |
| -32600 | `ErrCodeInvalidRequest` | Received JSON is not a valid Request |
| -32601 | `ErrCodeMethodNotFound` | Method does not exist |
| -32602 | `ErrCodeInvalidParams` | Invalid method parameters |
| -32603 | `ErrCodeInternalError` | Internal JSON-RPC error |

Custom error codes in the range -32000 to -32099 are reserved for application-specific errors.

### Stream Events

For long-running operations, use `StreamEvent` to send incremental updates:

```go
event := jsonrpc.StreamEvent{
    Type: jsonrpc.StreamEventProgress,
    Data: json.RawMessage(`{"percent": 50}`),
}
```

Available stream event types:
- `StreamEventText` - Text message
- `StreamEventComplete` - Operation complete
- `StreamEventError` - Error occurred
- `StreamEventProgress` - Progress update

## Example: JSON-RPC Server

```go
func handleRequest(conn net.Conn) {
    decoder := json.NewDecoder(conn)
    encoder := json.NewEncoder(conn)

    var req jsonrpc.Request
    if err := decoder.Decode(&req); err != nil {
        sendError(encoder, 0, jsonrpc.ErrCodeParseError, "Parse error")
        return
    }

    // Process request...
    result := processMethod(req.Method, req.Params)

    resp := jsonrpc.Response{
        JSONRPC: "2.0",
        ID:      req.ID,
        Result:  result,
    }
    encoder.Encode(resp)
}
```

## Example: JSON-RPC Client

```go
func callMethod(method string, params any) (json.RawMessage, error) {
    req := jsonrpc.NewRequest(nextID(), method, params)
    data, _ := json.Marshal(req)

    resp, err := httpClient.Post(url, "application/json", bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var rpcResp jsonrpc.Response
    if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
        return nil, err
    }

    if rpcResp.Error != nil {
        return nil, rpcResp.Error
    }

    return rpcResp.Result, nil
}
```

## Specification Compliance

This package follows the JSON-RPC 2.0 specification:

- Request objects must have `jsonrpc`, `method`, and `id` members
- Response objects must have `jsonrpc` and `id` members
- Either `result` or `error` must be present, but not both
- Notifications omit the `id` member
