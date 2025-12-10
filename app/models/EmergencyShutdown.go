package models

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"github.com/ethanburkett/goadmin/app/logger"
	"go.uber.org/zap"
)

// EmergencyShutdownManager handles automatic shutdown of commands on abuse detection
type EmergencyShutdownManager struct {
	mu               sync.RWMutex
	disabledCommands map[string]*ShutdownInfo
	userAlerts       map[uint]int // userID -> alert count
	cleanupTicker    *time.Ticker
	cleanupStop      chan bool
}

// ShutdownInfo contains information about why a command was shut down
type ShutdownInfo struct {
	Command     string
	Reason      string
	DisabledAt  time.Time
	DisabledBy  string // "system" or user ID
	UserID      *uint  // Optional user who triggered shutdown
	ReenableAt  time.Time
	AutoRenable bool
}

var GlobalEmergencyShutdown *EmergencyShutdownManager

// InitEmergencyShutdown initializes the global emergency shutdown manager
func InitEmergencyShutdown() {
	GlobalEmergencyShutdown = &EmergencyShutdownManager{
		disabledCommands: make(map[string]*ShutdownInfo),
		userAlerts:       make(map[uint]int),
		cleanupTicker:    time.NewTicker(1 * time.Minute),
		cleanupStop:      make(chan bool),
	}

	// Start cleanup goroutine for auto-reenable
	go GlobalEmergencyShutdown.cleanupLoop()

	logger.Info("Emergency shutdown manager initialized")
}

// cleanupLoop periodically checks for commands to re-enable
func (esm *EmergencyShutdownManager) cleanupLoop() {
	for {
		select {
		case <-esm.cleanupTicker.C:
			esm.checkAutoReenable()
		case <-esm.cleanupStop:
			return
		}
	}
}

// checkAutoReenable re-enables commands that have passed their reenable time
func (esm *EmergencyShutdownManager) checkAutoReenable() {
	esm.mu.Lock()
	defer esm.mu.Unlock()

	now := time.Now()
	for cmd, info := range esm.disabledCommands {
		if info.AutoRenable && now.After(info.ReenableAt) {
			delete(esm.disabledCommands, cmd)
			logger.Info("Auto-reenabled command after cooldown",
				zap.String("command", cmd),
				zap.String("reason", info.Reason))
		}
	}
}

// DisableCommand disables a command due to abuse detection
func (esm *EmergencyShutdownManager) DisableCommand(command, reason string, userID *uint, duration time.Duration) {
	esm.mu.Lock()
	defer esm.mu.Unlock()

	info := &ShutdownInfo{
		Command:     command,
		Reason:      reason,
		DisabledAt:  time.Now(),
		DisabledBy:  "system",
		UserID:      userID,
		AutoRenable: duration > 0,
		ReenableAt:  time.Now().Add(duration),
	}

	esm.disabledCommands[command] = info

	logger.Warn("Command disabled due to abuse",
		zap.String("command", command),
		zap.String("reason", reason),
		zap.Duration("duration", duration))

	// Notify super admins
	go esm.notifySuperAdmins(command, reason, userID)
}

// IsCommandDisabled checks if a command is currently disabled
func (esm *EmergencyShutdownManager) IsCommandDisabled(command string) (bool, *ShutdownInfo) {
	esm.mu.RLock()
	defer esm.mu.RUnlock()

	info, disabled := esm.disabledCommands[command]
	return disabled, info
}

// EnableCommand manually re-enables a disabled command
func (esm *EmergencyShutdownManager) EnableCommand(command string, adminUserID uint) error {
	esm.mu.Lock()
	defer esm.mu.Unlock()

	info, exists := esm.disabledCommands[command]
	if !exists {
		return fmt.Errorf("command %s is not disabled", command)
	}

	delete(esm.disabledCommands, command)

	logger.Info("Command manually re-enabled",
		zap.String("command", command),
		zap.Uint("admin_user_id", adminUserID),
		zap.String("original_reason", info.Reason))

	return nil
}

// GetDisabledCommands returns a list of currently disabled commands
func (esm *EmergencyShutdownManager) GetDisabledCommands() map[string]*ShutdownInfo {
	esm.mu.RLock()
	defer esm.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make(map[string]*ShutdownInfo)
	for k, v := range esm.disabledCommands {
		result[k] = v
	}
	return result
}

// IncrementUserAlert increments the alert count for a user
func (esm *EmergencyShutdownManager) IncrementUserAlert(userID uint) int {
	esm.mu.Lock()
	defer esm.mu.Unlock()

	esm.userAlerts[userID]++
	count := esm.userAlerts[userID]

	logger.Warn("User abuse alert incremented",
		zap.Uint("user_id", userID),
		zap.Int("alert_count", count))

	return count
}

// GetUserAlertCount returns the current alert count for a user
func (esm *EmergencyShutdownManager) GetUserAlertCount(userID uint) int {
	esm.mu.RLock()
	defer esm.mu.RUnlock()

	return esm.userAlerts[userID]
}

// ResetUserAlerts resets the alert count for a user
func (esm *EmergencyShutdownManager) ResetUserAlerts(userID uint) {
	esm.mu.Lock()
	defer esm.mu.Unlock()

	delete(esm.userAlerts, userID)

	logger.Info("User alerts reset",
		zap.Uint("user_id", userID))
}

// notifySuperAdmins notifies super admins about command shutdown
func (esm *EmergencyShutdownManager) notifySuperAdmins(command, reason string, userID *uint) {
	// Get all super admin users
	var users []User
	err := database.DB.Preload("Roles").Find(&users).Error
	if err != nil {
		logger.Error("Failed to fetch users for notification", zap.Error(err))
		return
	}

	var superAdmins []User
	for _, user := range users {
		for _, role := range user.Roles {
			if role.Name == "super_admin" {
				superAdmins = append(superAdmins, user)
				break
			}
		}
	}

	if len(superAdmins) == 0 {
		logger.Warn("No super admins found to notify about command shutdown")
		return
	}

	// Create audit log entry for the shutdown
	var userIDVal uint
	if userID != nil {
		userIDVal = *userID
	}

	metadata := map[string]interface{}{
		"command":      command,
		"reason":       reason,
		"disabled_at":  time.Now(),
		"super_admins": len(superAdmins),
	}

	metadataJSON, _ := json.Marshal(metadata)

	CreateAuditLog(
		database.DB,
		&userIDVal,
		"SYSTEM",
		"",
		ActionSecurityViolation,
		SourceSystem,
		true,
		"",
		"command",
		command,
		command,
		string(metadataJSON),
		fmt.Sprintf("Command '%s' automatically disabled: %s", command, reason),
	)

	logger.Info("Super admins notified about command shutdown",
		zap.String("command", command),
		zap.Int("admin_count", len(superAdmins)))
}

// Stop stops the emergency shutdown manager
func (esm *EmergencyShutdownManager) Stop() {
	esm.cleanupTicker.Stop()
	close(esm.cleanupStop)
	logger.Info("Emergency shutdown manager stopped")
}
