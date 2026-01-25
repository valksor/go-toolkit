package eventbus

import (
	"sync"
	"testing"
	"time"
)

// mockEventer implements Eventer for testing.
type mockEventer struct {
	event Event
}

func (m mockEventer) ToEvent() Event {
	return m.event
}

func TestBus_SubscribeAndPublish(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	var received Event
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe(Type("test"), func(e Event) {
		received = e
		wg.Done()
	})

	event := Event{Type: Type("test"), Data: map[string]any{"key": "value"}}
	bus.PublishRaw(event)

	wg.Wait()

	if received.Type != Type("test") {
		t.Fatalf("expected type 'test', got %v", received.Type)
	}
	if received.Data["key"] != "value" {
		t.Fatalf("expected data key 'value', got %v", received.Data["key"])
	}
}

func TestBus_SubscribeWithEventer(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	var received Event
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe(Type("test"), func(e Event) {
		received = e
		wg.Done()
	})

	eventer := mockEventer{
		event: Event{Type: Type("test"), Data: map[string]any{"key": "value"}},
	}
	bus.Publish(eventer)

	wg.Wait()

	if received.Type != Type("test") {
		t.Fatalf("expected type 'test', got %v", received.Type)
	}
}

func TestBus_SubscribeAll(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	count := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	bus.Subscribe(Type("specific"), func(e Event) {
		mu.Lock()
		count++
		mu.Unlock()
		wg.Done()
	})

	bus.SubscribeAll(func(e Event) {
		mu.Lock()
		count++
		mu.Unlock()
		wg.Done()
	})

	// Publish to specific type
	bus.PublishRaw(Event{Type: Type("specific")})
	wg.Wait()

	mu.Lock()
	if count != 2 {
		mu.Unlock()
		t.Fatalf("expected 2 handlers called, got %d", count)
	}
	mu.Unlock()
}

func TestBus_Unsubscribe(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	count := 0
	var mu sync.Mutex

	bus.Subscribe(Type("test"), func(e Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	id := bus.Subscribe(Type("test"), func(e Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	// Unsubscribe the second handler
	bus.Unsubscribe(id)

	bus.PublishRaw(Event{Type: Type("test")})

	mu.Lock()
	if count != 1 {
		mu.Unlock()
		t.Fatalf("expected 1 handler called after unsubscribe, got %d", count)
	}
	mu.Unlock()
}

func TestBus_UnsubscribeAllHandler(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	count := 0
	var mu sync.Mutex

	id := bus.SubscribeAll(func(e Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	bus.Unsubscribe(id)

	bus.PublishRaw(Event{Type: Type("any")})

	mu.Lock()
	if count != 0 {
		mu.Unlock()
		t.Fatalf("expected 0 handlers called after unsubscribe, got %d", count)
	}
	mu.Unlock()
}

func TestBus_PublishAsync(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	var received Event
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe(Type("async"), func(e Event) {
		received = e
		wg.Done()
	})

	bus.PublishRawAsync(Event{Type: Type("async"), Data: map[string]any{"key": "value"}})

	wg.Wait()

	if received.Type != Type("async") {
		t.Fatalf("expected type 'async', got %v", received.Type)
	}
}

func TestBus_HasSubscribers(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	if bus.HasSubscribers(Type("test")) {
		t.Fatal("expected no subscribers for 'test'")
	}

	bus.Subscribe(Type("test"), func(e Event) {})

	if !bus.HasSubscribers(Type("test")) {
		t.Fatal("expected subscribers for 'test'")
	}
}

func TestBus_HasSubscribersWithAll(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	if bus.HasSubscribers(Type("any")) {
		t.Fatal("expected no subscribers")
	}

	bus.SubscribeAll(func(e Event) {})

	if !bus.HasSubscribers(Type("any")) {
		t.Fatal("expected subscribers via SubscribeAll")
	}
}

func TestBus_Clear(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	bus.Subscribe(Type("test"), func(e Event) {})
	bus.SubscribeAll(func(e Event) {})

	bus.Clear()

	if bus.HasSubscribers(Type("test")) {
		t.Fatal("expected no subscribers after Clear")
	}
}

func TestBus_Shutdown(t *testing.T) {
	bus := NewBus()

	var called bool
	var mu sync.Mutex

	// Add a single handler
	bus.Subscribe(Type("test"), func(e Event) {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		called = true
		mu.Unlock()
	})

	// Publish multiple async events
	for range 10 {
		bus.PublishRawAsync(Event{Type: Type("test")})
	}

	// Shutdown should wait for all async publishes to complete
	bus.Shutdown()

	mu.Lock()
	if !called {
		mu.Unlock()
		t.Fatal("expected handlers to be called before shutdown")
	}
	mu.Unlock()
}

func TestBus_ConcurrentPublish(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	var count int
	var mu sync.Mutex
	bus.SubscribeAll(func(e Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	// Publish concurrently
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bus.PublishRaw(Event{Type: Type("test")})
		}()
	}

	wg.Wait()

	mu.Lock()
	if count != 100 {
		mu.Unlock()
		t.Fatalf("expected 100 events, got %d", count)
	}
	mu.Unlock()
}

func TestBus_ConcurrentSubscribe(t *testing.T) {
	bus := NewBus()
	defer bus.Shutdown()

	// Subscribe concurrently
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bus.Subscribe(Type("test"), func(e Event) {})
		}()
	}

	wg.Wait()

	// Verify all subscriptions were registered
	if !bus.HasSubscribers(Type("test")) {
		t.Fatal("expected subscribers after concurrent subscription")
	}
}

func TestEventer_Interface(t *testing.T) {
	eventer := mockEventer{
		event: Event{
			Type:      Type("test"),
			Data:      map[string]any{"key": "value"},
			Timestamp: time.Now(),
		},
	}

	event := eventer.ToEvent()

	if event.Type != Type("test") {
		t.Fatalf("expected type 'test', got %v", event.Type)
	}
	if event.Data["key"] != "value" {
		t.Fatalf("expected data key 'value', got %v", event.Data["key"])
	}
}
