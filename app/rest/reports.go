package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

type ActionReportRequest struct {
	Action   string `json:"action" binding:"required"` // "dismiss", "ban", "tempban"
	Duration *int   `json:"duration"`                  // Duration in hours for tempban
	Reason   string `json:"reason"`                    // Additional reason/notes
}

func RegisterReportRoutes(r *gin.Engine, api *Api) {
	reports := r.Group("/reports")
	reports.Use(AuthMiddleware())
	{
		reports.GET("", RequirePermission("reports.view"), getAllReports(api))
		reports.GET("/pending", RequirePermission("reports.view"), getPendingReports(api))
		reports.GET("/:id", RequirePermission("reports.view"), getReport(api))
		reports.POST("/:id/action", RequirePermission("reports.action"), actionReport(api))
		reports.DELETE("/:id", RequirePermission("reports.action"), deleteReport(api))
	}

	tempBans := r.Group("/tempbans")
	tempBans.Use(AuthMiddleware())
	{
		tempBans.GET("", RequirePermission("reports.view"), getAllTempBans(api))
		tempBans.GET("/active", RequirePermission("reports.view"), getActiveTempBans(api))
		tempBans.POST("/:id/revoke", RequirePermission("reports.action"), revokeTempBan(api))
	}
}

func getAllReports(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		reports, err := models.GetAllReports()
		if err != nil {
			c.Set("error", "Failed to retrieve reports")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", reports)
		c.Status(http.StatusOK)
	}
}

func getPendingReports(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		reports, err := models.GetPendingReports()
		if err != nil {
			c.Set("error", "Failed to retrieve pending reports")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", reports)
		c.Status(http.StatusOK)
	}
}

func getReport(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid report ID")
			c.Status(http.StatusBadRequest)
			return
		}

		report, err := models.GetReportByID(uint(id))
		if err != nil {
			c.Set("error", "Report not found")
			c.Status(http.StatusNotFound)
			return
		}

		c.Set("data", report)
		c.Status(http.StatusOK)
	}
}

func actionReport(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid report ID")
			c.Status(http.StatusBadRequest)
			return
		}

		var req ActionReportRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		report, err := models.GetReportByID(uint(id))
		if err != nil {
			c.Set("error", "Report not found")
			c.Status(http.StatusNotFound)
			return
		}

		// Get current user ID
		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)
		uid := user.ID

		updates := make(map[string]interface{})
		updates["reviewed_by_user_id"] = uid

		switch req.Action {
		case "dismiss":
			updates["status"] = "dismissed"
			updates["action_taken"] = req.Reason

		case "ban":
			updates["status"] = "actioned"
			updates["action_taken"] = "Permanently banned: " + req.Reason
			// Execute permanent ban via RCON
			_, err := api.rcon.SendCommand("banClient " + report.ReportedGUID + " " + req.Reason)
			if err != nil {
				Audit.LogBan(c, report.ReportedName, report.ReportedGUID, req.Reason, false, err.Error())
				c.Set("error", "Failed to execute ban")
				c.Status(http.StatusInternalServerError)
				return
			}
			Audit.LogBan(c, report.ReportedName, report.ReportedGUID, req.Reason, true, "")

		case "tempban":
			if req.Duration == nil || *req.Duration <= 0 {
				c.Set("error", "Duration required for temporary ban")
				c.Status(http.StatusBadRequest)
				return
			}

			duration := time.Duration(*req.Duration) * time.Hour
			_, err := models.CreateTempBan(report.ReportedName, report.ReportedGUID, req.Reason, duration, &uid)
			if err != nil {
				Audit.LogTempBan(c, report.ReportedName, report.ReportedGUID, req.Reason, *req.Duration, false, err.Error())
				c.Set("error", "Failed to create temporary ban")
				c.Status(http.StatusInternalServerError)
				return
			}

			updates["status"] = "actioned"
			updates["action_taken"] = "Temporarily banned for " + strconv.Itoa(*req.Duration) + " hours: " + req.Reason

			// Kick the player
			_, err = api.rcon.SendCommand("clientkick " + report.ReportedGUID + " Temporarily banned: " + req.Reason)
			if err != nil {
				// Log error but continue - they might already be offline
			}

			Audit.LogTempBan(c, report.ReportedName, report.ReportedGUID, req.Reason, *req.Duration, true, "")

		default:
			c.Set("error", "Invalid action")
			c.Status(http.StatusBadRequest)
			return
		}

		err = models.UpdateReport(uint(id), updates)
		if err != nil {
			c.Set("error", "Failed to update report")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Report actioned successfully"})
		c.Status(http.StatusOK)
	}
}

func deleteReport(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid report ID")
			c.Status(http.StatusBadRequest)
			return
		}

		err = models.DeleteReport(uint(id))
		if err != nil {
			c.Set("error", "Failed to delete report")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Report deleted successfully"})
		c.Status(http.StatusOK)
	}
}

func getAllTempBans(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		bans, err := models.GetAllTempBans()
		if err != nil {
			c.Set("error", "Failed to retrieve temp bans")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", bans)
		c.Status(http.StatusOK)
	}
}

func getActiveTempBans(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		bans, err := models.GetActiveTempBans()
		if err != nil {
			c.Set("error", "Failed to retrieve active temp bans")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", bans)
		c.Status(http.StatusOK)
	}
}

func revokeTempBan(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid temp ban ID")
			c.Status(http.StatusBadRequest)
			return
		}

		err = models.RevokeTempBan(uint(id))
		if err != nil {
			c.Set("error", "Failed to revoke temp ban")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Temp ban revoked successfully"})
		c.Status(http.StatusOK)
	}
}
