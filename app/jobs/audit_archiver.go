package jobs

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
	"go.uber.org/zap"
)

// AuditLogArchiver handles automatic archiving of old audit logs
type AuditLogArchiver struct {
	retentionDays int
	stopChan      chan bool
}

// NewAuditLogArchiver creates a new audit log archiver
func NewAuditLogArchiver(retentionDays int) *AuditLogArchiver {
	if retentionDays <= 0 {
		retentionDays = 90 // Default to 90 days
	}

	return &AuditLogArchiver{
		retentionDays: retentionDays,
		stopChan:      make(chan bool),
	}
}

// Start begins the archiving process (runs daily at 2 AM)
func (a *AuditLogArchiver) Start() {
	logger.Info("Starting audit log archiver",
		zap.Int("retention_days", a.retentionDays))

	// Run immediately on startup
	a.archive()

	// Then run daily at 2 AM
	go a.schedule()
}

// Stop halts the archiving process
func (a *AuditLogArchiver) Stop() {
	logger.Info("Stopping audit log archiver")
	close(a.stopChan)
}

// schedule runs the archiving task daily at 2 AM
func (a *AuditLogArchiver) schedule() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			// Check if it's 2 AM
			if now.Hour() == 2 {
				a.archive()
			}
		case <-a.stopChan:
			return
		}
	}
}

// archive performs the actual archiving
func (a *AuditLogArchiver) archive() {
	logger.Info("Running audit log archival",
		zap.Int("retention_days", a.retentionDays))

	// Archive old logs
	archived, err := models.ArchiveOldAuditLogs(database.DB, a.retentionDays)
	if err != nil {
		logger.Error("Failed to archive audit logs", zap.Error(err))
		return
	}

	if archived > 0 {
		logger.Info("Archived old audit logs",
			zap.Int64("count", archived),
			zap.Int("retention_days", a.retentionDays))
	}

	// Get archival retention setting (how long to keep archived logs before purging)
	// Default to 1 year for archived log retention
	archiveRetentionDays := 365

	// Optionally purge very old archived logs after extended retention
	// This permanently deletes logs that have been archived for more than archiveRetentionDays
	cutoffDate := time.Now().AddDate(0, 0, -archiveRetentionDays)
	result := database.DB.Unscoped().
		Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoffDate).
		Delete(&models.AuditLog{})

	if result.Error != nil {
		logger.Error("Failed to purge old archived logs", zap.Error(result.Error))
	} else if result.RowsAffected > 0 {
		logger.Info("Purged old archived audit logs",
			zap.Int64("count", result.RowsAffected),
			zap.Int("archive_retention_days", archiveRetentionDays))
	}
}

// SetRetentionDays updates the retention period
func (a *AuditLogArchiver) SetRetentionDays(days int) {
	if days > 0 {
		a.retentionDays = days
		logger.Info("Updated audit log retention period", zap.Int("retention_days", days))
	}
}

// GetStats returns current archiving statistics
func (a *AuditLogArchiver) GetStats() map[string]interface{} {
	stats, err := models.GetAuditLogStats(database.DB)
	if err != nil {
		logger.Error("Failed to get audit log stats", zap.Error(err))
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	stats["retention_days"] = a.retentionDays
	return stats
}
