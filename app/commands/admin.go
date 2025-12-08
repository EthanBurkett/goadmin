package commands

import (
	"fmt"
	"strings"

	"github.com/ethanburkett/goadmin/app/models"
)

// handleAdminListCommand lists all online admins (power >= 80 or in Admin group)
func (ch *CommandHandler) handleAdminListCommand(ch2 *CommandHandler, playerName, playerGUID string, args []string) error {
	status, err := ch.rcon.Status()
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to get server status")
		return err
	}

	if len(status.Players) == 0 {
		ch.sendPlayerMessage(playerName, "No players online")
		return nil
	}

	groups, err := models.GetAllGroups()
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to fetch groups")
		return err
	}

	var adminGroupID *uint
	for _, g := range groups {
		if strings.EqualFold(g.Name, "Admin") {
			adminGroupID = &g.ID
			break
		}
	}

	var admins []string
	for _, player := range status.Players {
		inGamePlayer, err := models.GetInGamePlayerByGUID(player.Uuid)
		if err != nil || inGamePlayer == nil {
			continue
		}

		if inGamePlayer.GroupID == nil {
			continue
		}

		group, err := models.GetGroupByID(*inGamePlayer.GroupID)
		if err != nil || group == nil {
			continue
		}

		isAdmin := group.Power >= 80
		if adminGroupID != nil && *inGamePlayer.GroupID == *adminGroupID {
			isAdmin = true
		}

		if isAdmin {
			admins = append(admins, fmt.Sprintf("%s (^2%s^7)", player.StrippedName, group.Name))
		}
	}

	if len(admins) == 0 {
		ch.sendPlayerMessage(playerName, "No admins currently online")
		return nil
	}

	ch.sendPlayerMessage(playerName, "^3Online Admins:")
	for _, admin := range admins {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("^7%s", admin))
	}

	return nil
}

// handleHelpCommand shows paginated list of available commands
func (ch *CommandHandler) handleHelpCommand(ch2 *CommandHandler, playerName, playerGUID string, args []string) error {
	playerPower := models.GetPlayerPower(playerGUID)
	playerPermissions := ch.getPlayerPermissions(playerGUID)

	page := 1
	if len(args) > 0 {
		fmt.Sscanf(args[0], "%d", &page)
		if page < 1 {
			page = 1
		}
	}

	allCommands, err := models.GetAllCustomCommands()
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to fetch commands")
		return err
	}

	var availableCommands []models.CustomCommand
	for _, cmd := range allCommands {
		if !cmd.Enabled {
			continue
		}

		requirementType := cmd.RequirementType
		if requirementType == "" {
			requirementType = "both"
		}

		canUse := true

		if requirementType == "power" || requirementType == "both" {
			if playerPower < cmd.MinPower {
				canUse = false
			}
		}

		if canUse && (requirementType == "permission" || requirementType == "both") {
			if !ch.hasRequiredPermissions(playerPermissions, `["all"]`) && !ch.hasRequiredPermissions(playerPermissions, cmd.Permissions) {
				canUse = false
			}
		}

		if canUse {
			availableCommands = append(availableCommands, cmd)
		}
	}

	if !models.HasSetting("ingame_iamgod_used", "true") {
		iamgodCmd := models.CustomCommand{
			Name:        "iamgod",
			Usage:       "!iamgod",
			Description: "Claim Owner privileges (first use only)",
		}
		availableCommands = append(availableCommands, iamgodCmd)
	}

	if len(availableCommands) == 0 {
		ch.sendPlayerMessage(playerName, "No commands available")
		return nil
	}

	perPage := 2
	totalPages := (len(availableCommands) + perPage - 1) / perPage
	if page > totalPages {
		page = totalPages
	}

	startIdx := (page - 1) * perPage
	endIdx := startIdx + perPage
	if endIdx > len(availableCommands) {
		endIdx = len(availableCommands)
	}

	ch.sendPlayerMessage(playerName, fmt.Sprintf("^3Available Commands (Page %d/%d):", page, totalPages))

	for i := startIdx; i < endIdx; i++ {
		cmd := availableCommands[i]
		ch.sendPlayerMessage(playerName, fmt.Sprintf("^2!%s ^7- %s", cmd.Name, cmd.Usage))
		if cmd.Description != "" {
			ch.sendPlayerMessage(playerName, fmt.Sprintf("  ^7%s", cmd.Description))
		}
	}

	if totalPages > 1 {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("^7Use ^3!help <page>^7 to view other pages (1-%d)", totalPages))
	}

	return nil
}
