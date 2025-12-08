package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
)

// ServerStats stores server metrics over time for charting
type ServerStats struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Timestamp   time.Time `gorm:"index;not null" json:"timestamp"`
	PlayerCount int       `json:"playerCount"`
	MaxPlayers  int       `json:"maxPlayers"`
	MapName     string    `json:"mapName"`
	Gametype    string    `json:"gametype"`
	Hostname    string    `json:"hostname"`
	FPS         int       `json:"fps"`
	Uptime      int       `json:"uptime"` // in seconds
	CreatedAt   time.Time `json:"createdAt"`
}

// SystemStats stores system metrics over time
type SystemStats struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Timestamp   time.Time `gorm:"index;not null" json:"timestamp"`
	CPUUsage    float64   `json:"cpuUsage"`
	MemoryUsed  int64     `json:"memoryUsed"`  // in bytes
	MemoryTotal int64     `json:"memoryTotal"` // in bytes
	CreatedAt   time.Time `json:"createdAt"`
}

// PlayerStats stores individual player metrics over time
type PlayerStats struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Timestamp   time.Time `gorm:"index;not null" json:"timestamp"`
	TotalKills  int       `json:"totalKills"`
	TotalDeaths int       `json:"totalDeaths"`
	AvgPing     float64   `json:"avgPing"`
	AvgScore    float64   `json:"avgScore"`
	CreatedAt   time.Time `json:"createdAt"`
}

// CreateServerStats creates a new server stats entry
func CreateServerStats(playerCount, maxPlayers int, mapName, gametype, hostname string, fps, uptime int) error {
	db := database.DB
	stat := &ServerStats{
		Timestamp:   time.Now(),
		PlayerCount: playerCount,
		MaxPlayers:  maxPlayers,
		MapName:     mapName,
		Gametype:    gametype,
		Hostname:    hostname,
		FPS:         fps,
		Uptime:      uptime,
	}
	return db.Create(stat).Error
}

// CreateSystemStats creates a new system stats entry
func CreateSystemStats(cpuUsage float64, memoryUsed, memoryTotal int64) error {
	db := database.DB
	stat := &SystemStats{
		Timestamp:   time.Now(),
		CPUUsage:    cpuUsage,
		MemoryUsed:  memoryUsed,
		MemoryTotal: memoryTotal,
	}
	return db.Create(stat).Error
}

// CreatePlayerStats creates a new player stats entry
func CreatePlayerStats(totalKills, totalDeaths int, avgPing, avgScore float64) error {
	db := database.DB
	stat := &PlayerStats{
		Timestamp:   time.Now(),
		TotalKills:  totalKills,
		TotalDeaths: totalDeaths,
		AvgPing:     avgPing,
		AvgScore:    avgScore,
	}
	return db.Create(stat).Error
}

// GetServerStatsRange retrieves server stats within a time range
func GetServerStatsRange(start, end time.Time) ([]ServerStats, error) {
	db := database.DB
	var stats []ServerStats
	err := db.Where("timestamp BETWEEN ? AND ?", start, end).
		Order("timestamp ASC").
		Find(&stats).Error
	return stats, err
}

// GetSystemStatsRange retrieves system stats within a time range
func GetSystemStatsRange(start, end time.Time) ([]SystemStats, error) {
	db := database.DB
	var stats []SystemStats
	err := db.Where("timestamp BETWEEN ? AND ?", start, end).
		Order("timestamp ASC").
		Find(&stats).Error
	return stats, err
}

// GetPlayerStatsRange retrieves player stats within a time range
func GetPlayerStatsRange(start, end time.Time) ([]PlayerStats, error) {
	db := database.DB
	var stats []PlayerStats
	err := db.Where("timestamp BETWEEN ? AND ?", start, end).
		Order("timestamp ASC").
		Find(&stats).Error
	return stats, err
}

// CleanupOldStats removes stats older than the specified duration
func CleanupOldStats(olderThan time.Time) error {
	db := database.DB
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("timestamp < ?", olderThan).Delete(&ServerStats{}).Error; err != nil {
			return err
		}
		if err := tx.Where("timestamp < ?", olderThan).Delete(&SystemStats{}).Error; err != nil {
			return err
		}
		if err := tx.Where("timestamp < ?", olderThan).Delete(&PlayerStats{}).Error; err != nil {
			return err
		}
		return nil
	})
}
