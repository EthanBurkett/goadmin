package commands

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/ethanburkett/goadmin/app/rcon"
)

type CommandHandler struct {
	rcon      *rcon.Client
	callbacks map[string]CommandCallback
}

type CommandCallback func(ch *CommandHandler, playerName, playerGUID string, args []string) error

func NewCommandHandler(rconClient *rcon.Client) *CommandHandler {
	handler := &CommandHandler{
		rcon:      rconClient,
		callbacks: make(map[string]CommandCallback),
	}

	handler.registerBuiltInCallbacks()

	return handler
}

// ProcessChatCommand processes a chat message that starts with !
func (ch *CommandHandler) ProcessChatCommand(playerName, playerGUID, message string) error {
	message = strings.TrimPrefix(message, "!")

	parts := strings.Fields(message)
	if len(parts) == 0 {
		return nil
	}

	commandName := strings.ToLower(parts[0])
	args := parts[1:]

	if commandName == "iamgod" {
		return ch.processIamGod(playerName, playerGUID)
	}

	if callback, exists := ch.callbacks[commandName]; exists {
		cmd, err := models.GetCustomCommand(commandName)
		if err != nil {
			logger.Debug(fmt.Sprintf("Built-in command '%s' not found in database", commandName))
			return nil
		}

		if !cmd.Enabled {
			logger.Debug(fmt.Sprintf("Built-in command '%s' is disabled", commandName))
			return nil
		}

		// Validate argument count
		if len(args) < cmd.MinArgs {
			ch.sendPlayerMessage(playerName, fmt.Sprintf("Usage: %s", cmd.Usage))
			return nil
		}
		if cmd.MaxArgs >= 0 && len(args) > cmd.MaxArgs {
			ch.sendPlayerMessage(playerName, fmt.Sprintf("Usage: %s", cmd.Usage))
			return nil
		}

		// Get player's power level and permissions
		playerPower := models.GetPlayerPower(playerGUID)
		playerPermissions := ch.getPlayerPermissions(playerGUID)

		// Check requirements based on RequirementType
		requirementType := cmd.RequirementType
		if requirementType == "" {
			requirementType = "both"
		}

		if requirementType == "power" || requirementType == "both" {
			if playerPower < cmd.MinPower {
				ch.sendPlayerMessage(playerName, fmt.Sprintf("Insufficient power level (need %d, have %d)", cmd.MinPower, playerPower))
				logger.Info(fmt.Sprintf("Player %s (%s) denied access to built-in command '%s' - insufficient power (%d < %d)",
					playerName, playerGUID, commandName, playerPower, cmd.MinPower))
				return nil
			}
		}

		if requirementType == "permission" || requirementType == "both" {
			if !ch.hasRequiredPermissions(playerPermissions, `["all"]`) && !ch.hasRequiredPermissions(playerPermissions, cmd.Permissions) {
				ch.sendPlayerMessage(playerName, "You don't have permission to use this command")
				logger.Info(fmt.Sprintf("Player %s (%s) denied access to built-in command '%s' - missing permissions",
					playerName, playerGUID, commandName))
				return nil
			}
		}

		logger.Info(fmt.Sprintf("Executing built-in callback for command '%s' by player %s", commandName, playerName))
		if err := callback(ch, playerName, playerGUID, args); err != nil {
			logger.Error(fmt.Sprintf("Callback for command '%s' failed: %v", commandName, err))
			ch.sendPlayerMessage(playerName, "Command failed to execute")
			return err
		}
		return nil
	}

	cmd, err := models.GetCustomCommand(commandName)
	if err != nil {
		logger.Debug(fmt.Sprintf("Command '%s' not found in database", commandName))
		return nil
	}

	if !cmd.Enabled {
		logger.Debug(fmt.Sprintf("Command '%s' is disabled", commandName))
		return nil
	}

	if len(args) < cmd.MinArgs {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("Usage: %s", cmd.Usage))
		return nil
	}
	if cmd.MaxArgs >= 0 && len(args) > cmd.MaxArgs {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("Usage: %s", cmd.Usage))
		return nil
	}

	// Get player's power level and permissions
	playerPower := models.GetPlayerPower(playerGUID)
	playerPermissions := ch.getPlayerPermissions(playerGUID)

	// Check requirements based on RequirementType
	requirementType := cmd.RequirementType
	if requirementType == "" {
		requirementType = "both"
	}

	if requirementType == "power" || requirementType == "both" {
		if playerPower < cmd.MinPower {
			ch.sendPlayerMessage(playerName, fmt.Sprintf("Insufficient power level (need %d, have %d)", cmd.MinPower, playerPower))
			logger.Info(fmt.Sprintf("Player %s (%s) denied access to command '%s' - insufficient power (%d < %d)",
				playerName, playerGUID, commandName, playerPower, cmd.MinPower))
			return nil
		}
	}

	if requirementType == "permission" || requirementType == "both" {
		if !ch.hasRequiredPermissions(playerPermissions, `["all"]`) && !ch.hasRequiredPermissions(playerPermissions, cmd.Permissions) {
			ch.sendPlayerMessage(playerName, "You don't have permission to use this command")
			logger.Info(fmt.Sprintf("Player %s (%s) denied access to command '%s' - missing permissions",
				playerName, playerGUID, commandName))
			return nil
		}
	}

	rconCmd := ch.buildRconCommand(cmd.RconCommand, args, playerName, playerGUID)

	logger.Info(fmt.Sprintf("Executing custom command '%s' for player %s: %s", commandName, playerName, rconCmd))

	response, err := ch.rcon.SendCommand(rconCmd)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to execute command '%s': %v", commandName, err))
		ch.sendPlayerMessage(playerName, "Command failed to execute")
		return err
	}

	logger.Debug(fmt.Sprintf("Command response: %s", response))

	return nil
}

