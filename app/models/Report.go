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
	ReviewedByUserID *uint     `json:"reviewedByUserId"`                 // User who reviewed
	ReviewedBy       *User     `gorm:"foreignKey:ReviewedByUserID" json:"reviewedBy,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// CreateReport creates a new player report
func CreateReport(reporterName, reporterGUID, reportedName, reportedGUID, reason string) (*Report, error) {
	db := database.DB

	report := &Report{
		ReporterName: reporterName,
		ReporterGUID: reporterGUID,
		ReportedName: reportedName,
		ReportedGUID: reportedGUID,
		Reason:       reason,
		Status:       "pending",
	}

	err := db.Create(report).Error
	return report, err
}

// GetAllReports gets all reports
func GetAllReports() ([]Report, error) {
	db := database.DB
	var reports []Report
	err := db.Preload("ReviewedBy").Order("created_at DESC").Find(&reports).Error
	return reports, err
}

// GetPendingReports gets all pending reports
func GetPendingReports() ([]Report, error) {
	db := database.DB
	var reports []Report
	err := db.Where("status = ?", "pending").Order("created_at DESC").Find(&reports).Error
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
