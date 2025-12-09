package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	Users       []User         `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;" json:"users,omitempty"`
	Permissions []Permission   `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE;" json:"permissions,omitempty"`
}

type Permission struct {
	ID          uint           `json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	Roles       []Role         `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE;" json:"roles,omitempty"`
}

func CreateRole(name, description string) (*Role, error) {
	role := &Role{
		Name:        name,
		Description: description,
	}
	result := database.DB.Create(role)
	if result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}

func GetRoleByID(id uint) (*Role, error) {
	var role Role
	result := database.DB.Preload("Permissions").First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func GetRoleByName(name string) (*Role, error) {
	var role Role
	result := database.DB.Preload("Permissions").Where("name = ?", name).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func GetAllRoles() ([]Role, error) {
	var roles []Role
	result := database.DB.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func DeleteRole(id uint) error {
	result := database.DB.Delete(&Role{}, id)
	return result.Error
}

func CreatePermission(name, description string) (*Permission, error) {
	permission := &Permission{
		Name:        name,
		Description: description,
	}
	result := database.DB.Create(permission)
	if result.Error != nil {
		return nil, result.Error
	}
	return permission, nil
}

func GetPermissionByID(id uint) (*Permission, error) {
	var permission Permission
	result := database.DB.First(&permission, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &permission, nil
}

func GetPermissionByName(name string) (*Permission, error) {
	var permission Permission
	result := database.DB.Where("name = ?", name).First(&permission)
	if result.Error != nil {
		return nil, result.Error
	}
	return &permission, nil
}

func GetAllPermissions() ([]Permission, error) {
	var permissions []Permission
	result := database.DB.Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

func DeletePermission(id uint) error {
	result := database.DB.Delete(&Permission{}, id)
	return result.Error
}

func AddPermissionToRole(roleID, permissionID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.Preload("Permissions").First(&role, roleID).Error; err != nil {
			return err
		}
		var permission Permission
		if err := tx.First(&permission, permissionID).Error; err != nil {
			return err
		}
		return tx.Model(&role).Association("Permissions").Append(&permission)
	})
}

func RemovePermissionFromRole(roleID, permissionID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.Preload("Permissions").First(&role, roleID).Error; err != nil {
			return err
		}
		var permission Permission
		if err := tx.First(&permission, permissionID).Error; err != nil {
			return err
		}
		return tx.Model(&role).Association("Permissions").Delete(&permission)
	})
}
