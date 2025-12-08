package models

import (
	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex;not null"`
	Value string `gorm:"not null"`
}

// GetSetting retrieves a setting by key
func GetSetting(key string) (*Setting, error) {
	var setting Setting
	result := database.DB.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		return nil, result.Error
	}
	return &setting, nil
}

// SetSetting creates or updates a setting
func SetSetting(key, value string) error {
	var setting Setting
	result := database.DB.Where("key = ?", key).First(&setting)

	if result.Error != nil {
		// Setting doesn't exist, create it
		setting = Setting{Key: key, Value: value}
		return database.DB.Create(&setting).Error
	}

	// Setting exists, update it
	setting.Value = value
	return database.DB.Save(&setting).Error
}

// HasSetting checks if a setting exists and equals a specific value
func HasSetting(key, value string) bool {
	setting, err := GetSetting(key)
	if err != nil {
		return false
	}
	return setting.Value == value
}
