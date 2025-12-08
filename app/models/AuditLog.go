package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
)

// ActionType represents the type of action performed
type ActionType string

const (
	ActionBanPlayer        ActionType = "ban_player"
	ActionTempBanPlayer    ActionType = "tempban_player"
	ActionKickPlayer       ActionType = "kick_player"
	ActionUnbanPlayer      ActionType = "unban_player"
	ActionRconCommand      ActionType = "rcon_command"
	ActionRoleAssign       ActionType = "role_assign"
	ActionRoleRevoke       ActionType = "role_revoke"
	ActionGroupAssign      ActionType = "group_assign"
	ActionPermissionGrant  ActionType = "permission_grant"
	ActionPermissionRevoke ActionType = "permission_revoke"
	ActionUserApprove      ActionType = "user_approve"
	ActionUserReject       ActionType = "user_reject"
	ActionCommandCreate    ActionType = "command_create"
	ActionCommandUpdate    ActionType = "command_update"
	ActionCommandDelete    ActionType = "command_delete"
	ActionReportReview     ActionType = "report_review"
	ActionReportDismiss    ActionType = "report_dismiss"
	ActionLogin            ActionType = "login"
	ActionLogout           ActionType = "logout"
	ActionLoginFailed      ActionType = "login_failed"
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
	Success      bool         `gorm:"default:true;index" json:"success"`
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

	err := database.DB.Create(log).Error
	return log, err
}

// GetAuditLogs retrieves audit logs with optional filters
func GetAuditLogs(filters map[string]interface{}, limit int, offset int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := database.DB.Model(&AuditLog{}).Preload("User")

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
func GetAuditLogsByUser(userID uint, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetAuditLogsByAction retrieves audit logs for a specific action type
func GetAuditLogsByAction(action ActionType, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := database.DB.Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetAuditLogsByTarget retrieves audit logs for a specific target
func GetAuditLogsByTarget(targetType string, targetID string, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := database.DB.Where("target_type = ? AND target_id = ?", targetType, targetID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetRecentAuditLogs retrieves the most recent audit logs
func GetRecentAuditLogs(limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := database.DB.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
