package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
)

// Dispatcher handles webhook delivery
type Dispatcher struct {
	client *http.Client
}

// NewDispatcher creates a new webhook dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WebhookPayload represents the standard webhook payload structure
type WebhookPayload struct {
	Event     string                 `json:"event"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Dispatch sends a webhook for a specific event
func (d *Dispatcher) Dispatch(event models.WebhookEvent, data map[string]interface{}) error {
	webhooks, err := models.GetEnabledWebhooksForEvent(event)
	if err != nil {
		return fmt.Errorf("failed to get webhooks: %w", err)
	}

	if len(webhooks) == 0 {
		logger.Debug(fmt.Sprintf("No webhooks configured for event: %s", event))
		return nil
	}

	payload := WebhookPayload{
		Event:     string(event),
		Timestamp: time.Now(),
		Data:      data,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	for _, webhook := range webhooks {
		delivery := &models.WebhookDelivery{
			WebhookID:    webhook.ID,
			Event:        event,
			Payload:      string(payloadBytes),
			Status:       "pending",
			AttemptCount: 0,
		}

		if err := models.CreateWebhookDelivery(delivery); err != nil {
			logger.Error(fmt.Sprintf("Failed to create webhook delivery: %v", err))
			continue
		}

		// Attempt immediate delivery in background
		go d.deliverWebhook(&webhook, delivery, payloadBytes)
	}

	return nil
}

// deliverWebhook attempts to deliver a webhook
func (d *Dispatcher) deliverWebhook(webhook *models.Webhook, delivery *models.WebhookDelivery, payloadBytes []byte) {
	delivery.AttemptCount++

	// Create HTTP request
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		d.markDeliveryFailed(webhook, delivery, 0, "", fmt.Sprintf("Failed to create request: %v", err))
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoAdmin-Webhook/1.0")
	req.Header.Set("X-Webhook-Event", string(delivery.Event))
	req.Header.Set("X-Webhook-Delivery", fmt.Sprintf("%d", delivery.ID))
	req.Header.Set("X-Webhook-Timestamp", delivery.CreatedAt.Format(time.RFC3339))

	// Sign payload if secret is configured
	if webhook.Secret != "" {
		signature := d.signPayload(payloadBytes, webhook.Secret)
		req.Header.Set("X-Webhook-Signature", signature)
	}

	// Use custom timeout if configured
	client := d.client
	if webhook.TimeoutSeconds > 0 {
		client = &http.Client{
			Timeout: time.Duration(webhook.TimeoutSeconds) * time.Second,
		}
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		d.markDeliveryFailed(webhook, delivery, 0, "", fmt.Sprintf("Request failed: %v", err))
		d.scheduleRetry(webhook, delivery)
		return
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, _ := io.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	// Check response status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		d.markDeliverySuccess(webhook, delivery, resp.StatusCode, responseBody)
	} else {
		d.markDeliveryFailed(webhook, delivery, resp.StatusCode, responseBody,
			fmt.Sprintf("HTTP %d", resp.StatusCode))
		d.scheduleRetry(webhook, delivery)
	}
}

// markDeliverySuccess marks a delivery as successful
func (d *Dispatcher) markDeliverySuccess(webhook *models.Webhook, delivery *models.WebhookDelivery, statusCode int, responseBody string) {
	now := time.Now()
	updates := map[string]interface{}{
		"status":        "delivered",
		"response_code": statusCode,
		"response_body": responseBody,
		"delivered_at":  now,
	}

	if err := models.UpdateWebhookDelivery(delivery.ID, updates); err != nil {
		logger.Error(fmt.Sprintf("Failed to update delivery: %v", err))
	}

	if err := models.UpdateWebhookStats(webhook.ID, true, ""); err != nil {
		logger.Error(fmt.Sprintf("Failed to update webhook stats: %v", err))
	}

	logger.Info(fmt.Sprintf("Webhook delivered successfully: %s to %s (status: %d)",
		delivery.Event, webhook.URL, statusCode))
}

// markDeliveryFailed marks a delivery as failed
func (d *Dispatcher) markDeliveryFailed(webhook *models.Webhook, delivery *models.WebhookDelivery,
	statusCode int, responseBody string, errorMsg string) {

	updates := map[string]interface{}{
		"error_message": errorMsg,
		"attempt_count": delivery.AttemptCount,
	}

	if statusCode > 0 {
		updates["response_code"] = statusCode
		updates["response_body"] = responseBody
	}

	// Mark as failed if max retries exceeded
	if delivery.AttemptCount >= webhook.MaxRetries {
		updates["status"] = "failed"
		logger.Error(fmt.Sprintf("Webhook delivery failed after %d attempts: %s to %s - %s",
			delivery.AttemptCount, delivery.Event, webhook.URL, errorMsg))
	}

	if err := models.UpdateWebhookDelivery(delivery.ID, updates); err != nil {
		logger.Error(fmt.Sprintf("Failed to update delivery: %v", err))
	}

	if delivery.AttemptCount >= webhook.MaxRetries {
		if err := models.UpdateWebhookStats(webhook.ID, false, errorMsg); err != nil {
			logger.Error(fmt.Sprintf("Failed to update webhook stats: %v", err))
		}
	}
}

// scheduleRetry schedules a delivery retry with exponential backoff
func (d *Dispatcher) scheduleRetry(webhook *models.Webhook, delivery *models.WebhookDelivery) {
	if delivery.AttemptCount >= webhook.MaxRetries {
		return
	}

	// Exponential backoff: retryDelay * 2^(attemptCount-1)
	delay := time.Duration(webhook.RetryDelay) * time.Second
	backoffMultiplier := 1 << (delivery.AttemptCount - 1) // 2^(n-1)
	nextRetry := time.Now().Add(delay * time.Duration(backoffMultiplier))

	updates := map[string]interface{}{
		"next_retry_at": nextRetry,
	}

	if err := models.UpdateWebhookDelivery(delivery.ID, updates); err != nil {
		logger.Error(fmt.Sprintf("Failed to schedule retry: %v", err))
	}

	logger.Info(fmt.Sprintf("Scheduled webhook retry #%d for %s at %s",
		delivery.AttemptCount+1, webhook.URL, nextRetry.Format(time.RFC3339)))
}

// signPayload creates HMAC SHA256 signature
func (d *Dispatcher) signPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// ProcessRetries processes pending webhook deliveries
func (d *Dispatcher) ProcessRetries() error {
	deliveries, err := models.GetPendingDeliveries()
	if err != nil {
		return fmt.Errorf("failed to get pending deliveries: %w", err)
	}

	for _, delivery := range deliveries {
		if delivery.Webhook == nil {
			continue
		}

		payloadBytes := []byte(delivery.Payload)
		go d.deliverWebhook(delivery.Webhook, &delivery, payloadBytes)
	}

	return nil
}

// StartRetryWorker starts a background worker to process retries
func (d *Dispatcher) StartRetryWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			if err := d.ProcessRetries(); err != nil {
				logger.Error(fmt.Sprintf("Failed to process webhook retries: %v", err))
			}
		}
	}()
}

// Global dispatcher instance
var GlobalDispatcher = NewDispatcher()
