package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
)

// CustomCommand represents a user-defined chat command
type CustomCommand struct {
	ID              uint         `gorm:"primaryKey" json:"id"`
	Name            string       `gorm:"unique;not null" json:"name"`                                                             // Command trigger (without !)
	Usage           string       `json:"usage"`                                                                                   // How to use: !xp <player>
	Description     string       `json:"description"`                                                                             // What the command does
	RconCommand     string       `json:"rconCommand"`                                                                             // Template: adm xp:{arg0}
	MinArgs         int          `json:"minArgs"`                                                                                 // Minimum required arguments
	MaxArgs         int          `json:"maxArgs"`                                                                                 // Maximum allowed arguments (-1 for unlimited)
	MinPower        int          `json:"minPower"`                                                                                // Minimum power level required (0-100)
	Permissions     []Permission `gorm:"many2many:command_permissions;constraint:OnDelete:CASCADE;" json:"permissions,omitempty"` // Required permissions
	RequirementType string       `gorm:"default:'both'" json:"requirementType"`                                                   // "permission", "power", or "both"
	IsBuiltIn       bool         `gorm:"default:false" json:"isBuiltIn"`                                                          // Whether command is built-in (non-editable)
	Enabled         bool         `gorm:"default:true" json:"enabled"`                                                             // Whether command is active
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
}

// CreateCustomCommand creates a new custom command
func CreateCustomCommand(name, usage, description, rconCommand, requirementType string, minArgs, maxArgs, minPower int, isBuiltIn bool, permissionIDs []uint) error {
	db := database.DB

	// Default to "both" if not specified
	if requirementType == "" {
		requirementType = "both"
	}

	cmd := &CustomCommand{
		Name:            name,
		Usage:           usage,
		Description:     description,
		RconCommand:     rconCommand,
		MinArgs:         minArgs,
		MaxArgs:         maxArgs,
		MinPower:        minPower,
		RequirementType: requirementType,
		IsBuiltIn:       isBuiltIn,
		Enabled:         true,
	}

	err := db.Create(cmd).Error
	if err != nil {
		return err
	}

	// Associate permissions if provided
	if len(permissionIDs) > 0 {
		var permissions []Permission
		if err := db.Find(&permissions, permissionIDs).Error; err != nil {
			return err
		}
		if err := db.Model(cmd).Association("Permissions").Replace(&permissions); err != nil {
			return err
		}
	}

	return nil
}

// GetCustomCommand gets a command by name
func GetCustomCommand(name string) (*CustomCommand, error) {
	db := database.DB
	var cmd CustomCommand
	err := db.Preload("Permissions").Where("name = ? AND enabled = ?", name, true).First(&cmd).Error
	if err != nil {
		return nil, err
	}
	return &cmd, nil
}

// GetCustomCommandByID gets a command by ID
func GetCustomCommandByID(id uint) (*CustomCommand, error) {
	db := database.DB
	var cmd CustomCommand
	err := db.Preload("Permissions").First(&cmd, id).Error
	if err != nil {
		return nil, err
	}
	return &cmd, nil
}

// GetAllCustomCommands gets all custom commands
func GetAllCustomCommands() ([]CustomCommand, error) {
	db := database.DB
	var commands []CustomCommand
	err := db.Preload("Permissions").Order("name").Find(&commands).Error
	return commands, err
}

// UpdateCustomCommand updates an existing command
func UpdateCustomCommand(id uint, updates map[string]interface{}) error {
	db := database.DB
	return db.Model(&CustomCommand{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCustomCommand deletes a command
func DeleteCustomCommand(id uint) error {
	db := database.DB
	return db.Delete(&CustomCommand{}, id).Error
}

// AddPermissionToCommand adds a permission to a custom command
func (c *CustomCommand) AddPermissionToCommand(permissionID uint) error {
	db := database.DB

	permission, err := GetPermissionByID(permissionID)
	if err != nil {
		return err
	}

	return db.Model(c).Association("Permissions").Append(permission)
}

// RemovePermissionFromCommand removes a permission from a custom command
func (c *CustomCommand) RemovePermissionFromCommand(permissionID uint) error {
	db := database.DB

	permission, err := GetPermissionByID(permissionID)
	if err != nil {
		return err
	}

	return db.Model(c).Association("Permissions").Delete(permission)
}

// SetCommandPermissions replaces all permissions for a command
func (c *CustomCommand) SetCommandPermissions(permissionIDs []uint) error {
	db := database.DB

	if len(permissionIDs) == 0 {
		return db.Model(c).Association("Permissions").Clear()
	}

	var permissions []Permission
	if err := db.Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return db.Model(c).Association("Permissions").Replace(&permissions)
}

// HasPermission checks if a command requires a specific permission
func (c *CustomCommand) HasPermission(permissionName string) bool {
	for _, perm := range c.Permissions {
		if perm.Name == permissionName {
			return true
		}
	}
	return false
}

// TableName specifies the table name
func (CustomCommand) TableName() string {
	return "custom_commands"
}
