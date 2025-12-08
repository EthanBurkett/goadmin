package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/ethanburkett/goadmin/app/webhook"
	"gorm.io/gorm"
)

// handleReportCommand allows players to report others
func (ch *CommandHandler) handleReportCommand(ch2 *CommandHandler, playerName, playerGUID string, args []string) error {
	if len(args) < 2 {
		ch.sendPlayerMessage(playerName, "Usage: !report <player> <reason>")
		return nil
	}

	reportedPlayerName := args[0]
	reason := strings.Join(args[1:], " ")

	status, err := ch.rcon.Status()
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to get server status")
		return err
	}

	var reportedGUID string
	searchName := strings.ToLower(reportedPlayerName)
	for _, player := range status.Players {
		if strings.ToLower(player.StrippedName) == searchName || strings.Contains(strings.ToLower(player.StrippedName), searchName) {
			reportedGUID = player.Uuid
			reportedPlayerName = player.StrippedName
			break
		}
	}

	if reportedGUID == "" {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("Player '%s' not found online", reportedPlayerName))
		return nil
	}

	if reportedGUID == playerGUID {
		ch.sendPlayerMessage(playerName, "You cannot report yourself")
		return nil
	}

	report, err := models.CreateReport(playerName, playerGUID, reportedPlayerName, reportedGUID, reason)
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to submit report")
		return err
	}

	// Dispatch webhook event
	go webhook.GlobalDispatcher.Dispatch(models.WebhookEventReportCreated, map[string]interface{}{
		"report_id":     report.ID,
		"reporter_name": playerName,
		"reporter_guid": playerGUID,
		"reported_name": reportedPlayerName,
		"reported_guid": reportedGUID,
		"reason":        reason,
		"status":        report.Status,
		"source":        "in-game",
		"created_at":    report.CreatedAt.Format(time.RFC3339),
	})

	ch.sendPlayerMessage(playerName, fmt.Sprintf("^2Report submitted for %s (ID: #%d)", reportedPlayerName, report.ID))
	logger.Info(fmt.Sprintf("Player %s reported %s (GUID: %s) for: %s", playerName, reportedPlayerName, reportedGUID, reason))

	return nil
}

// parseDuration parses duration strings like "5m", "2h", "3d", "1M", "2y"
func parseDuration(input string) (time.Duration, error) {
	if len(input) < 2 {
		return 0, fmt.Errorf("invalid duration format")
	}

	numStr := input[:len(input)-1]
	unit := input[len(input)-1:]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number in duration: %s", numStr)
	}

	if num <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}

	switch unit {
	case "m":
		return time.Duration(num) * time.Minute, nil
	case "h":
		return time.Duration(num) * time.Hour, nil
	case "d":
		return time.Duration(num) * 24 * time.Hour, nil
	case "M":
		return time.Duration(num) * 30 * 24 * time.Hour, nil
	case "y":
		return time.Duration(num) * 365 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid duration unit: %s (use m, h, d, M, or y)", unit)
	}
}

