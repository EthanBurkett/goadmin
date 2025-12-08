package models

import (
	"sync"
	"time"
)

// CommandThrottler prevents command abuse by tracking command targets
type CommandThrottler struct {
	// Map of "adminGUID:targetGUID:commandType" -> last execution time
	executions map[string]time.Time
	mu         sync.RWMutex
}

// ThrottleResult indicates whether a command should be throttled
type ThrottleResult struct {
	Allowed       bool
	Reason        string
	TimeRemaining time.Duration
}

// NewCommandThrottler creates a new command throttler
func NewCommandThrottler() *CommandThrottler {
	ct := &CommandThrottler{
		executions: make(map[string]time.Time),
	}

	// Start cleanup goroutine
	go ct.cleanup()

	return ct
}

// CheckThrottle checks if a command against a target should be allowed
func (ct *CommandThrottler) CheckThrottle(adminGUID, targetGUID, commandType string, cooldown time.Duration) *ThrottleResult {
	key := adminGUID + ":" + targetGUID + ":" + commandType

	ct.mu.RLock()
	lastExec, exists := ct.executions[key]
	ct.mu.RUnlock()

	if exists {
		timeSince := time.Since(lastExec)
		if timeSince < cooldown {
			return &ThrottleResult{
				Allowed:       false,
				Reason:        "You are targeting this player too frequently. Please wait before trying again.",
				TimeRemaining: cooldown - timeSince,
			}
		}
	}

	// Record this execution
	ct.mu.Lock()
	ct.executions[key] = time.Now()
	ct.mu.Unlock()

	return &ThrottleResult{
		Allowed: true,
	}
}

// cleanup periodically removes old entries to prevent memory leaks
func (ct *CommandThrottler) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ct.mu.Lock()
		cutoff := time.Now().Add(-1 * time.Hour) // Remove entries older than 1 hour
		for key, timestamp := range ct.executions {
			if timestamp.Before(cutoff) {
				delete(ct.executions, key)
			}
		}
		ct.mu.Unlock()
	}
}

// GetTargetStats returns statistics about how often an admin targets a player
func (ct *CommandThrottler) GetTargetStats(adminGUID, targetGUID string) int {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	count := 0
	cutoff := time.Now().Add(-15 * time.Minute)

	for key, timestamp := range ct.executions {
		// Check if key starts with adminGUID:targetGUID
		if len(key) > len(adminGUID)+len(targetGUID)+1 {
			if key[:len(adminGUID)] == adminGUID &&
				key[len(adminGUID)+1:len(adminGUID)+1+len(targetGUID)] == targetGUID &&
				timestamp.After(cutoff) {
				count++
			}
		}
	}

	return count
}

// Global instance
var CommandThrottlerInstance = NewCommandThrottler()
