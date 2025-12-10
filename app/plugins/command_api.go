package plugins

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
	"go.uber.org/zap"
)

// CommandAPIImpl implements the CommandAPI interface for plugins
type CommandAPIImpl struct {
	mu               sync.RWMutex
	pluginCommands   map[string]*PluginCommand // command name -> plugin command
	commandCallbacks map[string]CommandHandler // command name -> handler
	rconAPI          RCONAPI                   // For sending messages to players
}

// PluginCommand represents a command registered by a plugin
type PluginCommand struct {
	PluginID   string
	Definition CommandDefinition
}

// NewCommandAPI creates a new Command API instance
func NewCommandAPI(rconAPI RCONAPI) *CommandAPIImpl {
	return &CommandAPIImpl{
		pluginCommands:   make(map[string]*PluginCommand),
		commandCallbacks: make(map[string]CommandHandler),
		rconAPI:          rconAPI,
	}
}

// RegisterCommand registers a custom in-game command
func (c *CommandAPIImpl) RegisterCommand(cmd CommandDefinition) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if command already exists
	if _, exists := c.pluginCommands[cmd.Name]; exists {
		return fmt.Errorf("command '%s' is already registered", cmd.Name)
	}

	// Validate command
	if cmd.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	if cmd.Handler == nil {
		return fmt.Errorf("command handler cannot be nil")
	}

	// Store the command
	pluginCmd := &PluginCommand{
		PluginID:   "", // Will be set by manager when plugin registers it
		Definition: cmd,
	}
	c.pluginCommands[cmd.Name] = pluginCmd
	c.commandCallbacks[cmd.Name] = cmd.Handler

	logger.Info("Plugin command registered", zap.String("command", cmd.Name))
	return nil
}

// UnregisterCommand removes a custom command
func (c *CommandAPIImpl) UnregisterCommand(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.pluginCommands[name]; !exists {
		return fmt.Errorf("command '%s' is not registered", name)
	}

	delete(c.pluginCommands, name)
	delete(c.commandCallbacks, name)

	logger.Info("Plugin command unregistered", zap.String("command", name))
	return nil
}

// ExecuteCommand executes a command programmatically
func (c *CommandAPIImpl) ExecuteCommand(playerName, playerGUID, command string, args []string) error {
	c.mu.RLock()
	handler, exists := c.commandCallbacks[command]
	c.mu.RUnlock()

	if !exists {
		return fmt.Errorf("command '%s' not found", command)
	}

	return handler(playerName, playerGUID, args)
}

// ProcessPluginCommand processes a command from chat
func (c *CommandAPIImpl) ProcessPluginCommand(playerName, playerGUID, commandName string, args []string) error {
	c.mu.RLock()
	pluginCmd, exists := c.pluginCommands[commandName]
	handler := c.commandCallbacks[commandName]
	c.mu.RUnlock()

	if !exists {
		return nil // Not a plugin command
	}

	cmd := pluginCmd.Definition

	// Validate argument count
	if len(args) < cmd.MinArgs {
		c.sendPlayerMessage(playerName, fmt.Sprintf("Usage: !%s", cmd.Usage))
		return nil
	}
	if cmd.MaxArgs >= 0 && len(args) > cmd.MaxArgs {
		c.sendPlayerMessage(playerName, fmt.Sprintf("Usage: !%s", cmd.Usage))
		return nil
	}

	// Check permissions if required
	if cmd.MinPower > 0 || len(cmd.Permissions) > 0 {
		// Get player's power level from their group
		playerPower := models.GetPlayerPower(playerGUID)

		// Check power level
		if cmd.MinPower > 0 && playerPower < cmd.MinPower {
			c.sendPlayerMessage(playerName, "You don't have permission to use this command")
			return nil
		}

		// Check permissions (if specified)
		if len(cmd.Permissions) > 0 {
			// Get player's group to check permissions
			player, err := models.GetInGamePlayerByGUID(playerGUID)
			if err != nil || player.Group == nil {
				c.sendPlayerMessage(playerName, "You don't have permission to use this command")
				return nil
			}

			// Parse group permissions (JSON array of strings)
			hasPermission := false
			if player.Group.Permissions != "" {
				// Simple check - if any required permission is in the group's permissions
				for _, reqPerm := range cmd.Permissions {
					if strings.Contains(player.Group.Permissions, reqPerm) {
						hasPermission = true
						break
					}
				}
			}
			if !hasPermission {
				c.sendPlayerMessage(playerName, "You don't have permission to use this command")
				return nil
			}
		}
	}

	// Execute handler
	if err := handler(playerName, playerGUID, args); err != nil {
		logger.Error("Plugin command handler error",
			zap.String("command", commandName),
			zap.String("player", playerName),
			zap.Error(err))
		c.sendPlayerMessage(playerName, "An error occurred while executing the command")
		return err
	}

	return nil
}

// GetRegisteredCommands returns all registered plugin commands
func (c *CommandAPIImpl) GetRegisteredCommands() map[string]CommandDefinition {
	c.mu.RLock()
	defer c.mu.RUnlock()

	commands := make(map[string]CommandDefinition)
	for name, pluginCmd := range c.pluginCommands {
		commands[name] = pluginCmd.Definition
	}
	return commands
}

// sendPlayerMessage sends a message to a specific player
func (c *CommandAPIImpl) sendPlayerMessage(playerName, message string) {
	if c.rconAPI == nil {
		return
	}

	// Use tell command to send private message
	cmd := fmt.Sprintf(`tell %s "^2%s"`, playerName, message)
	if _, err := c.rconAPI.SendCommand(cmd); err != nil {
		logger.Error("Failed to send player message", zap.Error(err))
	}
}

// GetCommandCount returns the number of registered plugin commands
func (c *CommandAPIImpl) GetCommandCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.pluginCommands)
}
