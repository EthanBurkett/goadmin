package commands

import (
	"fmt"

	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
)

// processIamGod handles the special !iamgod command
func (ch *CommandHandler) processIamGod(playerName, playerGUID string) error {
	if models.HasSetting("ingame_iamgod_used", "true") {
		ch.sendPlayerMessage(playerName, "This command has already been used and cannot be used again")
		return nil
	}

	groups, err := models.GetAllGroups()
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to fetch groups")
		return err
	}

	var ownerGroup *models.Group
	for _, g := range groups {
		if g.Name == "Owner" {
			ownerGroup = &g
			break
		}
	}

	if ownerGroup == nil {
		ch.sendPlayerMessage(playerName, "Owner group not found")
		return fmt.Errorf("owner group not found")
	}

	player, err := models.CreateOrUpdateInGamePlayer(playerGUID, playerName)
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to create player")
		return err
	}

	if err := models.AssignPlayerToGroup(player.ID, ownerGroup.ID); err != nil {
		ch.sendPlayerMessage(playerName, "Failed to assign Owner group")
		return err
	}

	err = models.SetSetting("ingame_iamgod_used", "true")
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to mark iamgod as used: %v", err))
		return err
	}

	ch.sendPlayerMessage(playerName, "^2Owner privileges granted! This command has been permanently disabled.")
	logger.Info(fmt.Sprintf("Player %s (%s) used !iamgod and was granted Owner privileges", playerName, playerGUID))

	return nil
}
