//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ethanburkett/goadmin/app/config"
	"github.com/ethanburkett/goadmin/app/database"
	"github.com/ethanburkett/goadmin/app/logger"
	"gorm.io/gorm"
)

type ValidationResult struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Count       int64  `json:"count,omitempty"`
	Severity    string `json:"severity"` // "info", "warning", "error"
}

type ValidationReport struct {
	Results []ValidationResult `json:"results"`
	Summary struct {
		TotalIssues  int `json:"totalIssues"`
		Errors       int `json:"errors"`
		Warnings     int `json:"warnings"`
		InfoMessages int `json:"infoMessages"`
	} `json:"summary"`
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	logger.Init("validate_db", cfg.Environment == "development")

	// Initialize database
	database.Init()

	// Run validation
	report := runValidation(database.DB)

	// Print report
	printReport(report)

	// Exit with appropriate code
	if report.Summary.Errors > 0 {
		os.Exit(1)
	}
}

func runValidation(db *gorm.DB) ValidationReport {
	var report ValidationReport

	// Check for orphaned sessions
	report.Results = append(report.Results, checkOrphanedSessions(db)...)

	// Check for orphaned role assignments
	report.Results = append(report.Results, checkOrphanedRoleAssignments(db)...)

	// Check for orphaned permission assignments
	report.Results = append(report.Results, checkOrphanedPermissionAssignments(db)...)

	// Check for orphaned reports
	report.Results = append(report.Results, checkOrphanedReports(db)...)

	// Check for orphaned temp bans
	report.Results = append(report.Results, checkOrphanedTempBans(db)...)

	// Check for orphaned command history
	report.Results = append(report.Results, checkOrphanedCommandHistory(db)...)

	// Check for orphaned in-game players
	report.Results = append(report.Results, checkOrphanedInGamePlayers(db)...)

	// Check for orphaned webhook deliveries
	report.Results = append(report.Results, checkOrphanedWebhookDeliveries(db)...)

	// Check for orphaned audit logs
	report.Results = append(report.Results, checkOrphanedAuditLogs(db)...)

	// Check for missing indexes
	report.Results = append(report.Results, checkMissingIndexes(db)...)

	// Calculate summary
	for _, result := range report.Results {
		report.Summary.TotalIssues++
		switch result.Severity {
		case "error":
			report.Summary.Errors++
		case "warning":
			report.Summary.Warnings++
		case "info":
			report.Summary.InfoMessages++
		}
	}

	return report
}

func checkOrphanedSessions(db *gorm.DB) []ValidationResult {
	var results []ValidationResult
	var count int64

	// Find sessions with non-existent users
	err := db.Table("sessions").
		Joins("LEFT JOIN users ON sessions.user_id = users.id").
		Where("users.id IS NULL").
		Count(&count).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_sessions",
			Description: fmt.Sprintf("Error checking orphaned sessions: %v", err),
			Severity:    "error",
		})
	} else if count > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_sessions",
			Description: "Found sessions with non-existent users",
			Count:       count,
			Severity:    "error",
		})
	} else {
		results = append(results, ValidationResult{
			Type:        "orphaned_sessions",
			Description: "No orphaned sessions found",
			Severity:    "info",
		})
	}

	return results
}

func checkOrphanedRoleAssignments(db *gorm.DB) []ValidationResult {
	var results []ValidationResult
	var userCount, roleCount int64

	// Check for user_roles entries with non-existent users
	err := db.Table("user_roles").
		Joins("LEFT JOIN users ON user_roles.user_id = users.id").
		Where("users.id IS NULL").
		Count(&userCount).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_user_roles",
			Description: fmt.Sprintf("Error checking orphaned user roles: %v", err),
			Severity:    "error",
		})
	} else if userCount > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_user_roles",
			Description: "Found user_roles entries with non-existent users",
			Count:       userCount,
			Severity:    "error",
		})
	}

	// Check for user_roles entries with non-existent roles
	err = db.Table("user_roles").
		Joins("LEFT JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.id IS NULL").
		Count(&roleCount).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_role_assignments",
			Description: fmt.Sprintf("Error checking orphaned role assignments: %v", err),
			Severity:    "error",
		})
	} else if roleCount > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_role_assignments",
			Description: "Found user_roles entries with non-existent roles",
			Count:       roleCount,
			Severity:    "error",
		})
	}

	if userCount == 0 && roleCount == 0 && err == nil {
		results = append(results, ValidationResult{
			Type:        "user_role_assignments",
			Description: "No orphaned user role assignments found",
			Severity:    "info",
		})
	}

	return results
}

