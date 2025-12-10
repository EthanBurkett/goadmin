package webhook

import (
	"context"
	"sync"
	"time"

	"github.com/ethanburkett/goadmin/app/logger"
	"go.uber.org/zap"
)

// EventFilter is a function that determines if an event should be processed
type EventFilter func(event string, payload map[string]interface{}) bool

// EventTransformer is a function that can modify event payload before dispatch
type EventTransformer func(event string, payload map[string]interface{}) map[string]interface{}

// EventMiddleware represents a middleware that can filter or transform events
type EventMiddleware struct {
	Name        string
	Filter      EventFilter
	Transformer EventTransformer
	Priority    int // Lower numbers run first
}

// MiddlewareManager manages event middleware
type MiddlewareManager struct {
	middlewares []EventMiddleware
	mu          sync.RWMutex
}

// NewMiddlewareManager creates a new middleware manager
func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{
		middlewares: make([]EventMiddleware, 0),
	}
}

// AddMiddleware adds a new middleware to the chain
func (mm *MiddlewareManager) AddMiddleware(middleware EventMiddleware) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.middlewares = append(mm.middlewares, middleware)

	// Sort by priority (lower numbers first)
	for i := 0; i < len(mm.middlewares); i++ {
		for j := i + 1; j < len(mm.middlewares); j++ {
			if mm.middlewares[j].Priority < mm.middlewares[i].Priority {
				mm.middlewares[i], mm.middlewares[j] = mm.middlewares[j], mm.middlewares[i]
			}
		}
	}

	logger.Info("Event middleware added",
		zap.String("name", middleware.Name),
		zap.Int("priority", middleware.Priority))
}

// RemoveMiddleware removes a middleware by name
func (mm *MiddlewareManager) RemoveMiddleware(name string) bool {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	for i, m := range mm.middlewares {
		if m.Name == name {
			mm.middlewares = append(mm.middlewares[:i], mm.middlewares[i+1:]...)
			logger.Info("Event middleware removed", zap.String("name", name))
			return true
		}
	}
	return false
}

// ProcessEvent runs the event through all middleware
// Returns (shouldDispatch, transformedPayload)
func (mm *MiddlewareManager) ProcessEvent(ctx context.Context, event string, payload map[string]interface{}) (bool, map[string]interface{}) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	currentPayload := payload

	for _, middleware := range mm.middlewares {
		// Check context cancellation
		select {
		case <-ctx.Done():
			logger.Warn("Event processing cancelled", zap.String("event", event))
			return false, currentPayload
		default:
		}

		// Apply filter if present
		if middleware.Filter != nil {
			if !middleware.Filter(event, currentPayload) {
				logger.Debug("Event filtered out by middleware",
					zap.String("event", event),
					zap.String("middleware", middleware.Name))
				return false, currentPayload
			}
		}

		// Apply transformer if present
		if middleware.Transformer != nil {
			transformed := middleware.Transformer(event, currentPayload)
			if transformed != nil {
				currentPayload = transformed
			}
		}
	}

	return true, currentPayload
}

// GetMiddlewares returns a copy of all registered middlewares
func (mm *MiddlewareManager) GetMiddlewares() []EventMiddleware {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	result := make([]EventMiddleware, len(mm.middlewares))
	copy(result, mm.middlewares)
	return result
}

// Common filter functions

// FilterByEventType creates a filter that only allows specific event types
func FilterByEventType(allowedEvents ...string) EventFilter {
	allowed := make(map[string]bool)
	for _, event := range allowedEvents {
		allowed[event] = true
	}

	return func(event string, payload map[string]interface{}) bool {
		return allowed[event]
	}
}

// FilterByPayloadField creates a filter that checks a specific payload field value
func FilterByPayloadField(field string, expectedValue interface{}) EventFilter {
	return func(event string, payload map[string]interface{}) bool {
		if value, exists := payload[field]; exists {
			return value == expectedValue
		}
		return false
	}
}

// FilterByPayloadExists creates a filter that requires specific fields to exist
func FilterByPayloadExists(fields ...string) EventFilter {
	return func(event string, payload map[string]interface{}) bool {
		for _, field := range fields {
			if _, exists := payload[field]; !exists {
				return false
			}
		}
		return true
	}
}

// Common transformer functions

// AddTimestamp adds a timestamp field to the payload
func AddTimestamp() EventTransformer {
	return func(event string, payload map[string]interface{}) map[string]interface{} {
		if payload == nil {
			payload = make(map[string]interface{})
		}
		payload["middleware_timestamp"] = time.Now()
		return payload
	}
}

// AddEventType adds the event type to the payload
func AddEventType() EventTransformer {
	return func(event string, payload map[string]interface{}) map[string]interface{} {
		if payload == nil {
			payload = make(map[string]interface{})
		}
		payload["event_type"] = event
		return payload
	}
}

// RedactSensitiveFields removes sensitive data from the payload
func RedactSensitiveFields(fields ...string) EventTransformer {
	return func(event string, payload map[string]interface{}) map[string]interface{} {
		if payload == nil {
			return payload
		}

		// Create a copy to avoid modifying original
		result := make(map[string]interface{})
		for k, v := range payload {
			result[k] = v
		}

		// Redact specified fields
		for _, field := range fields {
			if _, exists := result[field]; exists {
				result[field] = "[REDACTED]"
			}
		}

		return result
	}
}

// EnrichPayload adds additional fields to the payload
func EnrichPayload(additionalFields map[string]interface{}) EventTransformer {
	return func(event string, payload map[string]interface{}) map[string]interface{} {
		if payload == nil {
			payload = make(map[string]interface{})
		}

		// Add additional fields
		for k, v := range additionalFields {
			// Don't overwrite existing fields
			if _, exists := payload[k]; !exists {
				payload[k] = v
			}
		}

		return payload
	}
}
