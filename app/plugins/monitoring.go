package plugins

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/ethanburkett/goadmin/app/logger"
	"go.uber.org/zap"
)

// ResourceMonitor monitors plugin resource usage
type ResourceMonitor struct {
	mu              sync.RWMutex
	pluginMetrics   map[string]*PluginMetrics
	monitorInterval time.Duration
	ctx             context.Context
	cancel          context.CancelFunc
}

// PluginMetrics tracks resource usage for a plugin
type PluginMetrics struct {
	PluginID       string
	MemoryUsageMB  float64
	GoroutineCount int
	LastChecked    time.Time
	ViolationCount int
	Throttled      bool
}

// NewResourceMonitor creates a new resource monitor
func NewResourceMonitor(interval time.Duration) *ResourceMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &ResourceMonitor{
		pluginMetrics:   make(map[string]*PluginMetrics),
		monitorInterval: interval,
		ctx:             ctx,
		cancel:          cancel,
	}
}

// Start begins monitoring plugin resources
func (m *ResourceMonitor) Start() {
	// Do an immediate check before starting the monitor loop
	m.checkResources()
	go m.monitorLoop()
}

// Stop stops the resource monitor
func (m *ResourceMonitor) Stop() {
	m.cancel()
}

// monitorLoop periodically checks plugin resource usage
func (m *ResourceMonitor) monitorLoop() {
	ticker := time.NewTicker(m.monitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.checkResources()
		}
	}
}

// checkResources checks current resource usage
func (m *ResourceMonitor) checkResources() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.mu.Lock()
	defer m.mu.Unlock()

	// Update metrics for all monitored plugins
	for pluginID := range m.pluginMetrics {
		goroutineCount := runtime.NumGoroutine()

		m.pluginMetrics[pluginID] = &PluginMetrics{
			PluginID:       pluginID,
			MemoryUsageMB:  float64(memStats.Alloc) / 1024 / 1024,
			GoroutineCount: goroutineCount,
			LastChecked:    time.Now(),
		}
	}
}

// RegisterPlugin registers a plugin for resource monitoring
func (m *ResourceMonitor) RegisterPlugin(pluginID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.pluginMetrics[pluginID]; !exists {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		goroutineCount := runtime.NumGoroutine()

		m.pluginMetrics[pluginID] = &PluginMetrics{
			PluginID:       pluginID,
			MemoryUsageMB:  float64(memStats.Alloc) / 1024 / 1024,
			GoroutineCount: goroutineCount,
			LastChecked:    time.Now(),
		}
	}
}

// UnregisterPlugin removes a plugin from monitoring
func (m *ResourceMonitor) UnregisterPlugin(pluginID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.pluginMetrics, pluginID)
}

// GetMetrics returns current metrics for a plugin
func (m *ResourceMonitor) GetMetrics(pluginID string) (*PluginMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics, exists := m.pluginMetrics[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin not being monitored: %s", pluginID)
	}

	return metrics, nil
}

// CheckLimits checks if a plugin is violating its resource limits
func (m *ResourceMonitor) CheckLimits(pluginID string, limits *ResourceLimits) error {
	if limits == nil {
		return nil
	}

	metrics, err := m.GetMetrics(pluginID)
	if err != nil {
		return err
	}

	violations := []string{}

	if limits.MaxMemoryMB > 0 && metrics.MemoryUsageMB > float64(limits.MaxMemoryMB) {
		violations = append(violations, fmt.Sprintf("memory usage %.2fMB exceeds limit %dMB",
			metrics.MemoryUsageMB, limits.MaxMemoryMB))
	}

	if limits.MaxGoroutines > 0 && metrics.GoroutineCount > limits.MaxGoroutines {
		violations = append(violations, fmt.Sprintf("goroutine count %d exceeds limit %d",
			metrics.GoroutineCount, limits.MaxGoroutines))
	}

	if len(violations) > 0 {
		m.mu.Lock()
		metrics.ViolationCount++
		m.mu.Unlock()

		return fmt.Errorf("resource limit violations for plugin %s: %v", pluginID, violations)
	}

	return nil
}

// GetAllMetrics returns metrics for all monitored plugins
func (m *ResourceMonitor) GetAllMetrics() map[string]*PluginMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make(map[string]*PluginMetrics, len(m.pluginMetrics))
	for id, m := range m.pluginMetrics {
		metricsCopy := *m
		metrics[id] = &metricsCopy
	}

	return metrics
}

// HotReloader manages hot-reloading of plugins
type HotReloader struct {
	manager *Manager
	mu      sync.Mutex
}

// NewHotReloader creates a new hot reloader
func NewHotReloader(manager *Manager) *HotReloader {
	return &HotReloader{
		manager: manager,
	}
}

// Reload performs a hot reload of a plugin
// This stops the plugin, reloads its configuration, and starts it again
func (r *HotReloader) Reload(pluginID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	logger.Info("Hot-reloading plugin", zap.String("id", pluginID))

	// Get current plugin status
	status, err := r.manager.GetPluginStatus(pluginID)
	if err != nil {
		return fmt.Errorf("failed to get plugin status: %w", err)
	}

	wasStarted := status.State == PluginStateStarted

	// Stop the plugin if it's running
	if wasStarted {
		if err := r.manager.Stop(pluginID); err != nil {
			return fmt.Errorf("failed to stop plugin during reload: %w", err)
		}
	}

	// Call the plugin's Reload method
	if err := r.manager.Reload(pluginID); err != nil {
		// Try to restart if it was running before
		if wasStarted {
			if startErr := r.manager.Start(pluginID); startErr != nil {
				logger.Error("Failed to restart plugin after reload failure",
					zap.String("id", pluginID),
					zap.Error(startErr))
			}
		}
		return fmt.Errorf("failed to reload plugin: %w", err)
	}

	// Restart the plugin if it was running before
	if wasStarted {
		if err := r.manager.Start(pluginID); err != nil {
			return fmt.Errorf("failed to restart plugin after reload: %w", err)
		}
	}

	logger.Info("Plugin hot-reloaded successfully", zap.String("id", pluginID))
	return nil
}

// ReloadAll performs a hot reload of all running plugins
func (r *HotReloader) ReloadAll() error {
	statuses := r.manager.GetStatus()

	errors := []error{}
	for _, status := range statuses {
		if status.State == PluginStateStarted {
			if err := r.Reload(status.ID); err != nil {
				errors = append(errors, err)
				logger.Error("Failed to reload plugin",
					zap.String("id", status.ID),
					zap.Error(err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to reload %d plugin(s)", len(errors))
	}

	return nil
}
