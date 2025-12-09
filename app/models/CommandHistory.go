package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
)

type CommandHistory struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	UserID    uint           `gorm:"not null;index" json:"userId"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	ServerID  *uint          `gorm:"index" json:"serverId,omitempty"`
	Server    *Server        `gorm:"foreignKey:ServerID;constraint:OnDelete:SET NULL" json:"server,omitempty"`
	Command   string         `gorm:"not null" json:"command"`
	Response  string         `gorm:"type:text" json:"response"`
	Success   bool           `gorm:"default:true" json:"success"`
}

func CreateCommandHistory(userID uint, command string, response string, success bool, serverID *uint) (*CommandHistory, error) {
	history := &CommandHistory{
		UserID:   userID,
		Command:  command,
		Response: response,
		Success:  success,
		ServerID: serverID,
	}

	result := database.DB.Create(history)
	if result.Error != nil {
		return nil, result.Error
	}

	return history, nil
}

func GetCommandHistoryByUser(userID uint, limit int, serverID *uint) ([]CommandHistory, error) {
	var history []CommandHistory
	query := database.DB.Preload("Server").Where("user_id = ?", userID)

	if serverID != nil {
		query = query.Where("server_id = ?", *serverID)
	}

	query = query.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&history)
	if result.Error != nil {
		return nil, result.Error
	}

	return history, nil
}

func GetAllCommandHistory(limit int, serverID *uint) ([]CommandHistory, error) {
	var history []CommandHistory
	query := database.DB.Preload("User").Preload("Server")

	if serverID != nil {
		query = query.Where("server_id = ?", *serverID)
	}

	query = query.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&history)
	if result.Error != nil {
		return nil, result.Error
	}

	return history, nil
}

func DeleteCommandHistory(id uint) error {
	return database.DB.Delete(&CommandHistory{}, id).Error
}
