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
	Command   string         `gorm:"not null" json:"command"`
	Response  string         `gorm:"type:text" json:"response"`
	Success   bool           `gorm:"default:true" json:"success"`
}

func CreateCommandHistory(userID uint, command string, response string, success bool) (*CommandHistory, error) {
	history := &CommandHistory{
		UserID:   userID,
		Command:  command,
		Response: response,
		Success:  success,
	}

	result := database.DB.Create(history)
	if result.Error != nil {
		return nil, result.Error
	}

	return history, nil
}

func GetCommandHistoryByUser(userID uint, limit int) ([]CommandHistory, error) {
	var history []CommandHistory
	query := database.DB.Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&history)
	if result.Error != nil {
		return nil, result.Error
	}

	return history, nil
}

func GetAllCommandHistory(limit int) ([]CommandHistory, error) {
	var history []CommandHistory
	query := database.DB.Preload("User").Order("created_at DESC")

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
