package rest

import (
	"net/http"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

func RegisterIamGodRoute(r *gin.Engine, api *Api) {
	r.GET("/auth/iamgod", iamGod(api))
}

func iamGod(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if iamgod has already been used
		if models.HasSetting("iamgod_used", "true") {
			c.Set("error", "This route has already been used")
			c.Status(http.StatusForbidden)
			return
		}

		token, err := c.Cookie("session_token")
		if err != nil {
			c.Set("error", "Not authenticated")
			c.Status(http.StatusUnauthorized)
			return
		}

		session, err := models.GetSessionByToken(token)
		if err != nil {
			c.Set("error", "Invalid session")
			c.Status(http.StatusUnauthorized)
			return
		}

		user, err := models.GetUserByID(session.UserID)
		if err != nil {
			c.Set("error", "User not found")
			c.Status(http.StatusInternalServerError)
			return
		}

		// Get the super_admin role (should already exist from startup initialization)
		superAdminRole, err := models.GetRoleByName("super_admin")
		if err != nil {
			c.Set("error", "Super admin role not found")
			c.Status(http.StatusInternalServerError)
			return
		}

		// Assign super admin role to user
		if err := models.AddRoleToUser(user.ID, superAdminRole.ID); err != nil {
			c.Set("error", "Failed to assign super admin role")
			c.Status(http.StatusInternalServerError)
			return
		}

		// Mark as used in database - persist across restarts
		if err := models.SetSetting("iamgod_used", "true"); err != nil {
			c.Set("error", "Failed to mark iamgod as used")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{
			"message": "Super admin privileges granted successfully",
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
			},
		})
		c.Status(http.StatusOK)
	}
}
