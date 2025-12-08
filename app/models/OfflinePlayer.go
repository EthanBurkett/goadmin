package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OfflinePlayer struct {
	ID            uint           `json:"id"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	PlayerID      string         `gorm:"uniqueIndex;not null" json:"playerId"`
	PlayerSteamID string         `json:"playerSteamId"`
	Name          string         `json:"name"`
	IP            string         `json:"ip"`
	PBGuid        string         `json:"pbGuid"`
}

type OptionalOfflinePlayer struct {
	IP     *string
	PBGuid *string
}

func GetOfflinePlayer(identifier string) *OfflinePlayer {
	var player OfflinePlayer
	result := database.DB.Where("player_id = ? OR pb_guid = ?", identifier, identifier).First(&player)
	if result.Error != nil {
		return nil
	}
	return &player
}

func UpdateOfflinePlayer(player *OfflinePlayer, optional *OptionalOfflinePlayer) error {
	updates := map[string]interface{}{
		"player_steam_id": player.PlayerSteamID,
		"name":            player.Name,
	}

	if optional != nil {
		if optional.IP != nil {
			updates["ip"] = *optional.IP
		}
		if optional.PBGuid != nil {
			updates["pb_guid"] = *optional.PBGuid
		}
	}

	// Upsert: Insert or update based on player_id uniqueness
	result := database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "player_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"player_steam_id", "name", "ip", "pb_guid"}),
	}).Create(player)

	return result.Error
}
