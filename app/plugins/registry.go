package plugins

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethanburkett/goadmin/app/logger"
	"go.uber.org/zap"
)

// Registry holds all registered plugins
var Registry = &PluginRegistry{
	plugins: make(map[string]Plugin),
	mu:      sync.RWMutex{},
}

// PluginRegistry manages plugin registration
type PluginRegistry struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

// Register adds a plugin to the registry
// Plugins should call this in their init() function
func (r *PluginRegistry) Register(plugin Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	metadata := plugin.Metadata()
	if metadata.ID == "" {
		return fmt.Errorf("plugin ID cannot be empty")
	}

	if _, exists := r.plugins[metadata.ID]; exists {
		return fmt.Errorf("plugin %s already registered", metadata.ID)
	}

	r.plugins[metadata.ID] = plugin
	// Don't log here - logger may not be initialized yet during init()
	return nil
}

// GetAll returns all registered plugins
func (r *PluginRegistry) GetAll() map[string]Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make(map[string]Plugin, len(r.plugins))
	for id, p := range r.plugins {
		plugins[id] = p
	}
	return plugins
}

// Get returns a specific plugin by ID
func (r *PluginRegistry) Get(id string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.plugins[id]
	return p, exists
}

// Manager manages plugin lifecycle
type Manager struct {
	plugins      map[string]*LoadedPlugin
	pluginStates map[string]PluginState
	mu           sync.RWMutex
	rconAPI      RCONAPI
	commandAPI   *CommandAPIImpl
}

// LoadedPlugin represents a loaded plugin instance
type LoadedPlugin struct {
	Plugin     Plugin
	Metadata   PluginMetadata
	Context    *PluginContext
	State      PluginState
	LoadedAt   time.Time
	Error      string
	cancelFunc context.CancelFunc
}

// NewManager creates a new plugin manager
func NewManager() *Manager {
	return &Manager{
		plugins:      make(map[string]*LoadedPlugin),
		pluginStates: make(map[string]PluginState),
	}
}

// SetRCONClient sets the RCON client for the manager
func (m *Manager) SetRCONClient(rconAPI RCONAPI) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rconAPI = rconAPI
	m.commandAPI = NewCommandAPI(rconAPI)
}

// GetCommandAPI returns the command API instance
func (m *Manager) GetCommandAPI() *CommandAPIImpl {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.commandAPI
}

// LoadAll loads all registered plugins from the registry
func (m *Manager) LoadAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	registeredPlugins := Registry.GetAll()
	if len(registeredPlugins) == 0 {
		logger.Info("No plugins registered")
		return nil
	}

	logger.Info(fmt.Sprintf("Found %d registered plugin(s)", len(registeredPlugins)))

	for id, plugin := range registeredPlugins {
		if _, exists := m.plugins[id]; exists {
			continue // Already loaded
		}

		metadata := plugin.Metadata()
		logger.Info("Loading plugin", zap.String("id", id), zap.String("name", metadata.Name), zap.String("version", metadata.Version))

		if err := m.loadPlugin(id, plugin); err != nil {
			logger.Error("Failed to load plugin", zap.String("id", id), zap.Error(err))
			continue
		}
	}

	logger.Info(fmt.Sprintf("Successfully loaded %d plugin(s)", len(m.plugins)))
	return nil
}

// loadPlugin loads and initializes a single plugin
func (m *Manager) loadPlugin(id string, pluginInstance Plugin) error {
	// Get metadata
	metadata := pluginInstance.Metadata()

	// Check if plugin ID is already loaded
	if _, exists := m.plugins[metadata.ID]; exists {
		return fmt.Errorf("plugin with ID %s is already loaded", metadata.ID)
	}

	// Create plugin context
	ctx, cancel := context.WithCancel(context.Background())
	pluginCtx := &PluginContext{
		PluginID:   metadata.ID,
		Context:    ctx,
		CancelFunc: cancel,
		EventBus:   GlobalEventBus,
		RCONAPI:    m.rconAPI,
		CommandAPI: m.commandAPI,
	}

	// Initialize plugin
	if err := pluginInstance.Init(pluginCtx); err != nil {
		cancel()
		return fmt.Errorf("plugin initialization failed: %w", err)
	}

	// Store loaded plugin
	loaded := &LoadedPlugin{
		Plugin:     pluginInstance,
		Metadata:   metadata,
		Context:    pluginCtx,
		State:      PluginStateLoaded,
		LoadedAt:   time.Now(),
		cancelFunc: cancel,
	}

	m.plugins[metadata.ID] = loaded
	m.pluginStates[metadata.ID] = PluginStateLoaded

	logger.Info("Plugin loaded", zap.String("id", metadata.ID), zap.String("version", metadata.Version))
	return nil
}

