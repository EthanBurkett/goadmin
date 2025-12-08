package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
)

type Session struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Token     string         `gorm:"uniqueIndex;not null" json:"token"`
	UserID    uint           `gorm:"not null;index" json:"userId"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	ExpiresAt time.Time      `gorm:"not null" json:"expiresAt"`
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(userID uint) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	result := database.DB.Create(session)
	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}

func GetSessionByToken(token string) (*Session, error) {
	var session Session
	result := database.DB.Preload("User").Where("token = ? AND expires_at > ?", token, time.Now()).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func DeleteSession(token string) error {
	result := database.DB.Where("token = ?", token).Delete(&Session{})
	return result.Error
}

func DeleteExpiredSessions() error {
	result := database.DB.Where("expires_at < ?", time.Now()).Delete(&Session{})
	return result.Error
}
