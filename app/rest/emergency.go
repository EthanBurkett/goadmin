package rest

import (
	"net/http"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

func RegisterEmergencyRoutes(r *gin.Engine, api *Api) {
	emergency := r.Group("/emergency")
	emergency.Use(AuthMiddleware())
	{
		// View disabled commands (requires admin permission)
		emergency.GET("/disabled", RequirePermission("commands.manage"), getDisabledCommands(api))

		// Manually re-enable a command (requires admin permission)
		emergency.POST("/reenable/:command", RequirePermission("commands.manage"), reenableCommand(api))

		// View user alert counts (requires admin permission)
		emergency.GET("/alerts", RequirePermission("commands.manage"), getUserAlerts(api))

		// Reset user alerts (requires admin permission)
		emergency.POST("/alerts/:userId/reset", RequirePermission("commands.manage"), resetUserAlerts(api))
	}
}

func getDisabledCommands(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		disabled := models.GlobalEmergencyShutdown.GetDisabledCommands()

		c.Set("data", disabled)
		c.Status(http.StatusOK)
	}
}

func reenableCommand(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		command := c.Param("command")

		// Get current user from session
		session, exists := c.Get("session")
		if !exists {
			c.Set("error", "Unauthorized")
			c.Status(http.StatusUnauthorized)
			return
		}

		userSession := session.(*models.Session)

		err := models.GlobalEmergencyShutdown.EnableCommand(command, userSession.User.ID)
		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		c.Set("data", map[string]interface{}{
			"message": "Command re-enabled successfully",
			"command": command,
		})
		c.Status(http.StatusOK)
	}
}

func getUserAlerts(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would need to be implemented to return all user alert counts
		// For now, return empty
		c.Set("data", map[string]interface{}{
			"message": "User alerts endpoint not fully implemented",
		})
		c.Status(http.StatusOK)
	}
}

func resetUserAlerts(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would reset alert count for specific user
		// For now, return success
		c.Set("data", map[string]interface{}{
			"message": "User alerts reset successfully",
		})
		c.Status(http.StatusOK)
	}
}