// StartAll starts all loaded plugins
func (m *Manager) StartAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, loaded := range m.plugins {
		if loaded.State == PluginStateStarted {
			continue
		}

		if err := m.startPlugin(id); err != nil {
			logger.Error("Failed to start plugin", zap.String("id", id), zap.Error(err))
			loaded.Error = err.Error()
			loaded.State = PluginStateError
			m.pluginStates[id] = PluginStateError
		}
	}

	return nil
}

// startPlugin starts a single plugin (must be called with lock held)
func (m *Manager) startPlugin(id string) error {
	loaded, exists := m.plugins[id]
	if !exists {
		return fmt.Errorf("plugin not found")
	}

	if err := loaded.Plugin.Start(); err != nil {
		return err
	}

	loaded.State = PluginStateStarted
	m.pluginStates[id] = PluginStateStarted
	loaded.Error = ""

	logger.Info("Plugin started", zap.String("id", id))
	return nil
}

// Start starts a specific plugin by ID
func (m *Manager) Start(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.startPlugin(id)
}

// StopAll stops all running plugins
func (m *Manager) StopAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, loaded := range m.plugins {
		if loaded.State != PluginStateStarted {
			continue
		}

		if err := m.stopPlugin(id); err != nil {
			logger.Error("Failed to stop plugin", zap.String("id", id), zap.Error(err))
		}
	}

	return nil
}

// stopPlugin stops a single plugin (must be called with lock held)
func (m *Manager) stopPlugin(id string) error {
	loaded, exists := m.plugins[id]
	if !exists {
		return fmt.Errorf("plugin not found")
	}

	if err := loaded.Plugin.Stop(); err != nil {
		return err
	}

	loaded.State = PluginStateStopped
	m.pluginStates[id] = PluginStateStopped

	logger.Info("Plugin stopped", zap.String("id", id))
	return nil
}

// Stop stops a specific plugin by ID
func (m *Manager) Stop(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.stopPlugin(id)
}

// Reload reloads a specific plugin's configuration
func (m *Manager) Reload(id string) error {
	m.mu.RLock()
	loaded, exists := m.plugins[id]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("plugin not found")
	}

	if err := loaded.Plugin.Reload(); err != nil {
		return err
	}

	logger.Info("Plugin reloaded", zap.String("id", id))
	return nil
}

// GetStatus returns the status of all plugins
func (m *Manager) GetStatus() []PluginStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make([]PluginStatus, 0, len(m.plugins))
	for id, loaded := range m.plugins {
		statuses = append(statuses, PluginStatus{
			ID:       id,
			Name:     loaded.Metadata.Name,
			Version:  loaded.Metadata.Version,
			State:    loaded.State,
			Enabled:  loaded.State == PluginStateStarted,
			LoadedAt: loaded.LoadedAt,
			Error:    loaded.Error,
		})
	}

	return statuses
}

// GetPluginStatus returns the status of a specific plugin
func (m *Manager) GetPluginStatus(id string) (PluginStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	loaded, exists := m.plugins[id]
	if !exists {
		return PluginStatus{}, fmt.Errorf("plugin not found")
	}

	return PluginStatus{
		ID:       id,
		Name:     loaded.Metadata.Name,
		Version:  loaded.Metadata.Version,
		State:    loaded.State,
		Enabled:  loaded.State == PluginStateStarted,
		LoadedAt: loaded.LoadedAt,
		Error:    loaded.Error,
	}, nil
}

// GlobalPluginManager is the global instance of the plugin manager
var GlobalPluginManager *Manager
