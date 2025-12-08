package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
)

// Group represents an in-game admin/user group with power level (like B3)
type Group struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`  // e.g., "SuperAdmin", "Admin", "Moderator"
	Power       int       `gorm:"not null" json:"power"`        // Power level 0-100 (100 = highest)
	Permissions string    `gorm:"type:text" json:"permissions"` // JSON array of permission strings
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// InGamePlayer represents a player identified by their PB GUID (separate from web auth)
type InGamePlayer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GUID      string    `gorm:"unique;not null;index" json:"guid"` // PB GUID/UUID/XUID
	Name      string    `gorm:"index" json:"name"`                 // Last known name
	GroupID   *uint     `json:"groupId"`                           // Optional group assignment
	Group     *Group    `gorm:"foreignKey:GroupID;constraint:OnDelete:SET NULL" json:"group,omitempty"`
	Enabled   bool      `gorm:"default:true" json:"enabled"` // Can be disabled/banned
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateGroup creates a new in-game group
func CreateGroup(name string, power int, permissions, description string) error {
	db := database.DB
	group := &Group{
		Name:        name,
		Power:       power,
		Permissions: permissions,
		Description: description,
	}
	return db.Create(group).Error
}

// GetAllGroups gets all groups ordered by power level (highest first)
func GetAllGroups() ([]Group, error) {
	db := database.DB
	var groups []Group
	err := db.Order("power DESC").Find(&groups).Error
	return groups, err
}

// GetGroupByID gets a group by ID
func GetGroupByID(id uint) (*Group, error) {
	db := database.DB
	var group Group
	err := db.First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// UpdateGroup updates a group
func UpdateGroup(id uint, updates map[string]interface{}) error {
	db := database.DB
	return db.Model(&Group{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteGroup deletes a group
func DeleteGroup(id uint) error {
	db := database.DB
	// Remove group assignment from all players
	db.Model(&InGamePlayer{}).Where("group_id = ?", id).Update("group_id", nil)
	return db.Delete(&Group{}, id).Error
}

// CreateOrUpdateInGamePlayer creates or updates an in-game player by GUID
func CreateOrUpdateInGamePlayer(guid, name string) (*InGamePlayer, error) {
	db := database.DB
	var player InGamePlayer

	err := db.Where("guid = ?", guid).First(&player).Error
	if err == gorm.ErrRecordNotFound {
		// Create new player
		player = InGamePlayer{
			GUID:    guid,
			Name:    name,
			Enabled: true,
		}
		if err := db.Create(&player).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// Update name if changed
		if player.Name != name {
			db.Model(&player).Update("name", name)
		}
	}

	// Load group if assigned
	if player.GroupID != nil {
		db.Preload("Group").First(&player, player.ID)
	}

	return &player, nil
}

// GetInGamePlayerByGUID gets a player by their GUID
func GetInGamePlayerByGUID(guid string) (*InGamePlayer, error) {
	db := database.DB
	var player InGamePlayer
	err := db.Preload("Group").Where("guid = ?", guid).First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// GetAllInGamePlayers gets all in-game players
func GetAllInGamePlayers() ([]InGamePlayer, error) {
	db := database.DB
	var players []InGamePlayer
	err := db.Preload("Group").Order("name").Find(&players).Error
	return players, err
}

// AssignPlayerToGroup assigns a player to a group
func AssignPlayerToGroup(playerID, groupID uint) error {
	db := database.DB
	return db.Model(&InGamePlayer{}).Where("id = ?", playerID).Update("group_id", groupID).Error
}

// RemovePlayerFromGroup removes a player from their group
func RemovePlayerFromGroup(playerID uint) error {
	db := database.DB
	return db.Model(&InGamePlayer{}).Where("id = ?", playerID).Update("group_id", nil).Error
}

// GetPlayerPower gets the effective power level of a player (0 if no group)
func GetPlayerPower(guid string) int {
	player, err := GetInGamePlayerByGUID(guid)
	if err != nil || player.Group == nil || !player.Enabled {
		return 0
	}
	return player.Group.Power
}

// TableName specifies the table name for Group
func (Group) TableName() string {
	return "groups"
}

// TableName specifies the table name for InGamePlayer
func (InGamePlayer) TableName() string {
	return "in_game_players"
}
