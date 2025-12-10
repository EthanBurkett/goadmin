package plugins

import (
	"context"
	"time"
)

// Plugin represents the interface all plugins must implement
type Plugin interface {
	// Metadata returns information about the plugin
	Metadata() PluginMetadata

	// Init initializes the plugin with the provided context
	// This is called once when the plugin is loaded
	Init(ctx *PluginContext) error

	// Start starts the plugin
	// This is called after Init and when the plugin is enabled
	Start() error

	// Stop stops the plugin
	// This is called when the plugin is disabled or server shuts down
	Stop() error

	// Reload reloads the plugin configuration
	// This is called when plugin config changes
	Reload() error
}

// PluginMetadata contains information about a plugin
type PluginMetadata struct {
	ID             string          `json:"id"`                       // Unique identifier
	Name           string          `json:"name"`                     // Display name
	Version        string          `json:"version"`                  // Semantic version
	Author         string          `json:"author"`                   // Plugin author
	Description    string          `json:"description"`              // What the plugin does
	Website        string          `json:"website"`                  // Plugin website/repo
	Dependencies   []string        `json:"dependencies"`             // Required plugin IDs
	Permissions    []string        `json:"permissions"`              // Required permissions
	MinAPIVersion  string          `json:"minApiVersion,omitempty"`  // Minimum GoAdmin API version required
	MaxAPIVersion  string          `json:"maxApiVersion,omitempty"`  // Maximum GoAdmin API version supported
	ResourceLimits *ResourceLimits `json:"resourceLimits,omitempty"` // Optional resource limits
}

// ResourceLimits defines resource constraints for a plugin
type ResourceLimits struct {
	MaxMemoryMB   int64         `json:"maxMemoryMb,omitempty"`   // Maximum memory in MB
	MaxCPUPercent float64       `json:"maxCpuPercent,omitempty"` // Maximum CPU usage percentage
	MaxGoroutines int           `json:"maxGoroutines,omitempty"` // Maximum number of goroutines
	Timeout       time.Duration `json:"timeout,omitempty"`       // Maximum execution time for operations
}

// PluginContext provides access to GoAdmin APIs
type PluginContext struct {
	// API access
	EventBus    EventBusAPI
	CommandAPI  CommandAPI
	RCONAPI     RCONAPI
	DatabaseAPI DatabaseAPI
	WebhookAPI  WebhookAPI
	ConfigAPI   ConfigAPI

	// Plugin metadata
	PluginID   string
	PluginDir  string
	ConfigPath string

	// Cancellation context for graceful shutdown
	Context    context.Context
	CancelFunc context.CancelFunc
}

// EventBusAPI provides access to the event system
type EventBusAPI interface {
	// Subscribe registers a callback for an event type
	Subscribe(eventType string, callback EventCallback) error

	// Unsubscribe removes a callback for an event type
	Unsubscribe(eventType string, callback EventCallback) error

	// Publish dispatches an event to all subscribers
	Publish(eventType string, data map[string]interface{}) error
}

// EventCallback is called when an event is triggered
type EventCallback func(eventType string, data map[string]interface{}) error

// CommandAPI provides access to command registration
type CommandAPI interface {
	// RegisterCommand registers a custom in-game command
	RegisterCommand(cmd CommandDefinition) error

	// UnregisterCommand removes a custom command
	UnregisterCommand(name string) error

	// ExecuteCommand executes a command programmatically
	ExecuteCommand(playerName, playerGUID, command string, args []string) error
}

// CommandDefinition defines a custom command
type CommandDefinition struct {
	Name            string
	Usage           string
	Description     string
	MinArgs         int
	MaxArgs         int
	MinPower        int
	Permissions     []string
	RequirementType string
	Handler         CommandHandler
}

// CommandHandler is called when a command is executed
type CommandHandler func(playerName, playerGUID string, args []string) error

// RCONAPI provides access to RCON communication
type RCONAPI interface {
	// SendCommand sends a raw RCON command
	SendCommand(command string) (string, error)

	// SendCommandWithTimeout sends a command with a custom timeout
	SendCommandWithTimeout(command string, timeout time.Duration) (string, error)

	// GetStatus gets server status
	GetStatus() (map[string]interface{}, error)
}

// DatabaseAPI provides access to database operations
type DatabaseAPI interface {
	// GetDB returns the GORM database instance
	GetDB() interface{}

	// Query executes a raw SQL query
	Query(sql string, args ...interface{}) ([]map[string]interface{}, error)

	// Exec executes a raw SQL statement
	Exec(sql string, args ...interface{}) error
}

// WebhookAPI provides access to webhook system
type WebhookAPI interface {
	// Dispatch sends a webhook for a custom event
	Dispatch(event string, data map[string]interface{}) error

	// RegisterEvent registers a custom webhook event type
	RegisterEvent(eventType string, description string) error
}

// ConfigAPI provides access to plugin configuration
type ConfigAPI interface {
	// Get retrieves a configuration value
	Get(key string) (interface{}, error)

	// Set stores a configuration value
	Set(key string, value interface{}) error

	// GetString retrieves a string configuration value
	GetString(key string, defaultValue string) string

	// GetInt retrieves an integer configuration value
	GetInt(key string, defaultValue int) int

	// GetBool retrieves a boolean configuration value
	GetBool(key string, defaultValue bool) bool
}

// PluginState represents the current state of a plugin
type PluginState string

const (
	PluginStateLoaded  PluginState = "loaded"
	PluginStateStarted PluginState = "started"
	PluginStateStopped PluginState = "stopped"
	PluginStateError   PluginState = "error"
)

// PluginStatus represents the runtime status of a plugin
type PluginStatus struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Version  string      `json:"version"`
	State    PluginState `json:"state"`
	Enabled  bool        `json:"enabled"`
	LoadedAt time.Time   `json:"loadedAt"`
	Error    string      `json:"error,omitempty"`
}