func checkOrphanedPermissionAssignments(db *gorm.DB) []ValidationResult {
	var results []ValidationResult
	var roleCount, permCount int64

	// Check for role_permissions entries with non-existent roles
	err := db.Table("role_permissions").
		Joins("LEFT JOIN roles ON role_permissions.role_id = roles.id").
		Where("roles.id IS NULL").
		Count(&roleCount).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_role_permissions",
			Description: fmt.Sprintf("Error checking orphaned role permissions: %v", err),
			Severity:    "error",
		})
	} else if roleCount > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_role_permissions",
			Description: "Found role_permissions entries with non-existent roles",
			Count:       roleCount,
			Severity:    "error",
		})
	}

	// Check for role_permissions entries with non-existent permissions
	err = db.Table("role_permissions").
		Joins("LEFT JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("permissions.id IS NULL").
		Count(&permCount).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_permission_assignments",
			Description: fmt.Sprintf("Error checking orphaned permission assignments: %v", err),
			Severity:    "error",
		})
	} else if permCount > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_permission_assignments",
			Description: "Found role_permissions entries with non-existent permissions",
			Count:       permCount,
			Severity:    "error",
		})
	}

	if roleCount == 0 && permCount == 0 && err == nil {
		results = append(results, ValidationResult{
			Type:        "role_permission_assignments",
			Description: "No orphaned role permission assignments found",
			Severity:    "info",
		})
	}

	return results
}

func checkOrphanedReports(db *gorm.DB) []ValidationResult {
	var results []ValidationResult
	var count int64

	// Find reports with non-existent reviewer users (allowing NULL)
	err := db.Table("reports").
		Joins("LEFT JOIN users ON reports.reviewed_by_user_id = users.id").
		Where("reports.reviewed_by_user_id IS NOT NULL AND users.id IS NULL").
		Count(&count).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_reports",
			Description: fmt.Sprintf("Error checking orphaned reports: %v", err),
			Severity:    "error",
		})
	} else if count > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_reports",
			Description: "Found reports with non-existent reviewer users",
			Count:       count,
			Severity:    "warning",
		})
	} else {
		results = append(results, ValidationResult{
			Type:        "orphaned_reports",
			Description: "No orphaned reports found",
			Severity:    "info",
		})
	}

	return results
}

func checkOrphanedTempBans(db *gorm.DB) []ValidationResult {
	var results []ValidationResult
	var count int64

	// Find temp bans with non-existent users (allowing NULL)
	err := db.Table("temp_bans").
		Joins("LEFT JOIN users ON temp_bans.banned_by_user = users.id").
		Where("temp_bans.banned_by_user IS NOT NULL AND users.id IS NULL").
		Count(&count).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_temp_bans",
			Description: fmt.Sprintf("Error checking orphaned temp bans: %v", err),
			Severity:    "error",
		})
	} else if count > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_temp_bans",
			Description: "Found temp bans with non-existent users",
			Count:       count,
			Severity:    "warning",
		})
	} else {
		results = append(results, ValidationResult{
			Type:        "orphaned_temp_bans",
			Description: "No orphaned temp bans found",
			Severity:    "info",
		})
	}

	return results
}

func checkOrphanedCommandHistory(db *gorm.DB) []ValidationResult {
	var results []ValidationResult
	var count int64

	// Find command history with non-existent users
	err := db.Table("command_histories").
		Joins("LEFT JOIN users ON command_histories.user_id = users.id").
		Where("users.id IS NULL").
		Count(&count).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_command_history",
			Description: fmt.Sprintf("Error checking orphaned command history: %v", err),
			Severity:    "error",
		})
	} else if count > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_command_history",
			Description: "Found command history entries with non-existent users",
			Count:       count,
			Severity:    "error",
		})
	} else {
		results = append(results, ValidationResult{
			Type:        "orphaned_command_history",
			Description: "No orphaned command history found",
			Severity:    "info",
		})
	}

	return results
}

