package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/ethanburkett/goadmin/app/webhook"
	"github.com/gin-gonic/gin"
)

// WebhookRequest represents the request to create/update a webhook
type WebhookRequest struct {
	Name           string   `json:"name" binding:"required"`
	URL            string   `json:"url" binding:"required,url"`
	Secret         string   `json:"secret"`
	Events         []string `json:"events" binding:"required"`
	Description    string   `json:"description"`
	MaxRetries     *int     `json:"maxRetries"`
	RetryDelay     *int     `json:"retryDelay"`
	TimeoutSeconds *int     `json:"timeoutSeconds"`
	Enabled        *bool    `json:"enabled"`
}

// RegisterWebhookRoutes registers webhook routes
func RegisterWebhookRoutes(router *gin.Engine, api *Api) {
	webhooks := router.Group("/webhooks")
	webhooks.Use(AuthMiddleware())
	webhooks.Use(RequirePermission("rbac.manage")) // Only admins can manage webhooks

	webhooks.GET("", getWebhooks(api))
	webhooks.POST("", createWebhook(api))
	webhooks.GET("/:id", getWebhook(api))
	webhooks.PUT("/:id", updateWebhook(api))
	webhooks.DELETE("/:id", deleteWebhook(api))
	webhooks.GET("/:id/deliveries", getWebhookDeliveries(api))
	webhooks.POST("/:id/test", testWebhook(api))
}

// getWebhooks retrieves all webhooks
func getWebhooks(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		webhooks, err := models.GetWebhooks()
		if err != nil {
			c.Set("error", "Failed to retrieve webhooks")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"webhooks": webhooks})
		c.Status(http.StatusOK)
	}
}

// createWebhook creates a new webhook
func createWebhook(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req WebhookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		// Serialize events to JSON
		eventsJSON, err := json.Marshal(req.Events)
		if err != nil {
			c.Set("error", "Failed to serialize events")
			c.Status(http.StatusBadRequest)
			return
		}

		webhook := &models.Webhook{
			Name:           req.Name,
			URL:            req.URL,
			Secret:         req.Secret,
			Events:         string(eventsJSON),
			Description:    req.Description,
			MaxRetries:     3,
			RetryDelay:     60,
			TimeoutSeconds: 10,
			Enabled:        true,
		}

		if req.MaxRetries != nil {
			webhook.MaxRetries = *req.MaxRetries
		}
		if req.RetryDelay != nil {
			webhook.RetryDelay = *req.RetryDelay
		}
		if req.TimeoutSeconds != nil {
			webhook.TimeoutSeconds = *req.TimeoutSeconds
		}
		if req.Enabled != nil {
			webhook.Enabled = *req.Enabled
		}

		if err := models.CreateWebhook(webhook); err != nil {
			c.Set("error", "Failed to create webhook")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"webhook": webhook})
		c.Status(http.StatusCreated)
	}
}

// getWebhook retrieves a single webhook
func getWebhook(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid webhook ID")
			c.Status(http.StatusBadRequest)
			return
		}

		webhook, err := models.GetWebhookByID(uint(id))
		if err != nil {
			c.Set("error", "Webhook not found")
			c.Status(http.StatusNotFound)
			return
		}

		c.Set("data", gin.H{"webhook": webhook})
		c.Status(http.StatusOK)
	}
}

// updateWebhook updates a webhook
func updateWebhook(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid webhook ID")
			c.Status(http.StatusBadRequest)
			return
		}

		var req WebhookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		updates := make(map[string]interface{})
		updates["name"] = req.Name
		updates["url"] = req.URL
		updates["description"] = req.Description

		if req.Secret != "" {
			updates["secret"] = req.Secret
		}

		// Serialize events to JSON
		eventsJSON, err := json.Marshal(req.Events)
		if err != nil {
			c.Set("error", "Failed to serialize events")
			c.Status(http.StatusBadRequest)
			return
		}
		updates["events"] = string(eventsJSON)

		if req.MaxRetries != nil {
			updates["max_retries"] = *req.MaxRetries
		}
		if req.RetryDelay != nil {
			updates["retry_delay"] = *req.RetryDelay
		}
		if req.TimeoutSeconds != nil {
			updates["timeout_seconds"] = *req.TimeoutSeconds
		}
		if req.Enabled != nil {
			updates["enabled"] = *req.Enabled
		}

		if err := models.UpdateWebhook(uint(id), updates); err != nil {
			c.Set("error", "Failed to update webhook")
			c.Status(http.StatusInternalServerError)
			return
		}

		webhook, _ := models.GetWebhookByID(uint(id))
		c.Set("data", gin.H{"webhook": webhook})
		c.Status(http.StatusOK)
	}
}

// deleteWebhook deletes a webhook
func deleteWebhook(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid webhook ID")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.DeleteWebhook(uint(id)); err != nil {
			c.Set("error", "Failed to delete webhook")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Webhook deleted successfully"})
		c.Status(http.StatusOK)
	}
}

// getWebhookDeliveries retrieves delivery history for a webhook
func getWebhookDeliveries(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid webhook ID")
			c.Status(http.StatusBadRequest)
			return
		}

		limit := 50
		if limitStr := c.Query("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
				if limit > 500 {
					limit = 500
				}
			}
		}

		deliveries, err := models.GetWebhookDeliveries(uint(id), limit)
		if err != nil {
			c.Set("error", "Failed to retrieve deliveries")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"deliveries": deliveries})
		c.Status(http.StatusOK)
	}
}

// testWebhook sends a test webhook delivery
func testWebhook(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid webhook ID")
			c.Status(http.StatusBadRequest)
			return
		}

		wh, err := models.GetWebhookByID(uint(id))
		if err != nil {
			c.Set("error", "Webhook not found")
			c.Status(http.StatusNotFound)
			return
		}

		// Send test event
		testData := map[string]interface{}{
			"test":         true,
			"message":      "This is a test webhook delivery from GoAdmin",
			"webhook_id":   wh.ID,
			"webhook_name": wh.Name,
		}

		go webhook.GlobalDispatcher.Dispatch(models.WebhookEventSecurityAlert, testData)

		c.Set("data", gin.H{"message": "Test webhook queued for delivery"})
		c.Status(http.StatusOK)
	}
}
