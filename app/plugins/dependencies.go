package plugins

import (
	"fmt"
)

// DependencyValidator validates plugin dependencies
type DependencyValidator struct {
	registry *PluginRegistry
	manager  *Manager
}

// NewDependencyValidator creates a new dependency validator
func NewDependencyValidator(registry *PluginRegistry, manager *Manager) *DependencyValidator {
	return &DependencyValidator{
		registry: registry,
		manager:  manager,
	}
}

// ValidateDependencies checks if all dependencies for a plugin are satisfied
func (v *DependencyValidator) ValidateDependencies(metadata PluginMetadata) error {
	if len(metadata.Dependencies) == 0 {
		return nil
	}

	missing := []string{}
	incompatible := []string{}

	for _, depID := range metadata.Dependencies {
		// Check if dependency is registered
		depPlugin, exists := v.registry.Get(depID)
		if !exists {
			missing = append(missing, depID)
			continue
		}

		// Check if dependency is loaded and started
		if v.manager != nil {
			status, err := v.manager.GetPluginStatus(depID)
			if err != nil {
				incompatible = append(incompatible, fmt.Sprintf("%s (not found)", depID))
				continue
			}
			// Dependency must be at least loaded (not necessarily started yet during initial load)
			if status.State != PluginStateLoaded && status.State != PluginStateStarted {
				incompatible = append(incompatible, fmt.Sprintf("%s (state: %s)", depID, status.State))
				continue
			}
		}

		// Additional version compatibility checks could go here
		_ = depPlugin
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing dependencies: %v", missing)
	}

	if len(incompatible) > 0 {
		return fmt.Errorf("incompatible dependencies: %v", incompatible)
	}

	return nil
}

// GetDependencyTree builds a dependency tree for a plugin
func (v *DependencyValidator) GetDependencyTree(pluginID string) (map[string][]string, error) {
	tree := make(map[string][]string)
	visited := make(map[string]bool)

	if err := v.buildTree(pluginID, tree, visited); err != nil {
		return nil, err
	}

	return tree, nil
}

// buildTree recursively builds the dependency tree
func (v *DependencyValidator) buildTree(pluginID string, tree map[string][]string, visited map[string]bool) error {
	if visited[pluginID] {
		return fmt.Errorf("circular dependency detected: %s", pluginID)
	}

	visited[pluginID] = true
	defer func() { visited[pluginID] = false }()

	plugin, exists := v.registry.Get(pluginID)
	if !exists {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	metadata := plugin.Metadata()
	tree[pluginID] = metadata.Dependencies

	for _, depID := range metadata.Dependencies {
		if err := v.buildTree(depID, tree, visited); err != nil {
			return err
		}
	}

	return nil
}

// GetLoadOrder determines the correct order to load plugins based on dependencies
func (v *DependencyValidator) GetLoadOrder(pluginIDs []string) ([]string, error) {
	// Build dependency graph
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize all plugins with 0 in-degree
	for _, id := range pluginIDs {
		inDegree[id] = 0
	}

	for _, id := range pluginIDs {
		plugin, exists := v.registry.Get(id)
		if !exists {
			return nil, fmt.Errorf("plugin not found: %s", id)
		}

		metadata := plugin.Metadata()
		graph[id] = metadata.Dependencies

		// For each dependency, increment the in-degree of THIS plugin (not the dependency)
		for range metadata.Dependencies {
			inDegree[id]++
		}
	}

	// Topological sort using Kahn's algorithm
	order := []string{}
	queue := []string{}

	// Find all plugins with no dependencies (in-degree = 0)
	for id := range inDegree {
		if inDegree[id] == 0 {
			queue = append(queue, id)
		}
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)

		// For each plugin that depends on current, decrement its in-degree
		for otherID, deps := range graph {
			for _, depID := range deps {
				if depID == current {
					inDegree[otherID]--
					if inDegree[otherID] == 0 {
						queue = append(queue, otherID)
					}
				}
			}
		}
	}

	// Check for circular dependencies
	if len(order) != len(pluginIDs) {
		return nil, fmt.Errorf("circular dependency detected in plugins")
	}

	return order, nil
}