// handleTempBanCommand temporarily bans a player for a specified duration
func (ch *CommandHandler) handleTempBanCommand(ch2 *CommandHandler, playerName, playerGUID string, args []string) error {
	if len(args) < 3 {
		ch.sendPlayerMessage(playerName, "Usage: !tempban <player> <duration> <reason>")
		ch.sendPlayerMessage(playerName, "Duration format: {number}{m/h/d/M/y} (e.g., 5m, 2h, 3d, 1M, 2y)")
		return nil
	}

	bannedPlayerName := args[0]
	durationStr := args[1]
	reason := strings.Join(args[2:], " ")

	duration, err := parseDuration(durationStr)
	if err != nil {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("Invalid duration: %v", err))
		ch.sendPlayerMessage(playerName, "Duration format: {number}{m/h/d/M/y} (e.g., 5m, 2h, 3d, 1M, 2y)")
		return nil
	}

	status, err := ch.rcon.Status()
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to get server status")
		return err
	}

	var bannedGUID string
	var bannedEntityID int
	searchName := strings.ToLower(bannedPlayerName)
	for _, player := range status.Players {
		if strings.ToLower(player.StrippedName) == searchName || strings.Contains(strings.ToLower(player.StrippedName), searchName) {
			bannedGUID = player.Uuid
			bannedPlayerName = player.StrippedName
			bannedEntityID = player.ID
			break
		}
	}

	if bannedGUID == "" {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("Player '%s' not found online", bannedPlayerName))
		return nil
	}

	// Check command throttling (prevent targeting same player too frequently - 30 second cooldown)
	throttleResult := models.CommandThrottlerInstance.CheckThrottle(playerGUID, bannedGUID, "tempban", 30*time.Second)
	if !throttleResult.Allowed {
		ch.sendPlayerMessage(playerName, fmt.Sprintf("^1%s", throttleResult.Reason))
		ch.sendPlayerMessage(playerName, fmt.Sprintf("^1Please wait %d seconds", int(throttleResult.TimeRemaining.Seconds())))
		return nil
	}

	// Check for ban loop abuse (5 bans in 15 minutes)
	banLoopResult, err := models.BanLoopDetectorInstance.CheckBanLoop(bannedGUID, 15*time.Minute, 5)
	if err == nil && banLoopResult.IsAbuse {
		ch.sendPlayerMessage(playerName, "^1Warning: This player has been banned multiple times recently")
		ch.sendPlayerMessage(playerName, fmt.Sprintf("^1%s", banLoopResult.Reason))
		logger.Info(fmt.Sprintf("[BAN LOOP DETECTED] %s (GUID: %s) - %d bans in %v",
			bannedPlayerName, bannedGUID, banLoopResult.RecentBanCount, banLoopResult.TimeWindow))
	}

	tempBan, err := models.CreateTempBan(bannedPlayerName, bannedGUID, reason, duration, nil)
	if err != nil {
		ch.sendPlayerMessage(playerName, "Failed to create temp ban")
		return err
	}

	// Dispatch webhook event
	go webhook.GlobalDispatcher.Dispatch(models.WebhookEventPlayerBanned, map[string]interface{}{
		"player_name":    bannedPlayerName,
		"player_guid":    bannedGUID,
		"banned_by":      playerName,
		"reason":         reason,
		"duration":       durationStr,
		"expires_at":     tempBan.ExpiresAt.Format(time.RFC3339),
		"ban_type":       "temporary",
		"source":         "in-game",
		"recent_bans":    banLoopResult.RecentBanCount,
		"time_window":    banLoopResult.TimeWindow.String(),
		"abuse_detected": banLoopResult.IsAbuse,
	})

	// Log audit entry for in-game temp ban
	models.CreateAuditLog(
		ch.db.(*gorm.DB),
		nil,
		playerName,
		"",
		models.ActionTempBanPlayer,
		models.SourceInGame,
		true,
		"",
		"player",
		bannedGUID,
		bannedPlayerName,
		fmt.Sprintf(`{"reason": "%s", "duration": "%s", "issued_by": "%s"}`, reason, durationStr, playerName),
		fmt.Sprintf("Temporarily banned for %s: %s", durationStr, reason),
	)

	kickCmd := fmt.Sprintf("clientkick %d \"Temp banned: %s (Expires: %s)\"",
		bannedEntityID,
		reason,
		tempBan.ExpiresAt.Format("2006-01-02 15:04"))
	ch.rcon.SendCommand(kickCmd)

	ch.sendPlayerMessage(playerName, fmt.Sprintf("^2%s has been temp banned for %s", bannedPlayerName, durationStr))
	logger.Info(fmt.Sprintf("Player %s temp banned %s (GUID: %s) for %s: %s", playerName, bannedPlayerName, bannedGUID, durationStr, reason))

	return nil
}
