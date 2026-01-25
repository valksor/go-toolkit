# eventbus

Generic publish/subscribe event system for Go.

## Overview

The `eventbus` package provides a clean, thread-safe event bus for decoupled communication between components. It supports:

- Type-based event routing
- Synchronous and asynchronous publishing
- Wildcard subscriptions (subscribe to all events)
- Semaphore-based limiting for async operations
- Graceful shutdown

## Installation

```bash
go get github.com/valksor/go-toolkit/eventbus
```

## Usage

### Basic Publishing

```go
bus := eventbus.NewBus()
defer bus.Shutdown()

// Subscribe to a specific event type
bus.Subscribe(eventbus.Type("user_created"), func(e eventbus.Event) {
    fmt.Printf("User created: %+v\n", e.Data)
})

// Publish an event
bus.PublishRaw(eventbus.Event{
    Type: "user_created",
    Data: map[string]any{"id": 123, "name": "Alice"},
})
```

### Subscribe to All Events

```go
bus.SubscribeAll(func(e eventbus.Event) {
    fmt.Printf("Event: %s\n", e.Type)
})
```

### Asynchronous Publishing

```go
// Publish asynchronously (in a goroutine)
bus.PublishRawAsync(eventbus.Event{Type: "background_task"})

// Shutdown waits for all async publishes to complete
bus.Shutdown()
```

### Unsubscribing

```go
id := bus.Subscribe(eventbus.Type("task"), handler)
// Later...
bus.Unsubscribe(id)
```

### Typed Events

Define your own event types by implementing the `Eventer` interface:

```go
type UserCreatedEvent struct {
    UserID    int
    UserName  string
    Timestamp time.Time
}

func (e UserCreatedEvent) ToEvent() eventbus.Event {
    if e.Timestamp.IsZero() {
        e.Timestamp = time.Now()
    }
    return eventbus.Event{
        Type: "user_created",
        Data: map[string]any{
            "user_id": e.UserID,
            "name":    e.UserName,
        },
    }
}

// Publish using the typed event
bus.Publish(UserCreatedEvent{UserID: 123, UserName: "Alice"})
```

## Thread Safety

All methods are safe for concurrent use. Internal state is protected by a read-write mutex.

## Semantics

- **Synchronous Publish**: Handlers are called in the same goroutine as the caller. Lock is held only while collecting handlers, not during execution.
- **Asynchronous Publish**: Handlers are executed in goroutines limited by a semaphore (100 concurrent by default).
- **Shutdown**: Waits for all in-flight async publishes to complete before returning.

## Best Practices

1. **Keep handlers fast**: Async handlers can slow down the event bus if they block.
2. **Avoid panic in handlers**: Panics in handlers will crash the goroutine.
3. **Unsubscribe when done**: Prevents memory leaks from dangling subscriptions.
4. **Use typed events**: Improves type safety and documentation.
