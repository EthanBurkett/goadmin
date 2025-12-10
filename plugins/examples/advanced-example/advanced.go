package advancedexample

import (
	"fmt"
	"time"

	"github.com/ethanburkett/goadmin/app/plugins"
)

// AdvancedExamplePlugin demonstrates advanced plugin features
type AdvancedExamplePlugin struct {
	ctx *plugins.PluginContext
}

func (p *AdvancedExamplePlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		ID:            "advanced-example",
		Name:          "Advanced Example Plugin",
		Version:       "2.0.0",
		Author:        "GoAdmin Team",
		Description:   "Advanced plugin demonstrating versioning, dependencies, and resource limits",
		Website:       "https://github.com/ethanburkett/goadmin",
		Dependencies:  []string{"example-plugin"}, // Example: []string{"example-plugin"} to require base plugin
		Permissions:   []string{"rcon.execute", "events.subscribe", "commands.register"},
		MinAPIVersion: "1.0.0", // Minimum required API version
		MaxAPIVersion: "2.0.0", // Maximum supported API version
		ResourceLimits: &plugins.ResourceLimits{
			MaxMemoryMB:   100,              // Limit to 100MB
			MaxCPUPercent: 50.0,             // Limit to 50% CPU
			MaxGoroutines: 50,               // Limit to 50 goroutines
			Timeout:       30 * time.Second, // 30s timeout for operations
		},
	}
}

func (p *AdvancedExamplePlugin) Init(ctx *plugins.PluginContext) error {
	p.ctx = ctx
	fmt.Printf("[AdvancedExample] Plugin initialized with API version constraints\n")
	return nil
}

func (p *AdvancedExamplePlugin) Start() error {
	fmt.Printf("[AdvancedExample] Starting plugin with resource monitoring...\n")

	// Subscribe to events
	if p.ctx.EventBus != nil {
		// Player connect event
		if err := p.ctx.EventBus.Subscribe("player.connect", func(eventType string, data map[string]interface{}) error {
			playerName, _ := data["playerName"].(string)
			fmt.Printf("[AdvancedExample] Player %s connected\n", playerName)

			// Send personalized welcome message
			if p.ctx.RCONAPI != nil {
				welcomeMsg := fmt.Sprintf(`say "^2[Advanced] Welcome %s! This server uses advanced plugin features."`, playerName)
				p.ctx.RCONAPI.SendCommand(welcomeMsg)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to subscribe to player.connect: %w", err)
		}

		// Player disconnect event
		if err := p.ctx.EventBus.Subscribe("player.disconnect", func(eventType string, data map[string]interface{}) error {
			playerName, _ := data["playerName"].(string)
			fmt.Printf("[AdvancedExample] Player %s disconnected\n", playerName)
			return nil
		}); err != nil {
			return fmt.Errorf("failed to subscribe to player.disconnect: %w", err)
		}
	}

	// Register commands
	if p.ctx.CommandAPI != nil {
		// !status command - Shows plugin info
		if err := p.ctx.CommandAPI.RegisterCommand(plugins.CommandDefinition{
			Name:        "status",
			Usage:       "status",
			Description: "Shows advanced plugin status",
			MinArgs:     0,
			MaxArgs:     0,
			Handler: func(playerName, playerGUID string, args []string) error {
				if p.ctx.RCONAPI != nil {
					msg := fmt.Sprintf(`tell %s "^2[Advanced] Version 2.0.0 | Monitoring: Active | Dependencies: OK"`, playerName)
					p.ctx.RCONAPI.SendCommand(msg)
				}
				return nil
			},
		}); err != nil {
			return fmt.Errorf("failed to register status command: %w", err)
		}

		// !info command - Shows server info
		if err := p.ctx.CommandAPI.RegisterCommand(plugins.CommandDefinition{
			Name:        "info",
			Usage:       "info",
			Description: "Shows detailed server information",
			MinArgs:     0,
			MaxArgs:     0,
			Handler: func(playerName, playerGUID string, args []string) error {
				if p.ctx.RCONAPI != nil {
					// This demonstrates resource-aware operation
					statusData, err := p.ctx.RCONAPI.GetStatus()
					if err == nil {
						msg := fmt.Sprintf(`tell %s "^2Server Info: ^7%v"`, playerName, statusData)
						p.ctx.RCONAPI.SendCommand(msg)
					} else {
						msg := fmt.Sprintf(`tell %s "^1Failed to get server info"`, playerName)
						p.ctx.RCONAPI.SendCommand(msg)
					}
				}
				return nil
			},
		}); err != nil {
			return fmt.Errorf("failed to register info command: %w", err)
		}
	}

	// Example: Dispatch a webhook for plugin startup
	if p.ctx.WebhookAPI != nil {
		p.ctx.WebhookAPI.Dispatch("plugin.started", map[string]interface{}{
			"plugin_id":   "advanced-example",
			"plugin_name": "Advanced Example Plugin",
			"version":     "2.0.0",
			"timestamp":   time.Now().Unix(),
		})
	}

	fmt.Printf("[AdvancedExample] Plugin started successfully with enhanced features\n")
	return nil
}

func (p *AdvancedExamplePlugin) Stop() error {
	fmt.Printf("[AdvancedExample] Stopping plugin...\n")

	// Cleanup: Unsubscribe from events
	if p.ctx.EventBus != nil {
		// Events are automatically cleaned up by the manager
		fmt.Printf("[AdvancedExample] Unsubscribed from events\n")
	}

	// Cleanup: Unregister commands
	if p.ctx.CommandAPI != nil {
		p.ctx.CommandAPI.UnregisterCommand("status")
		p.ctx.CommandAPI.UnregisterCommand("info")
		fmt.Printf("[AdvancedExample] Unregistered commands\n")
	}

	// Dispatch shutdown webhook
	if p.ctx.WebhookAPI != nil {
		p.ctx.WebhookAPI.Dispatch("plugin.stopped", map[string]interface{}{
			"plugin_id":   "advanced-example",
			"plugin_name": "Advanced Example Plugin",
			"timestamp":   time.Now().Unix(),
		})
	}

	fmt.Printf("[AdvancedExample] Plugin stopped successfully\n")
	return nil
}

func (p *AdvancedExamplePlugin) Reload() error {
	fmt.Printf("[AdvancedExample] Reloading plugin configuration...\n")

	// Example: Reload configuration from ConfigAPI
	if p.ctx.ConfigAPI != nil {
		// You could reload settings here
		welcomeEnabled := p.ctx.ConfigAPI.GetBool("welcome_enabled", true)
		fmt.Printf("[AdvancedExample] Welcome messages enabled: %v\n", welcomeEnabled)
	}

	fmt.Printf("[AdvancedExample] Plugin reloaded successfully\n")
	return nil
}

// Export the plugin
func New() plugins.Plugin {
	return &AdvancedExamplePlugin{}
}

func init() {
	// Register plugin on startup
	plugins.Registry.Register(&AdvancedExamplePlugin{})
}