func checkOrphanedInGamePlayers(db *gorm.DB) []ValidationResult {
	var results []ValidationResult
	var count int64

	// Find in-game players with non-existent groups (allowing NULL)
	err := db.Table("in_game_players").
		Joins("LEFT JOIN groups ON in_game_players.group_id = groups.id").
		Where("in_game_players.group_id IS NOT NULL AND groups.id IS NULL").
		Count(&count).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_in_game_players",
			Description: fmt.Sprintf("Error checking orphaned in-game players: %v", err),
			Severity:    "error",
		})
	} else if count > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_in_game_players",
			Description: "Found in-game players with non-existent groups",
			Count:       count,
			Severity:    "warning",
		})
	} else {
		results = append(results, ValidationResult{
			Type:        "orphaned_in_game_players",
			Description: "No orphaned in-game players found",
			Severity:    "info",
		})
	}

	return results
}

func checkOrphanedWebhookDeliveries(db *gorm.DB) []ValidationResult {
	var results []ValidationResult
	var count int64

	// Find webhook deliveries with non-existent webhooks
	err := db.Table("webhook_deliveries").
		Joins("LEFT JOIN webhooks ON webhook_deliveries.webhook_id = webhooks.id").
		Where("webhooks.id IS NULL").
		Count(&count).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_webhook_deliveries",
			Description: fmt.Sprintf("Error checking orphaned webhook deliveries: %v", err),
			Severity:    "error",
		})
	} else if count > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_webhook_deliveries",
			Description: "Found webhook deliveries with non-existent webhooks",
			Count:       count,
			Severity:    "warning",
		})
	} else {
		results = append(results, ValidationResult{
			Type:        "orphaned_webhook_deliveries",
			Description: "No orphaned webhook deliveries found",
			Severity:    "info",
		})
	}

	return results
}

func checkOrphanedAuditLogs(db *gorm.DB) []ValidationResult {
	var results []ValidationResult
	var count int64

	// Find audit logs with non-existent users (allowing NULL for system actions)
	err := db.Table("audit_logs").
		Joins("LEFT JOIN users ON audit_logs.user_id = users.id").
		Where("audit_logs.user_id IS NOT NULL AND users.id IS NULL").
		Count(&count).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "orphaned_audit_logs",
			Description: fmt.Sprintf("Error checking orphaned audit logs: %v", err),
			Severity:    "error",
		})
	} else if count > 0 {
		results = append(results, ValidationResult{
			Type:        "orphaned_audit_logs",
			Description: "Found audit logs with non-existent users",
			Count:       count,
			Severity:    "warning",
		})
	} else {
		results = append(results, ValidationResult{
			Type:        "orphaned_audit_logs",
			Description: "No orphaned audit logs found",
			Severity:    "info",
		})
	}

	return results
}

func checkMissingIndexes(db *gorm.DB) []ValidationResult {
	var results []ValidationResult

	// This is informational - GORM handles indexes via tags
	// We can check if critical indexes exist
	type IndexInfo struct {
		Name string
	}

	var indexes []IndexInfo

	// Check for indexes on sessions.user_id
	err := db.Raw("SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='sessions' AND sql LIKE '%user_id%'").
		Scan(&indexes).Error

	if err != nil {
		results = append(results, ValidationResult{
			Type:        "index_check",
			Description: fmt.Sprintf("Error checking indexes: %v", err),
			Severity:    "warning",
		})
	} else if len(indexes) == 0 {
		results = append(results, ValidationResult{
			Type:        "missing_index",
			Description: "No index found on sessions.user_id (may impact performance)",
			Severity:    "warning",
		})
	} else {
		results = append(results, ValidationResult{
			Type:        "index_check",
			Description: "Critical indexes present",
			Severity:    "info",
		})
	}

	return results
}

func printReport(report ValidationReport) {
	// Print as JSON for machine readability
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal report:", err)
	}
	fmt.Println(string(jsonData))

	// Print summary
	fmt.Println("\n=== VALIDATION SUMMARY ===")
	fmt.Printf("Total Issues: %d\n", report.Summary.TotalIssues)
	fmt.Printf("Errors: %d\n", report.Summary.Errors)
	fmt.Printf("Warnings: %d\n", report.Summary.Warnings)
	fmt.Printf("Info Messages: %d\n", report.Summary.InfoMessages)

	if report.Summary.Errors > 0 {
		fmt.Println("\n⚠️  Database integrity issues found! Please review and fix.")
	} else if report.Summary.Warnings > 0 {
		fmt.Println("\n⚠️  Some warnings detected, but no critical errors.")
	} else {
		fmt.Println("\n✅ Database integrity check passed!")
	}
}
