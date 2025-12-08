package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
)

// WebhookEvent represents the type of event that triggers a webhook
type WebhookEvent string

const (
	WebhookEventPlayerBanned   WebhookEvent = "player.banned"
	WebhookEventPlayerUnbanned WebhookEvent = "player.unbanned"
	WebhookEventPlayerKicked   WebhookEvent = "player.kicked"
	WebhookEventReportCreated  WebhookEvent = "report.created"
	WebhookEventReportActioned WebhookEvent = "report.actioned"
	WebhookEventUserApproved   WebhookEvent = "user.approved"
	WebhookEventUserRejected   WebhookEvent = "user.rejected"
	WebhookEventServerOnline   WebhookEvent = "server.online"
	WebhookEventServerOffline  WebhookEvent = "server.offline"
	WebhookEventSecurityAlert  WebhookEvent = "security.alert"
)

// Webhook represents a webhook configuration
type Webhook struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	URL         string `gorm:"type:varchar(500);not null" json:"url"`
	Secret      string `gorm:"type:varchar(100)" json:"-"`       // HMAC secret for signing
	Events      string `gorm:"type:text;not null" json:"events"` // JSON array of WebhookEvent
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	Description string `gorm:"type:text" json:"description"`

	// Retry configuration
	MaxRetries     int `gorm:"default:3" json:"maxRetries"`
	RetryDelay     int `gorm:"default:60" json:"retryDelay"` // Seconds
	TimeoutSeconds int `gorm:"default:10" json:"timeoutSeconds"`

	// Statistics
	TotalDeliveries   int        `gorm:"default:0" json:"totalDeliveries"`
	FailedDeliveries  int        `gorm:"default:0" json:"failedDeliveries"`
	LastDeliveryAt    *time.Time `json:"lastDeliveryAt,omitempty"`
	LastFailureAt     *time.Time `json:"lastFailureAt,omitempty"`
	LastFailureReason string     `gorm:"type:text" json:"lastFailureReason,omitempty"`
}

// WebhookDelivery tracks webhook delivery attempts
type WebhookDelivery struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	WebhookID uint     `gorm:"not null;index" json:"webhookId"`
	Webhook   *Webhook `gorm:"foreignKey:WebhookID;constraint:OnDelete:CASCADE" json:"webhook,omitempty"`

	Event   WebhookEvent `gorm:"type:varchar(50);not null;index" json:"event"`
	Payload string       `gorm:"type:text;not null" json:"payload"` // JSON payload

	// Delivery status
	Status       string     `gorm:"type:varchar(20);not null;index" json:"status"` // pending, delivered, failed
	AttemptCount int        `gorm:"default:0" json:"attemptCount"`
	ResponseCode *int       `json:"responseCode,omitempty"`
	ResponseBody string     `gorm:"type:text" json:"responseBody,omitempty"`
	ErrorMessage string     `gorm:"type:text" json:"errorMessage,omitempty"`
	DeliveredAt  *time.Time `json:"deliveredAt,omitempty"`
	NextRetryAt  *time.Time `json:"nextRetryAt,omitempty"`
}

// CreateWebhook creates a new webhook
func CreateWebhook(webhook *Webhook) error {
	return database.DB.Create(webhook).Error
}

// GetWebhooks retrieves all webhooks
func GetWebhooks() ([]Webhook, error) {
	var webhooks []Webhook
	err := database.DB.Order("created_at DESC").Find(&webhooks).Error
	return webhooks, err
}

// GetWebhookByID retrieves a webhook by ID
func GetWebhookByID(id uint) (*Webhook, error) {
	var webhook Webhook
	err := database.DB.First(&webhook, id).Error
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// UpdateWebhook updates a webhook
func UpdateWebhook(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&Webhook{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteWebhook soft deletes a webhook
func DeleteWebhook(id uint) error {
	return database.DB.Delete(&Webhook{}, id).Error
}

// GetEnabledWebhooksForEvent retrieves enabled webhooks for a specific event
func GetEnabledWebhooksForEvent(event WebhookEvent) ([]Webhook, error) {
	var webhooks []Webhook
	err := database.DB.Where("enabled = ? AND events LIKE ?", true, "%"+string(event)+"%").Find(&webhooks).Error
	return webhooks, err
}

// CreateWebhookDelivery creates a webhook delivery record
func CreateWebhookDelivery(delivery *WebhookDelivery) error {
	return database.DB.Create(delivery).Error
}

// GetWebhookDeliveries retrieves delivery history for a webhook
func GetWebhookDeliveries(webhookID uint, limit int) ([]WebhookDelivery, error) {
	var deliveries []WebhookDelivery
	query := database.DB.Where("webhook_id = ?", webhookID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&deliveries).Error
	return deliveries, err
}

// GetPendingDeliveries retrieves deliveries pending retry
func GetPendingDeliveries() ([]WebhookDelivery, error) {
	var deliveries []WebhookDelivery
	now := time.Now()
	err := database.DB.Preload("Webhook").
		Where("status = ? AND (next_retry_at IS NULL OR next_retry_at <= ?)", "pending", now).
		Find(&deliveries).Error
	return deliveries, err
}

// UpdateWebhookDelivery updates a delivery record
func UpdateWebhookDelivery(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&WebhookDelivery{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateWebhookStats updates webhook statistics
func UpdateWebhookStats(webhookID uint, success bool, failureReason string) error {
	updates := map[string]interface{}{
		"total_deliveries": gorm.Expr("total_deliveries + 1"),
		"last_delivery_at": time.Now(),
	}

	if !success {
		updates["failed_deliveries"] = gorm.Expr("failed_deliveries + 1")
		updates["last_failure_at"] = time.Now()
		updates["last_failure_reason"] = failureReason
	}

	return database.DB.Model(&Webhook{}).Where("id = ?", webhookID).Updates(updates).Error
}

// TableName specifies the table name for Webhook
func (Webhook) TableName() string {
	return "webhooks"
}

// TableName specifies the table name for WebhookDelivery
func (WebhookDelivery) TableName() string {
	return "webhook_deliveries"
}
