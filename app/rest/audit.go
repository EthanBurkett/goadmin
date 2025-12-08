package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

func RegisterAuditRoutes(r *gin.Engine, api *Api) {
	audit := r.Group("/audit")
	audit.Use(AuthMiddleware())
	audit.Use(RequirePermission("rbac.manage")) // Only admins can view audit logs
	{
		audit.GET("/logs", getAuditLogs(api))
		audit.GET("/logs/recent", getRecentAuditLogs(api))
		audit.GET("/logs/user/:userId", getAuditLogsByUser(api))
		audit.GET("/logs/action/:action", getAuditLogsByAction(api))
	}
}

type AuditLogsQuery struct {
	UserID     *uint   `form:"user_id"`
	Action     *string `form:"action"`
	Source     *string `form:"source"`
	Success    *bool   `form:"success"`
	TargetType *string `form:"target_type"`
	TargetID   *string `form:"target_id"`
	StartDate  *string `form:"start_date"`
	EndDate    *string `form:"end_date"`
	Limit      int     `form:"limit"`
	Offset     int     `form:"offset"`
}

func getAuditLogs(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query AuditLogsQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			c.Set("error", "Invalid query parameters")
			c.Status(http.StatusBadRequest)
			return
		}

		// Set defaults
		if query.Limit == 0 {
			query.Limit = 50
		}
		if query.Limit > 500 {
			query.Limit = 500
		}

		// Build filters
		filters := make(map[string]interface{})
		if query.UserID != nil {
			filters["user_id"] = *query.UserID
		}
		if query.Action != nil {
			filters["action"] = *query.Action
		}
		if query.Source != nil {
			filters["source"] = *query.Source
		}
		if query.Success != nil {
			filters["success"] = *query.Success
		}
		if query.TargetType != nil {
			filters["target_type"] = *query.TargetType
		}
		if query.TargetID != nil {
			filters["target_id"] = *query.TargetID
		}
		if query.StartDate != nil {
			startDate, err := time.Parse(time.RFC3339, *query.StartDate)
			if err == nil {
				filters["start_date"] = startDate
			}
		}
		if query.EndDate != nil {
			endDate, err := time.Parse(time.RFC3339, *query.EndDate)
			if err == nil {
				filters["end_date"] = endDate
			}
		}

		logs, total, err := models.GetAuditLogs(api.DB, filters, query.Limit, query.Offset)
		if err != nil {
			c.Set("error", "Failed to fetch audit logs")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{
			"logs":   logs,
			"total":  total,
			"limit":  query.Limit,
			"offset": query.Offset,
		})
		c.Status(http.StatusOK)
	}
}

func getRecentAuditLogs(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "100")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 100
		}
		if limit > 500 {
			limit = 500
		}

		logs, err := models.GetRecentAuditLogs(api.DB, limit)
		if err != nil {
			c.Set("error", "Failed to fetch recent audit logs")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{
			"logs": logs,
		})
		c.Status(http.StatusOK)
	}
}

func getAuditLogsByUser(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := strconv.ParseUint(c.Param("userId"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid user ID")
			c.Status(http.StatusBadRequest)
			return
		}

		limitStr := c.DefaultQuery("limit", "100")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 100
		}
		if limit > 500 {
			limit = 500
		}

		logs, err := models.GetAuditLogsByUser(api.DB, uint(userId), limit)
		if err != nil {
			c.Set("error", "Failed to fetch audit logs for user")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{
			"logs": logs,
		})
		c.Status(http.StatusOK)
	}
}

func getAuditLogsByAction(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		action := c.Param("action")
		if action == "" {
			c.Set("error", "Action parameter required")
			c.Status(http.StatusBadRequest)
			return
		}

		limitStr := c.DefaultQuery("limit", "100")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 100
		}
		if limit > 500 {
			limit = 500
		}

		logs, err := models.GetAuditLogsByAction(api.DB, models.ActionType(action), limit)
		if err != nil {
			c.Set("error", "Failed to fetch audit logs for action")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{
			"logs": logs,
		})
		c.Status(http.StatusOK)
	}
}