// buildRconCommand replaces placeholders in the command template
func (ch *CommandHandler) buildRconCommand(template string, args []string, playerName, playerGUID string) string {
	result := template

	result = ch.resolvePlayerIdPlaceholders(result, args)
	result = ch.resolveArgsFromPlaceholders(result, args)

	for i, arg := range args {
		placeholder := fmt.Sprintf("{arg%d}", i)
		result = strings.ReplaceAll(result, placeholder, arg)
	}

	result = strings.ReplaceAll(result, "{player}", playerName)
	result = strings.ReplaceAll(result, "{guid}", playerGUID)

	re := regexp.MustCompile(`\{arg\d+\}`)
	result = re.ReplaceAllString(result, "")

	return strings.TrimSpace(result)
}

// getPlayerPermissions gets all permissions for a player based on their group
func (ch *CommandHandler) getPlayerPermissions(guid string) []string {
	player, err := models.GetInGamePlayerByGUID(guid)
	if err != nil || player.Group == nil {
		return []string{}
	}

	var permissions []string
	if player.Group.Permissions != "" {
		if err := json.Unmarshal([]byte(player.Group.Permissions), &permissions); err != nil {
			logger.Error(fmt.Sprintf("Failed to parse permissions for group %s: %v", player.Group.Name, err))
			return []string{}
		}
	}

	return permissions
}

// hasRequiredPermissions checks if player has all required permissions
func (ch *CommandHandler) hasRequiredPermissions(playerPerms []string, requiredPermsJSON string) bool {
	if requiredPermsJSON == "" {
		return true
	}

	var requiredPerms []string
	if err := json.Unmarshal([]byte(requiredPermsJSON), &requiredPerms); err != nil {
		logger.Error(fmt.Sprintf("Failed to parse required permissions: %v", err))
		return false
	}

	if len(requiredPerms) == 0 {
		return true
	}

	playerPermSet := make(map[string]bool)
	for _, perm := range playerPerms {
		playerPermSet[perm] = true
	}

	for _, reqPerm := range requiredPerms {
		if !playerPermSet[reqPerm] {
			return false
		}
	}

	return true
}

// sendPlayerMessage sends a message to a specific player
func (ch *CommandHandler) sendPlayerMessage(playerName, message string) {
	cmd := fmt.Sprintf("tell %s ^3[Command] ^7%s", playerName, message)
	ch.rcon.SendCommand(cmd)
}

// resolveArgsFromPlaceholders resolves {argsFrom:N} placeholders by joining args from index N onwards
func (ch *CommandHandler) resolveArgsFromPlaceholders(template string, args []string) string {
	re := regexp.MustCompile(`\{argsFrom:(\d+)\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	result := template
	for _, match := range matches {
		fullPlaceholder := match[0]
		startIndex := match[1]

		var startNum int
		fmt.Sscanf(startIndex, "%d", &startNum)

		if startNum < len(args) {
			joinedArgs := strings.Join(args[startNum:], " ")
			result = strings.ReplaceAll(result, fullPlaceholder, joinedArgs)
		} else {
			result = strings.ReplaceAll(result, fullPlaceholder, "")
		}
	}

	return result
}

// resolvePlayerIdPlaceholders resolves {playerId:argN} placeholders to entity IDs
func (ch *CommandHandler) resolvePlayerIdPlaceholders(template string, args []string) string {
	re := regexp.MustCompile(`\{playerId:arg(\d+)\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	result := template
	for _, match := range matches {
		fullPlaceholder := match[0]
		argIndex := match[1]

		var argNum int
		fmt.Sscanf(argIndex, "%d", &argNum)

		if argNum < len(args) {
			playerName := args[argNum]
			playerID := ch.findPlayerByName(playerName)
			if playerID != "" {
				result = strings.ReplaceAll(result, fullPlaceholder, playerID)
			} else {
				result = strings.ReplaceAll(result, fullPlaceholder, playerName)
			}
		}
	}

	return result
}

// findPlayerByName searches for a player by name and returns their entity ID
func (ch *CommandHandler) findPlayerByName(searchName string) string {
	status, err := ch.rcon.Status()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get server status: %v", err))
		return ""
	}

	searchName = strings.ToLower(searchName)

	for _, player := range status.Players {
		if strings.ToLower(player.StrippedName) == searchName {
			return fmt.Sprintf("%d", player.ID)
		}
	}

	for _, player := range status.Players {
		if strings.Contains(strings.ToLower(player.StrippedName), searchName) {
			return fmt.Sprintf("%d", player.ID)
		}
	}

	return ""
}

// registerBuiltInCallbacks registers all built-in command callbacks
func (ch *CommandHandler) registerBuiltInCallbacks() {
	ch.callbacks["groups"] = ch.handleGroupsCommand
	ch.callbacks["mygroup"] = ch.handleMyGroupCommand
	ch.callbacks["putgroup"] = ch.handlePutGroupCommand
	ch.callbacks["adminlist"] = ch.handleAdminListCommand
	ch.callbacks["help"] = ch.handleHelpCommand
	ch.callbacks["report"] = ch.handleReportCommand
	ch.callbacks["tempban"] = ch.handleTempBanCommand
}
