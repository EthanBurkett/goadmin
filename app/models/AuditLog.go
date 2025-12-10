package models

import (
	"time"

	"gorm.io/gorm"
)

// ActionType represents the type of action performed
type ActionType string

const (
	ActionBanPlayer         ActionType = "ban_player"
	ActionTempBanPlayer     ActionType = "tempban_player"
	ActionKickPlayer        ActionType = "kick_player"
	ActionUnbanPlayer       ActionType = "unban_player"
	ActionRconCommand       ActionType = "rcon_command"
	ActionRoleAssign        ActionType = "role_assign"
	ActionRoleRevoke        ActionType = "role_revoke"
	ActionGroupAssign       ActionType = "group_assign"
	ActionPermissionGrant   ActionType = "permission_grant"
	ActionPermissionRevoke  ActionType = "permission_revoke"
	ActionUserApprove       ActionType = "user_approve"
	ActionUserReject        ActionType = "user_reject"
	ActionCommandCreate     ActionType = "command_create"
	ActionCommandUpdate     ActionType = "command_update"
	ActionCommandDelete     ActionType = "command_delete"
	ActionReportReview      ActionType = "report_review"
	ActionReportDismiss     ActionType = "report_dismiss"
	ActionLogin             ActionType = "login"
	ActionLogout            ActionType = "logout"
	ActionLoginFailed       ActionType = "login_failed"
	ActionSecurityViolation ActionType = "security_violation"
	ActionSystemChange      ActionType = "system_change"
)

// ActionSource represents where the action was initiated
type ActionSource string

const (
	SourceWebUI  ActionSource = "web_ui"
	SourceInGame ActionSource = "in_game"
	SourceAPI    ActionSource = "api"
	SourceSystem ActionSource = "system"
)

