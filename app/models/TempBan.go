package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
)

// TempBan represents a temporary ban
type TempBan struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PlayerName   string    `gorm:"not null" json:"playerName"`       // Name of banned player
	PlayerGUID   string    `gorm:"not null;index" json:"playerGuid"` // GUID of banned player
	Reason       string    `gorm:"type:text;not null" json:"reason"` // Reason for ban
	BannedByUser *uint     `json:"bannedByUser"`                     // User who issued ban
	BannedBy     *User     `gorm:"foreignKey:BannedByUser;constraint:OnDelete:SET NULL" json:"bannedBy,omitempty"`
	ExpiresAt    time.Time `gorm:"not null;index" json:"expiresAt"`  // When ban expires
	Active       bool      `gorm:"default:true;index" json:"active"` // Whether ban is still active
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// CreateTempBan creates a new temporary ban
func CreateTempBan(playerName, playerGUID, reason string, duration time.Duration, bannedByUser *uint) (*TempBan, error) {
	db := database.DB

	tempBan := &TempBan{
		PlayerName:   playerName,
		PlayerGUID:   playerGUID,
		Reason:       reason,
		BannedByUser: bannedByUser,
		ExpiresAt:    time.Now().Add(duration),
		Active:       true,
	}

	err := db.Create(tempBan).Error
	return tempBan, err
}

// GetActiveTempBans gets all active temporary bans
func GetActiveTempBans() ([]TempBan, error) {
	db := database.DB
	var bans []TempBan
	err := db.Preload("BannedBy").Where("active = ? AND expires_at > ?", true, time.Now()).Order("created_at DESC").Find(&bans).Error
	return bans, err
}

// GetAllTempBans gets all temporary bans (active and expired)
func GetAllTempBans() ([]TempBan, error) {
	db := database.DB
	var bans []TempBan
	err := db.Preload("BannedBy").Order("created_at DESC").Find(&bans).Error
	return bans, err
}

// GetTempBanByGUID checks if a player is temp banned
func GetTempBanByGUID(guid string) (*TempBan, error) {
	db := database.DB
	var ban TempBan
	err := db.Where("player_guid = ? AND active = ? AND expires_at > ?", guid, true, time.Now()).First(&ban).Error
	if err != nil {
		return nil, err
	}
	return &ban, nil
}

// IsPlayerTempBanned checks if a player is currently temp banned
func IsPlayerTempBanned(guid string) bool {
	ban, err := GetTempBanByGUID(guid)
	return err == nil && ban != nil
}

// ExpireTempBans marks expired bans as inactive
func ExpireTempBans() error {
	db := database.DB
	return db.Model(&TempBan{}).Where("active = ? AND expires_at <= ?", true, time.Now()).Update("active", false).Error
}

// RevokeTempBan manually revokes a temporary ban
func RevokeTempBan(id uint) error {
	db := database.DB
	return db.Model(&TempBan{}).Where("id = ?", id).Update("active", false).Error
}

// TableName specifies the table name
func (TempBan) TableName() string {
	return "temp_bans"
}
