package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
)

// BanLoopDetector provides methods to detect ban/unban abuse patterns
type BanLoopDetector struct{}

// BanLoopResult contains information about potential ban loop abuse
type BanLoopResult struct {
	IsAbuse        bool
	Reason         string
	RecentBanCount int
	TimeWindow     time.Duration
}

// CheckBanLoop detects if a player is being ban/unban cycled
func (bld *BanLoopDetector) CheckBanLoop(playerGUID string, timeWindow time.Duration, threshold int) (*BanLoopResult, error) {
	db := database.DB

	// Count recent bans within time window
	var count int64
	cutoffTime := time.Now().Add(-timeWindow)

	err := db.Model(&TempBan{}).
		Where("player_guid = ? AND created_at >= ?", playerGUID, cutoffTime).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	result := &BanLoopResult{
		IsAbuse:        int(count) >= threshold,
		RecentBanCount: int(count),
		TimeWindow:     timeWindow,
	}

	if result.IsAbuse {
		result.Reason = "Player has been banned multiple times in a short period, indicating potential ban loop abuse"
	}

	return result, nil
}

// CheckCircularBan detects if an admin is rapidly banning/unbanning the same player
func (bld *BanLoopDetector) CheckCircularBan(playerGUID string, adminUserID *uint, timeWindow time.Duration, threshold int) (*BanLoopResult, error) {
	db := database.DB

	// Count recent bans by this admin for this player
	var count int64
	cutoffTime := time.Now().Add(-timeWindow)

	query := db.Model(&TempBan{}).
		Where("player_guid = ? AND created_at >= ?", playerGUID, cutoffTime)

	if adminUserID != nil {
		query = query.Where("banned_by_user = ?", *adminUserID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, err
	}

	result := &BanLoopResult{
		IsAbuse:        int(count) >= threshold,
		RecentBanCount: int(count),
		TimeWindow:     timeWindow,
	}

	if result.IsAbuse {
		if adminUserID != nil {
			result.Reason = "Admin is repeatedly banning the same player in a short period"
		} else {
			result.Reason = "Player is being repeatedly banned in a short period"
		}
	}

	return result, nil
}

// CheckTargetAbuse detects if a specific player is being targeted by multiple bans
func (bld *BanLoopDetector) CheckTargetAbuse(playerGUID string, timeWindow time.Duration, threshold int) (*BanLoopResult, error) {
	return bld.CheckBanLoop(playerGUID, timeWindow, threshold)
}

// GetRecentBanHistory retrieves recent ban history for a player
func (bld *BanLoopDetector) GetRecentBanHistory(playerGUID string, timeWindow time.Duration) ([]TempBan, error) {
	db := database.DB
	var bans []TempBan
	cutoffTime := time.Now().Add(-timeWindow)

	err := db.Preload("BannedBy").
		Where("player_guid = ? AND created_at >= ?", playerGUID, cutoffTime).
		Order("created_at DESC").
		Find(&bans).Error

	return bans, err
}

// GetBanPatternStats provides statistics about ban patterns for analysis
func (bld *BanLoopDetector) GetBanPatternStats(playerGUID string, timeWindow time.Duration) (map[string]interface{}, error) {
	db := database.DB
	cutoffTime := time.Now().Add(-timeWindow)

	// Count total bans
	var totalBans int64
	db.Model(&TempBan{}).
		Where("player_guid = ? AND created_at >= ?", playerGUID, cutoffTime).
		Count(&totalBans)

	// Count active bans
	var activeBans int64
	db.Model(&TempBan{}).
		Where("player_guid = ? AND created_at >= ? AND active = ?", playerGUID, cutoffTime, true).
		Count(&activeBans)

	// Count revoked bans (inactive before expiry)
	var revokedBans int64
	db.Model(&TempBan{}).
		Where("player_guid = ? AND created_at >= ? AND active = ? AND expires_at > ?",
			playerGUID, cutoffTime, false, time.Now()).
		Count(&revokedBans)

	// Count unique admins who banned this player
	var uniqueAdmins int64
	db.Model(&TempBan{}).
		Where("player_guid = ? AND created_at >= ? AND banned_by_user IS NOT NULL", playerGUID, cutoffTime).
		Distinct("banned_by_user").
		Count(&uniqueAdmins)

	stats := map[string]interface{}{
		"total_bans":      totalBans,
		"active_bans":     activeBans,
		"revoked_bans":    revokedBans,
		"unique_admins":   uniqueAdmins,
		"time_window":     timeWindow.String(),
		"suspicion_score": calculateSuspicionScore(int(totalBans), int(revokedBans), int(uniqueAdmins)),
	}

	return stats, nil
}

// calculateSuspicionScore generates a score (0-100) indicating likelihood of abuse
func calculateSuspicionScore(totalBans, revokedBans, uniqueAdmins int) int {
	score := 0

	// More than 5 bans in time window is suspicious
	if totalBans > 5 {
		score += 30
	} else if totalBans > 3 {
		score += 15
	}

	// High revoke rate is suspicious (ban/unban cycling)
	if totalBans > 0 {
		revokeRate := float64(revokedBans) / float64(totalBans)
		if revokeRate > 0.5 {
			score += 40
		} else if revokeRate > 0.3 {
			score += 20
		}
	}

	// Single admin repeatedly banning is more suspicious than multiple admins
	if totalBans > 3 && uniqueAdmins == 1 {
		score += 30
	}

	if score > 100 {
		score = 100
	}

	return score
}

// Global instance
var BanLoopDetectorInstance = &BanLoopDetector{}
