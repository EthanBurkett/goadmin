package metrics

import (
	"fmt"
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/ethanburkett/goadmin/app/plugins"
)

// Metrics holds various system metrics
type Metrics struct {
	// Database metrics
	DBConnections   int     `json:"db_connections"`
	DBIdleConns     int     `json:"db_idle_conns"`
	DBOpenConns     int     `json:"db_open_conns"`
	DBWaitCount     int64   `json:"db_wait_count"`
	DBWaitDuration  float64 `json:"db_wait_duration_ms"`
	DBMaxIdleClosed int64   `json:"db_max_idle_closed"`
	DBMaxLifeClosed int64   `json:"db_max_life_closed"`

	// Audit log metrics
	TotalAuditLogs   int64   `json:"total_audit_logs"`
	ArchivedLogs     int64   `json:"archived_audit_logs"`
	AuditSuccessRate float64 `json:"audit_success_rate"`

	// User metrics
	TotalUsers   int64 `json:"total_users"`
	ActiveUsers  int64 `json:"active_users"`
	PendingUsers int64 `json:"pending_users"`

	// Report metrics
	TotalReports   int64 `json:"total_reports"`
	PendingReports int64 `json:"pending_reports"`

	// Ban metrics
	TotalBans  int64 `json:"total_bans"`
	ActiveBans int64 `json:"active_bans"`

	// Command metrics
	TotalCommands  int64 `json:"total_commands"`
	CustomCommands int64 `json:"custom_commands"`
	PluginCommands int64 `json:"plugin_commands"`

	// Cache metrics
	CacheSize int `json:"cache_size"`

	// Uptime
	UptimeSeconds int64 `json:"uptime_seconds"`
}

var startTime = time.Now()

// GetMetrics collects and returns current system metrics
func GetMetrics() (*Metrics, error) {
	m := &Metrics{}

	// Database connection pool stats
	sqlDB, err := database.DB.DB()
	if err == nil {
		stats := sqlDB.Stats()
		m.DBOpenConns = stats.OpenConnections
		m.DBIdleConns = stats.Idle
		m.DBWaitCount = stats.WaitCount
		m.DBWaitDuration = float64(stats.WaitDuration.Milliseconds())
		m.DBMaxIdleClosed = stats.MaxIdleClosed
		m.DBMaxLifeClosed = stats.MaxLifetimeClosed
	}

	// Audit log metrics
	var totalAudit, archivedAudit, successAudit int64
	database.DB.Model(&models.AuditLog{}).Count(&totalAudit)
	database.DB.Unscoped().Model(&models.AuditLog{}).Where("deleted_at IS NOT NULL").Count(&archivedAudit)
	database.DB.Model(&models.AuditLog{}).Where("success = ?", true).Count(&successAudit)

	m.TotalAuditLogs = totalAudit
	m.ArchivedLogs = archivedAudit
	if totalAudit > 0 {
		m.AuditSuccessRate = float64(successAudit) / float64(totalAudit) * 100
	}

	// User metrics
	database.DB.Model(&models.User{}).Count(&m.TotalUsers)
	database.DB.Model(&models.User{}).Where("approved = ?", true).Count(&m.ActiveUsers)
	database.DB.Model(&models.User{}).Where("approved = ?", false).Count(&m.PendingUsers)

	// Report metrics
	database.DB.Model(&models.Report{}).Count(&m.TotalReports)
	database.DB.Model(&models.Report{}).Where("status = ?", "pending").Count(&m.PendingReports)

	// Ban metrics
	database.DB.Model(&models.TempBan{}).Count(&m.TotalBans)
	database.DB.Model(&models.TempBan{}).Where("active = ? AND expires_at > ?", true, time.Now()).Count(&m.ActiveBans)

	// Command metrics
	database.DB.Model(&models.CustomCommand{}).Count(&m.TotalCommands)
	database.DB.Model(&models.CustomCommand{}).Where("is_built_in = ?", false).Count(&m.CustomCommands)

	// Plugin commands are stored in memory, not in the database
	if plugins.GlobalPluginManager != nil {
		commandAPI := plugins.GlobalPluginManager.GetCommandAPI()
		if commandAPI != nil {
			m.PluginCommands = int64(commandAPI.GetCommandCount())
		}
	}

	// Uptime
	m.UptimeSeconds = int64(time.Since(startTime).Seconds())

	return m, nil
}

// PrometheusFormat formats metrics in Prometheus exposition format
func (m *Metrics) PrometheusFormat() string {
	format := `# HELP goadmin_db_open_connections Number of open database connections
# TYPE goadmin_db_open_connections gauge
goadmin_db_open_connections %d

# HELP goadmin_db_idle_connections Number of idle database connections
# TYPE goadmin_db_idle_connections gauge
goadmin_db_idle_connections %d

# HELP goadmin_db_wait_count Total number of connections waited for
# TYPE goadmin_db_wait_count counter
goadmin_db_wait_count %d

# HELP goadmin_db_wait_duration_ms Total time blocked waiting for connections (ms)
# TYPE goadmin_db_wait_duration_ms counter
goadmin_db_wait_duration_ms %.2f

# HELP goadmin_audit_logs_total Total number of audit logs
# TYPE goadmin_audit_logs_total gauge
goadmin_audit_logs_total %d

# HELP goadmin_audit_logs_archived Number of archived audit logs
# TYPE goadmin_audit_logs_archived gauge
goadmin_audit_logs_archived %d

# HELP goadmin_audit_success_rate Audit log success rate percentage
# TYPE goadmin_audit_success_rate gauge
goadmin_audit_success_rate %.2f

# HELP goadmin_users_total Total number of users
# TYPE goadmin_users_total gauge
goadmin_users_total %d

# HELP goadmin_users_active Number of active (approved) users
# TYPE goadmin_users_active gauge
goadmin_users_active %d

# HELP goadmin_users_pending Number of pending (unapproved) users
# TYPE goadmin_users_pending gauge
goadmin_users_pending %d

# HELP goadmin_reports_total Total number of reports
# TYPE goadmin_reports_total gauge
goadmin_reports_total %d

# HELP goadmin_reports_pending Number of pending reports
# TYPE goadmin_reports_pending gauge
goadmin_reports_pending %d

# HELP goadmin_bans_total Total number of bans
# TYPE goadmin_bans_total gauge
goadmin_bans_total %d

# HELP goadmin_bans_active Number of active bans
# TYPE goadmin_bans_active gauge
goadmin_bans_active %d

# HELP goadmin_commands_total Total number of commands
# TYPE goadmin_commands_total gauge
goadmin_commands_total %d

# HELP goadmin_commands_custom Number of custom commands
# TYPE goadmin_commands_custom gauge
goadmin_commands_custom %d

# HELP goadmin_uptime_seconds System uptime in seconds
# TYPE goadmin_uptime_seconds counter
goadmin_uptime_seconds %d
`

	return fmt.Sprintf(format,
		m.DBOpenConns,
		m.DBIdleConns,
		m.DBWaitCount,
		m.DBWaitDuration,
		m.TotalAuditLogs,
		m.ArchivedLogs,
		m.AuditSuccessRate,
		m.TotalUsers,
		m.ActiveUsers,
		m.PendingUsers,
		m.TotalReports,
		m.PendingReports,
		m.TotalBans,
		m.ActiveBans,
		m.TotalCommands,
		m.CustomCommands,
		m.UptimeSeconds,
	)
}
