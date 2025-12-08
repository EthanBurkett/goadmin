package rest

import (
	"encoding/json"
	"fmt"

	"github.com/ethanburkett/goadmin/app/database"
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

// AuditHelper provides helper functions for creating audit logs
type AuditHelper struct{}

// LogAction creates an audit log entry from a gin context
func (ah *AuditHelper) LogAction(
	c *gin.Context,
	action models.ActionType,
	source models.ActionSource,
	success bool,
	errorMsg string,
	targetType string,
	targetID string,
	targetName string,
	metadata map[string]interface{},
	result string,
) error {
	var userID *uint
	var username string
	ipAddress := c.ClientIP()

	// Try to get user from context
	if userVal, exists := c.Get("user"); exists {
		if user, ok := userVal.(*models.User); ok {
			userID = &user.ID
			username = user.Username
		}
	}

	// Encode metadata as JSON
	var metadataJSON string
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err == nil {
			metadataJSON = string(metadataBytes)
		}
	}

	_, err := models.CreateAuditLog(
		database.DB,
		userID,
		username,
		ipAddress,
		action,
		source,
		success,
		errorMsg,
		targetType,
		targetID,
		targetName,
		metadataJSON,
		result,
	)

	return err
}

// LogBan logs a permanent ban action
func (ah *AuditHelper) LogBan(c *gin.Context, playerName, playerGUID, reason string, success bool, errorMsg string) error {
	return ah.LogAction(
		c,
		models.ActionBanPlayer,
		models.SourceWebUI,
		success,
		errorMsg,
		"player",
		playerGUID,
		playerName,
		map[string]interface{}{
			"reason": reason,
		},
		fmt.Sprintf("Permanently banned: %s", reason),
	)
}

// LogTempBan logs a temporary ban action
func (ah *AuditHelper) LogTempBan(c *gin.Context, playerName, playerGUID, reason string, durationHours int, success bool, errorMsg string) error {
	return ah.LogAction(
		c,
		models.ActionTempBanPlayer,
		models.SourceWebUI,
		success,
		errorMsg,
		"player",
		playerGUID,
		playerName,
		map[string]interface{}{
			"reason":         reason,
			"duration_hours": durationHours,
		},
		fmt.Sprintf("Temporarily banned for %d hours: %s", durationHours, reason),
	)
}

// LogKick logs a player kick action
func (ah *AuditHelper) LogKick(c *gin.Context, playerName, playerID, reason string, success bool, errorMsg string) error {
	return ah.LogAction(
		c,
		models.ActionKickPlayer,
		models.SourceWebUI,
		success,
		errorMsg,
		"player",
		playerID,
		playerName,
		map[string]interface{}{
			"reason": reason,
		},
		fmt.Sprintf("Kicked: %s", reason),
	)
}

// LogRconCommand logs an RCON command execution
func (ah *AuditHelper) LogRconCommand(c *gin.Context, command, response string, success bool, errorMsg string) error {
	return ah.LogAction(
		c,
		models.ActionRconCommand,
		models.SourceWebUI,
		success,
		errorMsg,
		"command",
		"",
		command,
		map[string]interface{}{
			"command": command,
		},
		response,
	)
}

// LogReportAction logs a report review/action
func (ah *AuditHelper) LogReportAction(c *gin.Context, reportID uint, action, reason string, success bool, errorMsg string) error {
	var actionType models.ActionType
	switch action {
	case "dismiss":
		actionType = models.ActionReportDismiss
	default:
		actionType = models.ActionReportReview
	}

	return ah.LogAction(
		c,
		actionType,
		models.SourceWebUI,
		success,
		errorMsg,
		"report",
		fmt.Sprintf("%d", reportID),
		"",
		map[string]interface{}{
			"action": action,
			"reason": reason,
		},
		fmt.Sprintf("Report %s: %s", action, reason),
	)
}

// LogSecurityViolation logs a security violation attempt
func (ah *AuditHelper) LogSecurityViolation(c *gin.Context, violationType, attemptedCommand, reason string) error {
	return ah.LogAction(
		c,
		models.ActionSecurityViolation,
		models.SourceWebUI,
		false, // Always marked as failed
		reason,
		"security",
		violationType,
		attemptedCommand,
		map[string]interface{}{
			"violation_type":    violationType,
			"attempted_command": attemptedCommand,
			"validation_error":  reason,
		},
		fmt.Sprintf("Security violation: %s - %s", violationType, reason),
	)
}

var Audit = &AuditHelper{}