// AuditLog represents a comprehensive audit trail entry
type AuditLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// Actor information
	UserID    *uint  `gorm:"index" json:"userId"`
	User      *User  `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	Username  string `gorm:"index" json:"username"` // Denormalized for when user is deleted
	IPAddress string `gorm:"index" json:"ipAddress"`

	// Action details
	Action       ActionType   `gorm:"type:varchar(50);not null;index" json:"action"`
	Source       ActionSource `gorm:"type:varchar(20);not null;index" json:"source"`
	Success      bool         `gorm:"index" json:"success"`
	ErrorMessage string       `gorm:"type:text" json:"errorMessage,omitempty"`

	// Target information
	TargetType string `gorm:"type:varchar(50);index" json:"targetType,omitempty"` // e.g., "user", "player", "command", "report"
	TargetID   string `gorm:"index" json:"targetId,omitempty"`                    // ID of the target entity
	TargetName string `json:"targetName,omitempty"`                               // Name for display

	// Additional context (stored as JSON)
	Metadata string `gorm:"type:text" json:"metadata,omitempty"` // JSON encoded additional data

	// Result data
	Result string `gorm:"type:text" json:"result,omitempty"` // Response or outcome
}

// CreateAuditLog creates a new audit log entry
func CreateAuditLog(
	db *gorm.DB,
	userID *uint,
	username string,
	ipAddress string,
	action ActionType,
	source ActionSource,
	success bool,
	errorMessage string,
	targetType string,
	targetID string,
	targetName string,
	metadata string,
	result string,
) (*AuditLog, error) {
	log := &AuditLog{
		UserID:       userID,
		Username:     username,
		IPAddress:    ipAddress,
		Action:       action,
		Source:       source,
		Success:      success,
		ErrorMessage: errorMessage,
		TargetType:   targetType,
		TargetID:     targetID,
		TargetName:   targetName,
		Metadata:     metadata,
		Result:       result,
	}

	err := db.Create(log).Error
	return log, err
}

// GetAuditLogs retrieves audit logs with optional filters
func GetAuditLogs(db *gorm.DB, filters map[string]interface{}, limit int, offset int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := db.Model(&AuditLog{}).Preload("User")

	// Apply filters
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if action, ok := filters["action"]; ok {
		query = query.Where("action = ?", action)
	}
	if source, ok := filters["source"]; ok {
		query = query.Where("source = ?", source)
	}
	if success, ok := filters["success"]; ok {
		query = query.Where("success = ?", success)
	}
	if targetType, ok := filters["target_type"]; ok {
		query = query.Where("target_type = ?", targetType)
	}
	if targetID, ok := filters["target_id"]; ok {
		query = query.Where("target_id = ?", targetID)
	}
	if startDate, ok := filters["start_date"]; ok {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"]; ok {
		query = query.Where("created_at <= ?", endDate)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

// GetAuditLogsByUser retrieves audit logs for a specific user
func GetAuditLogsByUser(db *gorm.DB, userID uint, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetAuditLogsByAction retrieves audit logs for a specific action type
func GetAuditLogsByAction(db *gorm.DB, action ActionType, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := db.Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetAuditLogsByTarget retrieves audit logs for a specific target
func GetAuditLogsByTarget(db *gorm.DB, targetType string, targetID string, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := db.Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetRecentAuditLogs retrieves the most recent audit logs
func GetRecentAuditLogs(db *gorm.DB, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := db.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// ArchiveOldAuditLogs moves audit logs older than the retention period to archive
// Returns the number of logs archived and any error
func ArchiveOldAuditLogs(db *gorm.DB, retentionDays int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	// Soft delete old logs (they'll remain in DB but be excluded from queries)
	result := db.Where("created_at < ?", cutoffDate).Delete(&AuditLog{})

	return result.RowsAffected, result.Error
}

// PurgeArchivedAuditLogs permanently deletes archived (soft-deleted) logs
// Returns the number of logs purged and any error
func PurgeArchivedAuditLogs(db *gorm.DB) (int64, error) {
	// Permanently delete soft-deleted logs
	result := db.Unscoped().Where("deleted_at IS NOT NULL").Delete(&AuditLog{})

	return result.RowsAffected, result.Error
}

// GetAuditLogStats returns statistics about audit logs
func GetAuditLogStats(db *gorm.DB) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total logs
	var total int64
	if err := db.Model(&AuditLog{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// Archived logs (soft-deleted)
	var archived int64
	if err := db.Unscoped().Model(&AuditLog{}).Where("deleted_at IS NOT NULL").Count(&archived).Error; err != nil {
		return nil, err
	}
	stats["archived"] = archived

	// Logs by action type
	var actionCounts []struct {
		Action ActionType
		Count  int64
	}
	if err := db.Model(&AuditLog{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Scan(&actionCounts).Error; err != nil {
		return nil, err
	}
	stats["by_action"] = actionCounts

	// Logs by source
	var sourceCounts []struct {
		Source ActionSource
		Count  int64
	}
	if err := db.Model(&AuditLog{}).
		Select("source, COUNT(*) as count").
		Group("source").
		Scan(&sourceCounts).Error; err != nil {
		return nil, err
	}
	stats["by_source"] = sourceCounts

	// Success rate
	var successCount int64
	if err := db.Model(&AuditLog{}).Where("success = ?", true).Count(&successCount).Error; err != nil {
		return nil, err
	}
	if total > 0 {
		stats["success_rate"] = float64(successCount) / float64(total) * 100
	} else {
		stats["success_rate"] = 0.0
	}

	// Oldest and newest log dates
	var oldestStr, newestStr string
	db.Model(&AuditLog{}).Select("MIN(created_at)").Row().Scan(&oldestStr)
	db.Model(&AuditLog{}).Select("MAX(created_at)").Row().Scan(&newestStr)

	if oldestStr != "" {
		if oldest, err := time.Parse(time.RFC3339, oldestStr); err == nil {
			stats["oldest_log"] = oldest
		} else if oldest, err := time.Parse("2006-01-02 15:04:05", oldestStr); err == nil {
			stats["oldest_log"] = oldest
		}
	}

	if newestStr != "" {
		if newest, err := time.Parse(time.RFC3339, newestStr); err == nil {
			stats["newest_log"] = newest
		} else if newest, err := time.Parse("2006-01-02 15:04:05", newestStr); err == nil {
			stats["newest_log"] = newest
		}
	}

	return stats, nil
}
