package plugins

import (
	"fmt"
	"sync"
)

// EventBus provides pub/sub event handling for plugins
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]EventCallback
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventCallback),
	}
}

// Subscribe subscribes to an event
func (eb *EventBus) Subscribe(eventType string, callback EventCallback) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], callback)
	fmt.Printf("[EventBus] Subscribed to event: %s (total subscribers: %d)\n", eventType, len(eb.subscribers[eventType]))
	return nil
}

// Unsubscribe unsubscribes from an event
func (eb *EventBus) Unsubscribe(eventType string, callback EventCallback) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// Note: Unsubscribe is limited in Go due to function pointer comparison
	// In a production system, you'd want to use a handle/ID-based approach
	return nil
}

// Publish publishes an event to all subscribers
func (eb *EventBus) Publish(eventType string, data map[string]interface{}) error {
	eb.mu.RLock()
	handlers := make([]EventCallback, len(eb.subscribers[eventType]))
	copy(handlers, eb.subscribers[eventType])
	subscriberCount := len(handlers)
	eb.mu.RUnlock()

	fmt.Printf("[EventBus] Publishing event: %s (subscribers: %d)\n", eventType, subscriberCount)

	// Call handlers in goroutines to prevent blocking
	for _, handler := range handlers {
		go handler(eventType, data)
	}

	return nil
}

// GlobalEventBus is the global event bus instance
var GlobalEventBus = NewEventBus()
