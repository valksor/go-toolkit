// Package eventbus provides a generic publish/subscribe event system for Go.
//
// The event bus enables decoupled communication between components through
// typed events with support for synchronous and asynchronous publishing.
//
// Thread safety:
//   - All methods are safe for concurrent use.
//   - Internal state is protected by a read-write mutex.
//
// Features:
//   - Type-based event routing with wildcard support
//   - Synchronous and asynchronous publishing
//   - Semaphore-based limiting for async operations
//   - Graceful shutdown with context cancellation
//
// Usage:
//
//	bus := eventbus.New()
//	bus.Subscribe(eventbus.Type("user_created"), func(e eventbus.Event) {
//	    fmt.Printf("User created: %+v\n", e.Data)
//	})
//	bus.Publish(eventbus.Event{Type: "user_created", Data: map[string]any{"id": 123}})
//	bus.Shutdown()
package eventbus

import "time"

// Type identifies event categories. Applications define their own types.
type Type string

// Event is the base event structure.
type Event struct {
	Timestamp time.Time
	Data      map[string]any
	Type      Type
}

// Eventer interface for typed events.
// Implement this interface to create strongly-typed event wrappers.
type Eventer interface {
	ToEvent() Event
}
