package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
)

// handleGroupsCommand shows all available groups with their power levels
func (ch *CommandHandler) handleGroupsCommand(ch2 *CommandHandler, playerName, playerGUID string, args []string) error {
	groups, err := models.GetAllGroups()
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to fetch groups")
		return err
	}

	if len(groups) == 0 {
		ch.sendPlayerMessage(playerName, "No groups configured")
		return nil
	}

	ch.sendPlayerMessage(playerName, "^3Available Groups:")
	for _, group := range groups {
		msg := fmt.Sprintf("^2%s ^7(Power: ^3%d^7)", group.Name, group.Power)
		ch.sendPlayerMessage(playerName, msg)
	}

	return nil
}

// handleMyGroupCommand shows the player's current group and permissions
func (ch *CommandHandler) handleMyGroupCommand(ch2 *CommandHandler, playerName, playerGUID string, args []string) error {
	player, err := models.GetInGamePlayerByGUID(playerGUID)
	if err != nil || player == nil {
		ch.sendPlayerMessage(playerName, "You are not assigned to any group")
		return nil
	}

	if player.GroupID == nil {
		ch.sendPlayerMessage(playerName, "You are not assigned to any group")
		return nil
	}

	group, err := models.GetGroupByID(*player.GroupID)
	if err != nil || group == nil {
		ch.sendPlayerMessage(playerName, "Group not found")
		return nil
	}

	ch.sendPlayerMessage(playerName, fmt.Sprintf("^3Your Group: ^2%s", group.Name))
	ch.sendPlayerMessage(playerName, fmt.Sprintf("^3Power Level: ^2%d", group.Power))

	var permissions []string
	if err := json.Unmarshal([]byte(group.Permissions), &permissions); err == nil && len(permissions) > 0 {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("^3Permissions: ^7%s", strings.Join(permissions, ", ")))
	} else {
		ch.sendPlayerMessage(playerName, "^3Permissions: ^7None")
	}

	return nil
}

// handlePutGroupCommand assigns a player to a group
func (ch *CommandHandler) handlePutGroupCommand(ch2 *CommandHandler, playerName, playerGUID string, args []string) error {
	if len(args) < 2 {
		ch.sendPlayerMessage(playerName, "Usage: !putgroup <player> <group>")
		return nil
	}

	targetPlayerName := args[0]
	groupName := args[1]

	status, err := ch.rcon.Status()
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to get server status")
		return err
	}

	var targetGUID string
	searchName := strings.ToLower(targetPlayerName)
	for _, player := range status.Players {
		if strings.ToLower(player.StrippedName) == searchName || strings.Contains(strings.ToLower(player.StrippedName), searchName) {
			targetGUID = player.Uuid
			break
		}
	}

	if targetGUID == "" {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("Player '%s' not found online", targetPlayerName))
		return nil
	}

	if targetGUID == playerGUID {
		ch.sendPlayerMessage(playerName, "You cannot change your own group")
		return nil
	}

	groups, err := models.GetAllGroups()
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to fetch groups")
		return err
	}

	var targetGroup *models.Group
	for _, g := range groups {
		if strings.EqualFold(g.Name, groupName) {
			targetGroup = &g
			break
		}
	}

	if targetGroup == nil {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("Group '%s' not found", groupName))
		return nil
	}

	player, err := models.CreateOrUpdateInGamePlayer(targetGUID, targetPlayerName)
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to create player record")
		return err
	}

	if targetGroup.Power > models.GetPlayerPower(targetGUID) {
		ch.sendPlayerMessage(playerName, "You cannot assign a group with higher power than your own")
		return nil
	}

	if err := models.AssignPlayerToGroup(player.ID, targetGroup.ID); err != nil {
		ch.sendPlayerMessage(playerName, "Failed to assign group")
		return err
	}

	ch.sendPlayerMessage(playerName, fmt.Sprintf("^2%s ^7assigned to group ^2%s", targetPlayerName, targetGroup.Name))
	logger.Info(fmt.Sprintf("%s assigned %s to group %s", playerName, targetPlayerName, targetGroup.Name))

	return nil
}
