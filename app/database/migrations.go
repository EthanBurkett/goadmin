package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// MigrationFunc represents a migration function
type MigrationFunc func(*gorm.DB) error

// MigrationDefinition represents a single migration
type MigrationDefinition struct {
	Version     string
	Name        string
	Description string
	Up          MigrationFunc
	Down        MigrationFunc
}

// MigrationRunner manages database migrations
type MigrationRunner struct {
	db         *gorm.DB
	migrations []MigrationDefinition
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *gorm.DB) *MigrationRunner {
	return &MigrationRunner{
		db:         db,
		migrations: []MigrationDefinition{},
	}
}

// Register adds a migration to the runner
func (mr *MigrationRunner) Register(migration MigrationDefinition) {
	mr.migrations = append(mr.migrations, migration)
}

// ApplyMigration applies a single migration
func (mr *MigrationRunner) ApplyMigration(migration MigrationDefinition) error {
	startTime := time.Now()

	// Check if migration already applied
	var count int64
	mr.db.Table("migrations").
		Where("version = ? AND rolled_back = ?", migration.Version, false).
		Count(&count)

	if count > 0 {
		log.Printf("Migration %s (%s) already applied, skipping", migration.Version, migration.Name)
		return nil
	}

	// Execute migration in transaction
	err := mr.db.Transaction(func(tx *gorm.DB) error {
		// Run the migration
		if err := migration.Up(tx); err != nil {
			return fmt.Errorf("migration up failed: %w", err)
		}

		// Record migration
		migrationRecord := map[string]interface{}{
			"version":     migration.Version,
			"name":        migration.Name,
			"description": migration.Description,
			"applied_at":  time.Now(),
			"rolled_back": false,
		}

		if err := tx.Table("migrations").Create(migrationRecord).Error; err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		return nil
	})

	duration := time.Since(startTime).Milliseconds()

	// Record history
	if err != nil {
		log.Printf("Migration %s (%s) failed: %v (duration: %dms)",
			migration.Version, migration.Name, err, duration)
		return err
	}

	log.Printf("Migration %s (%s) applied successfully (duration: %dms)",
		migration.Version, migration.Name, duration)
	return nil
}

// RollbackMigration rolls back a single migration
func (mr *MigrationRunner) RollbackMigration(migration MigrationDefinition) error {
	startTime := time.Now()

	// Check if migration was applied
	var count int64
	mr.db.Table("migrations").
		Where("version = ? AND rolled_back = ?", migration.Version, false).
		Count(&count)

	if count == 0 {
		return fmt.Errorf("migration %s was not applied or already rolled back", migration.Version)
	}

	// Execute rollback in transaction
	err := mr.db.Transaction(func(tx *gorm.DB) error {
		// Run the rollback
		if err := migration.Down(tx); err != nil {
			return fmt.Errorf("migration down failed: %w", err)
		}

		// Mark migration as rolled back
		now := time.Now()
		if err := tx.Table("migrations").
			Where("version = ?", migration.Version).
			Updates(map[string]interface{}{
				"rolled_back":    true,
				"rolled_back_at": now,
			}).Error; err != nil {
			return fmt.Errorf("failed to record rollback: %w", err)
		}

		return nil
	})

	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		log.Printf("Migration %s rollback failed: %v (duration: %dms)",
			migration.Version, err, duration)
		return err
	}

	log.Printf("Migration %s rolled back successfully (duration: %dms)",
		migration.Version, duration)
	return nil
}

// ApplyAll applies all pending migrations
func (mr *MigrationRunner) ApplyAll() error {
	for _, migration := range mr.migrations {
		if err := mr.ApplyMigration(migration); err != nil {
			return err
		}
	}
	return nil
}

// RollbackLast rolls back the most recent migration
func (mr *MigrationRunner) RollbackLast() error {
	var lastMigration struct {
		Version string
	}

	err := mr.db.Table("migrations").
		Where("rolled_back = ?", false).
		Order("applied_at DESC").
		First(&lastMigration).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no migrations to rollback")
		}
		return err
	}

	// Find the migration definition
	for _, migration := range mr.migrations {
		if migration.Version == lastMigration.Version {
			return mr.RollbackMigration(migration)
		}
	}

	return fmt.Errorf("migration definition not found for version %s", lastMigration.Version)
}

// GetStatus returns the current migration status
func (mr *MigrationRunner) GetStatus() (map[string]interface{}, error) {
	var applied []struct {
		Version    string
		Name       string
		AppliedAt  time.Time
		RolledBack bool
	}

	err := mr.db.Table("migrations").
		Order("version ASC").
		Find(&applied).Error

	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]bool)
	for _, m := range applied {
		if !m.RolledBack {
			appliedMap[m.Version] = true
		}
	}

	var pending []string
	for _, migration := range mr.migrations {
		if !appliedMap[migration.Version] {
			pending = append(pending, migration.Version)
		}
	}

	return map[string]interface{}{
		"total_migrations":   len(mr.migrations),
		"applied_migrations": len(appliedMap),
		"pending_migrations": len(pending),
		"pending_versions":   pending,
		"applied_versions":   applied,
	}, nil
}
