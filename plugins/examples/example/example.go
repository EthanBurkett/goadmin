package example

import (
	"fmt"
	"time"

	"github.com/ethanburkett/goadmin/app/plugins"
)

// ExamplePlugin is a simple plugin demonstrating the plugin API
type ExamplePlugin struct {
	ctx *plugins.PluginContext
}

// Metadata returns plugin information
func (p *ExamplePlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		ID:           "example-plugin",
		Name:         "Example Plugin",
		Version:      "1.0.0",
		Author:       "GoAdmin Team",
		Description:  "A simple example plugin demonstrating the plugin API",
		Website:      "https://github.com/ethanburkett/goadmin",
		Dependencies: []string{},
		Permissions: []string{
			"rcon.execute",
			"events.subscribe",
			"commands.register",
		},
	}
}

// Init initializes the plugin
func (p *ExamplePlugin) Init(ctx *plugins.PluginContext) error {
	p.ctx = ctx
	return nil
}

// Start starts the plugin
func (p *ExamplePlugin) Start() error {
	// Subscribe to player events if EventBus is available
	if p.ctx.EventBus != nil {
		if err := p.ctx.EventBus.Subscribe("player.connect", func(eventType string, data map[string]interface{}) error {
			playerName, _ := data["playerName"].(string)
			playerGUID, _ := data["playerGUID"].(string)
			fmt.Printf("[ExamplePlugin] Player connected: %s (GUID: %s)\n", playerName, playerGUID)

			// Send welcome message via RCON
			if p.ctx.RCONAPI != nil {
				welcomeMsg := fmt.Sprintf(`tell %s "^2Welcome to the server, %s!"`, playerName, playerName)
				p.ctx.RCONAPI.SendCommand(welcomeMsg)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to subscribe to player.connect: %w", err)
		}

		if err := p.ctx.EventBus.Subscribe("player.disconnect", func(eventType string, data map[string]interface{}) error {
			playerName, _ := data["playerName"].(string)
			playerGUID, _ := data["playerGUID"].(string)
			fmt.Printf("[ExamplePlugin] Player disconnected: %s (GUID: %s) at %s\n", playerName, playerGUID, time.Now().Format(time.RFC3339))
			return nil
		}); err != nil {
			return fmt.Errorf("failed to subscribe to player.disconnect: %w", err)
		}
	}

	// Register custom commands
	if p.ctx.CommandAPI != nil {
		// Example command: !hello - Sends a greeting
		if err := p.ctx.CommandAPI.RegisterCommand(plugins.CommandDefinition{
			Name:        "hello",
			Usage:       "hello",
			Description: "Says hello to the player",
			MinArgs:     0,
			MaxArgs:     0,
			Handler: func(playerName, playerGUID string, args []string) error {
				if p.ctx.RCONAPI != nil {
					msg := fmt.Sprintf(`tell %s "^2Hello ^7%s^2! This is a plugin command!"`, playerName, playerName)
					p.ctx.RCONAPI.SendCommand(msg)
				}
				return nil
			},
		}); err != nil {
			return fmt.Errorf("failed to register hello command: %w", err)
		}

		// Example command: !time - Shows current server time
		if err := p.ctx.CommandAPI.RegisterCommand(plugins.CommandDefinition{
			Name:        "time",
			Usage:       "time",
			Description: "Shows the current server time",
			MinArgs:     0,
			MaxArgs:     0,
			Handler: func(playerName, playerGUID string, args []string) error {
				if p.ctx.RCONAPI != nil {
					currentTime := time.Now().Format("15:04:05 MST")
					msg := fmt.Sprintf(`tell %s "^3Server time: ^7%s"`, playerName, currentTime)
					p.ctx.RCONAPI.SendCommand(msg)
				}
				return nil
			},
		}); err != nil {
			return fmt.Errorf("failed to register time command: %w", err)
		}

		// Example command: !echo <message> - Echoes back a message
		if err := p.ctx.CommandAPI.RegisterCommand(plugins.CommandDefinition{
			Name:        "echo",
			Usage:       "echo <message>",
			Description: "Echoes back your message",
			MinArgs:     1,
			MaxArgs:     -1, // unlimited
			Handler: func(playerName, playerGUID string, args []string) error {
				if p.ctx.RCONAPI != nil {
					message := ""
					for _, arg := range args {
						message += arg + " "
					}
					msg := fmt.Sprintf(`tell %s "^6Echo: ^7%s"`, playerName, message)
					p.ctx.RCONAPI.SendCommand(msg)
				}
				return nil
			},
		}); err != nil {
			return fmt.Errorf("failed to register echo command: %w", err)
		}
	}

	fmt.Printf("[ExamplePlugin] Started successfully\n")
	return nil
}

// Stop stops the plugin
func (p *ExamplePlugin) Stop() error {
	// Unregister commands
	if p.ctx.CommandAPI != nil {
		p.ctx.CommandAPI.UnregisterCommand("hello")
		p.ctx.CommandAPI.UnregisterCommand("time")
		p.ctx.CommandAPI.UnregisterCommand("echo")
	}

	fmt.Printf("[ExamplePlugin] Stopped\n")
	return nil
}

// Reload reloads the plugin configuration
func (p *ExamplePlugin) Reload() error {
	fmt.Printf("[ExamplePlugin] Reloaded\n")
	return nil
}

// Register the plugin with the global registry
func init() {
	plugins.Registry.Register(&ExamplePlugin{})
}
