package models

import (
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration record
type Migration struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Version      string     `gorm:"uniqueIndex;not null" json:"version"`
	Name         string     `gorm:"not null" json:"name"`
	Description  string     `json:"description"`
	AppliedAt    time.Time  `gorm:"autoCreateTime" json:"applied_at"`
	RolledBack   bool       `gorm:"default:false" json:"rolled_back"`
	RolledBackAt *time.Time `json:"rolled_back_at,omitempty"`
}

// MigrationHistory tracks all migration executions
type MigrationHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MigrationID uint      `json:"migration_id"`
	Action      string    `gorm:"not null" json:"action"` // "apply" or "rollback"
	ExecutedAt  time.Time `gorm:"autoCreateTime" json:"executed_at"`
	Success     bool      `gorm:"default:true" json:"success"`
	Error       string    `json:"error,omitempty"`
	Duration    int64     `json:"duration"` // Duration in milliseconds
}

// GetCurrentMigrationVersion returns the latest applied migration version
func GetCurrentMigrationVersion(db *gorm.DB) (string, error) {
	var migration Migration
	err := db.Where("rolled_back = ?", false).
		Order("applied_at DESC").
		First(&migration).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "none", nil
		}
		return "", err
	}

	return migration.Version, nil
}

// GetAllMigrations returns all migrations in order
func GetAllMigrations(db *gorm.DB) ([]Migration, error) {
	var migrations []Migration
	err := db.Order("version ASC").Find(&migrations).Error
	return migrations, err
}

// GetPendingMigrations returns migrations that haven't been applied
func GetPendingMigrations(db *gorm.DB, knownVersions []string) []string {
	var appliedVersions []string
	db.Model(&Migration{}).
		Where("rolled_back = ?", false).
		Pluck("version", &appliedVersions)

	appliedMap := make(map[string]bool)
	for _, v := range appliedVersions {
		appliedMap[v] = true
	}

	var pending []string
	for _, v := range knownVersions {
		if !appliedMap[v] {
			pending = append(pending, v)
		}
	}

	return pending
}

// RecordMigration records a successful migration
func RecordMigration(db *gorm.DB, version, name, description string) error {
	migration := &Migration{
		Version:     version,
		Name:        name,
		Description: description,
		AppliedAt:   time.Now(),
		RolledBack:  false,
	}

	return db.Create(migration).Error
}

// RecordMigrationRollback marks a migration as rolled back
func RecordMigrationRollback(db *gorm.DB, version string) error {
	now := time.Now()
	return db.Model(&Migration{}).
		Where("version = ?", version).
		Updates(map[string]interface{}{
			"rolled_back":    true,
			"rolled_back_at": now,
		}).Error
}

// RecordMigrationHistory records a migration execution attempt
func RecordMigrationHistory(db *gorm.DB, migrationID uint, action string, success bool, errorMsg string, duration int64) error {
	history := &MigrationHistory{
		MigrationID: migrationID,
		Action:      action,
		Success:     success,
		Error:       errorMsg,
		Duration:    duration,
	}

	return db.Create(history).Error
}

// GetMigrationHistory returns the history for a specific migration
func GetMigrationHistory(db *gorm.DB, migrationID uint) ([]MigrationHistory, error) {
	var history []MigrationHistory
	err := db.Where("migration_id = ?", migrationID).
		Order("executed_at DESC").
		Find(&history).Error
	return history, err
}

// CheckMigrationIntegrity validates that all expected migrations are present
func CheckMigrationIntegrity(db *gorm.DB, expectedVersions []string) (bool, []string, error) {
	var appliedVersions []string
	err := db.Model(&Migration{}).
		Where("rolled_back = ?", false).
		Pluck("version", &appliedVersions).Error

	if err != nil {
		return false, nil, err
	}

	appliedMap := make(map[string]bool)
	for _, v := range appliedVersions {
		appliedMap[v] = true
	}

	var missing []string
	for _, v := range expectedVersions {
		if !appliedMap[v] {
			missing = append(missing, v)
		}
	}

	return len(missing) == 0, missing, nil
}
