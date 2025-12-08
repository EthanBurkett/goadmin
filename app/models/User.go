package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	Approved  bool           `gorm:"default:false" json:"approved"`
	Sessions  []Session      `gorm:"foreignKey:UserID" json:"-"`
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func CreateUser(username, password string) (*User, error) {
	user := &User{
		Username: username,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	result := database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByID(id uint) (*User, error) {
	var user User
	result := database.DB.Preload("Roles.Permissions").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *User) HasPermission(permissionName string) bool {
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				return true
			}
		}
	}
	return false
}

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

func AddRoleToUser(userID, roleID uint) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}
	role, err := GetRoleByID(roleID)
	if err != nil {
		return err
	}
	return database.DB.Model(user).Association("Roles").Append(role)
}

func RemoveRoleFromUser(userID, roleID uint) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}
	role, err := GetRoleByID(roleID)
	if err != nil {
		return err
	}
	return database.DB.Model(user).Association("Roles").Delete(role)
}

func GetAllUsers() ([]User, error) {
	var users []User
	result := database.DB.Preload("Roles").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func GetPendingUsers() ([]User, error) {
	var users []User
	result := database.DB.Where("approved = ?", false).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func ApproveUser(userID uint, roleID uint) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	user.Approved = true
	if err := database.DB.Save(user).Error; err != nil {
		return err
	}

	// Assign role
	return AddRoleToUser(userID, roleID)
}

func DenyUser(userID uint) error {
	return database.DB.Delete(&User{}, userID).Error
}

func DeleteUser(userID uint) error {
	// Delete all sessions for this user first
	database.DB.Where("user_id = ?", userID).Delete(&Session{})
	// Delete the user
	return database.DB.Delete(&User{}, userID).Error
}
