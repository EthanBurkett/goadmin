package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
)

// Report represents a player report
type Report struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	ReporterName     string    `gorm:"not null" json:"reporterName"`     // Name of player who made report
	ReporterGUID     string    `gorm:"not null" json:"reporterGuid"`     // GUID of reporter
	ReportedName     string    `gorm:"not null" json:"reportedName"`     // Name of reported player
	ReportedGUID     string    `gorm:"not null" json:"reportedGuid"`     // GUID of reported player
	Reason           string    `gorm:"type:text;not null" json:"reason"` // Reason for report
	Status           string    `gorm:"default:'pending'" json:"status"`  // pending, reviewed, actioned, dismissed
	ActionTaken      string    `gorm:"type:text" json:"actionTaken"`     // What action was taken (if any)
	ReviewedByUserID *uint     `gorm:"index" json:"reviewedByUserId"`    // User who reviewed
	ReviewedBy       *User     `gorm:"foreignKey:ReviewedByUserID;constraint:OnDelete:SET NULL" json:"reviewedBy,omitempty"`
	ServerID         *uint     `gorm:"index" json:"serverId,omitempty"` // Server where report was created
	Server           *Server   `gorm:"foreignKey:ServerID;constraint:OnDelete:SET NULL" json:"server,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// CreateReport creates a new player report
func CreateReport(reporterName, reporterGUID, reportedName, reportedGUID, reason string, serverID *uint) (*Report, error) {
	db := database.DB

	report := &Report{
		ReporterName: reporterName,
		ReporterGUID: reporterGUID,
		ReportedName: reportedName,
		ReportedGUID: reportedGUID,
		Reason:       reason,
		ServerID:     serverID,
		Status:       "pending",
	}

	err := db.Create(report).Error
	return report, err
}

// GetAllReports gets all reports, optionally filtered by server ID
func GetAllReports(serverID *uint) ([]Report, error) {
	db := database.DB
	var reports []Report
	query := db.Preload("ReviewedBy").Preload("Server")

	if serverID != nil {
		query = query.Where("server_id = ?", *serverID)
	}

	err := query.Order("created_at DESC").Find(&reports).Error
	return reports, err
}

// GetPendingReports gets all pending reports, optionally filtered by server ID
func GetPendingReports(serverID *uint) ([]Report, error) {
	db := database.DB
	var reports []Report
	query := db.Preload("ReviewedBy").Preload("Server").Where("status = ?", "pending")

	if serverID != nil {
		query = query.Where("server_id = ?", *serverID)
	}

	err := query.Order("created_at DESC").Find(&reports).Error
	return reports, err
}

// GetReportByID gets a report by ID
func GetReportByID(id uint) (*Report, error) {
	db := database.DB
	var report Report
	err := db.Preload("ReviewedBy").First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// UpdateReport updates a report
func UpdateReport(id uint, updates map[string]interface{}) error {
	db := database.DB
	return db.Model(&Report{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteReport deletes a report
func DeleteReport(id uint) error {
	db := database.DB
	return db.Delete(&Report{}, id).Error
}

// TableName specifies the table name
func (Report) TableName() string {
	return "reports"
}
