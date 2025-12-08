package rest

import (
	"fmt"
	"net/http"

	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			c.Set("error", "Not authenticated")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, err := models.GetSessionByToken(token)
		if err != nil {
			c.Set("error", "Invalid session")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Reload user with roles and permissions
		user, err := models.GetUserByID(session.UserID)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get user by ID %d: %v", session.UserID, err))
			c.Set("error", "User not found")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)
		c.Set("session", session)
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.Set("error", "Not authenticated")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.Set("error", "Invalid user")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !user.HasPermission(permission) {
			c.Set("error", "Insufficient permissions")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.Set("error", "Not authenticated")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			c.Set("error", "Invalid user")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !user.HasRole(role) {
			c.Set("error", "Insufficient role")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
