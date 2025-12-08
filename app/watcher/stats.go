package watcher

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/ethanburkett/goadmin/app/rcon"
)

type StatsCollector struct {
	rcon   *rcon.Client
	ticker *time.Ticker
	done   chan bool
}

func NewStatsCollector(rconClient *rcon.Client) *StatsCollector {
	return &StatsCollector{
		rcon: rconClient,
		done: make(chan bool),
	}
}

func (sc *StatsCollector) Start() {
	logger.Info("Starting stats collection service (every 1 minute)")
	sc.ticker = time.NewTicker(1 * time.Minute)

	// Collect initial stats immediately
	go sc.collectStats()

	go func() {
		for {
			select {
			case <-sc.ticker.C:
				sc.collectStats()
			case <-sc.done:
				return
			}
		}
	}()
}

func (sc *StatsCollector) Stop() {
	logger.Info("Stopping stats collection service")
	sc.ticker.Stop()
	sc.done <- true
}

func (sc *StatsCollector) collectStats() {
	logger.Debug("Collecting server stats...")

	// Collect server stats
	if err := sc.collectServerStats(); err != nil {
		logger.Error(fmt.Sprintf("Failed to collect server stats: %v", err))
	}

	// Collect system stats
	if err := sc.collectSystemStats(); err != nil {
		logger.Error(fmt.Sprintf("Failed to collect system stats: %v", err))
	}

	// Collect player stats
	if err := sc.collectPlayerStats(); err != nil {
		logger.Error(fmt.Sprintf("Failed to collect player stats: %v", err))
	}

	logger.Debug("Stats collection completed")
}

func (sc *StatsCollector) collectServerStats() error {
	// Get status
	statusResp, err := sc.rcon.SendCommand("status")
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	// Get serverinfo
	serverinfoResp, err := sc.rcon.SendCommand("serverinfo")
	if err != nil {
		return fmt.Errorf("failed to get serverinfo: %w", err)
	}

	// Parse status for player count
	playerCount := parsePlayerCount(statusResp)

	// Parse serverinfo for settings (format: "key value" on each line)
	maxPlayers := parseCvarInt(serverinfoResp, "sv_maxclients", 32)
	mapName := parseCvarString(serverinfoResp, "mapname", "unknown")
	gametype := parseCvarString(serverinfoResp, "g_gametype", "unknown")
	hostname := parseCvarString(serverinfoResp, "sv_hostname", "CoD4 Server")
	fps := 20 // sv_fps not in serverinfo, default to 20

	// Parse uptime from serverinfo (format: "uptime               5 hours")
	uptime := parseUptime(serverinfoResp)

	return models.CreateServerStats(playerCount, maxPlayers, mapName, gametype, hostname, fps, uptime)
}

func (sc *StatsCollector) collectSystemStats() error {
	// Get meminfo
	meminfoResp, err := sc.rcon.SendCommand("meminfo")
	if err != nil {
		return fmt.Errorf("failed to get meminfo: %w", err)
	}

	// Parse memory info
	// Format: "10485760 bytes total hunk" and "66471 total hunk in use"
	memoryUsed := int64(0)
	memoryTotal := int64(0)

	lines := strings.Split(meminfoResp, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		lineLower := strings.ToLower(line)

		// Parse total hunk (total memory allocated by the game)
		// Format: "10485760 bytes total hunk"
		if strings.Contains(lineLower, "bytes total hunk") {
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				if total, err := strconv.ParseInt(fields[0], 10, 64); err == nil {
					memoryTotal = total
				} else {
					logger.Debug(fmt.Sprintf("Failed to parse total: %v", err))
				}
			}
		}

		// Parse total hunk in use (active memory usage)
		// Format: "66471 total hunk in use"
		if strings.Contains(lineLower, "total hunk in use") {
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				if used, err := strconv.ParseInt(fields[0], 10, 64); err == nil {
					memoryUsed = used
					logger.Debug(fmt.Sprintf("Parsed memoryUsed: %d bytes", memoryUsed))
				} else {
					logger.Debug(fmt.Sprintf("Failed to parse used: %v", err))
				}
			}
		}
	}

	cpuUsage := 0.0 // CPU usage not available in these commands

	return models.CreateSystemStats(cpuUsage, memoryUsed, memoryTotal)
}

func (sc *StatsCollector) collectPlayerStats() error {
	// Get status
	resp, err := sc.rcon.SendCommand("status")
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	// Parse player stats
	totalKills := 0
	totalDeaths := 0
	totalPing := 0
	totalScore := 0
	playerCount := 0

	lines := strings.Split(resp, "\n")
	for _, line := range lines {
		// Skip header lines
		if strings.Contains(line, "num") || strings.Contains(line, "---") || len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Parse player line (format: num score ping guid name lastmsg address)
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			if score, err := strconv.Atoi(fields[1]); err == nil {
				totalScore += score
				playerCount++
			}
			if ping, err := strconv.Atoi(fields[2]); err == nil {
				totalPing += ping
			}
		}
	}

	avgPing := 0.0
	avgScore := 0.0
	if playerCount > 0 {
		avgPing = float64(totalPing) / float64(playerCount)
		avgScore = float64(totalScore) / float64(playerCount)
	}

	return models.CreatePlayerStats(totalKills, totalDeaths, avgPing, avgScore)
}

// Helper functions
func parsePlayerCount(status string) int {
	count := 0
	lines := strings.Split(status, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines, headers, and separators
		if len(line) == 0 || strings.Contains(line, "num") || strings.Contains(line, "---") || strings.Contains(line, "map:") || strings.Contains(line, "players") {
			continue
		}
		// Player lines start with a number (slot ID)
		fields := strings.Fields(line)
		if len(fields) >= 7 {
			// Validate first field is a number (slot ID)
			if _, err := strconv.Atoi(fields[0]); err == nil {
				count++
			}
		}
	}
	return count
}

func parseCvarInt(response, cvar string, defaultValue int) int {
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Format: "cvar   value" (with variable spacing)
		if strings.HasPrefix(line, cvar) {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if val, err := strconv.Atoi(fields[1]); err == nil {
					return val
				}
			}
		}
	}
	return defaultValue
}

func parseCvarString(response, cvar, defaultValue string) string {
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Format: "cvar   value" (with variable spacing)
		if strings.HasPrefix(line, cvar) {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				// Join remaining fields in case value has spaces
				return strings.Join(fields[1:], " ")
			}
		}
	}
	return defaultValue
}

func parseUptime(response string) int {
	// Parse uptime from format: "uptime               5 hours"
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "uptime") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				// fields[1] is the number, fields[2] is the unit (hours, minutes, etc.)
				if uptimeVal, err := strconv.Atoi(fields[1]); err == nil {
					unit := fields[2]
					// Convert to seconds
					switch {
					case strings.HasPrefix(unit, "hour"):
						return uptimeVal * 3600
					case strings.HasPrefix(unit, "minute"):
						return uptimeVal * 60
					case strings.HasPrefix(unit, "second"):
						return uptimeVal
					case strings.HasPrefix(unit, "day"):
						return uptimeVal * 86400
					}
				}
			}
		}
	}
	return 0
}
