package rest

import (
	"net/http"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterAuthRoutes(r *gin.Engine, api *Api) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", RateLimitByIP(LoginRateLimiter), login(api))
		auth.POST("/register", RateLimitByIP(LoginRateLimiter), register(api))
		auth.POST("/logout", logout(api))
		auth.GET("/me", getMe(api))
	}
}

func login(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := models.GetUserByUsername(req.Username)
		if err != nil || !user.CheckPassword(req.Password) {
			c.Set("error", "Invalid username or password. Please check your credentials and try again.")
			c.Status(http.StatusUnauthorized)
			return
		}

		if !user.Approved {
			c.Set("error", "Your account is pending approval. Please wait for an administrator to approve your registration.")
			c.Status(http.StatusForbidden)
			return
		}

		session, err := models.CreateSession(user.ID)
		if err != nil {
			c.Set("error", "Failed to create session")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.SetCookie("session_token", session.Token, 30*24*60*60, "/", "", false, true)

		c.Set("data", gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
			},
		})
		c.Status(http.StatusOK)
	}
}

func register(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		if _, err := models.GetUserByUsername(req.Username); err == nil {
			c.Set("error", "This username is already taken. Please choose a different username.")
			c.Status(http.StatusConflict)
			return
		}

		user, err := models.CreateUser(req.Username, req.Password)
		if err != nil {
			c.Set("error", "Failed to create your account. Please try again or contact support.")
			c.Status(http.StatusInternalServerError)
			return
		}

		session, err := models.CreateSession(user.ID)
		if err != nil {
			c.Set("error", "Failed to create session")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.SetCookie("session_token", session.Token, 30*24*60*60, "/", "", false, true)

		c.Set("data", gin.H{
			"message": "Registration successful. Your account is pending approval.",
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"approved": user.Approved,
			},
		})
		c.Status(http.StatusCreated)
	}
}

func logout(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err == nil && token != "" {
			models.DeleteSession(token)
		}

		c.SetCookie("session_token", "", -1, "/", "", false, true)
		c.Set("data", gin.H{"message": "Logged out successfully"})
		c.Status(http.StatusOK)
	}
}

func getMe(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			c.Set("error", "You are not logged in. Please sign in to continue.")
			c.Status(http.StatusUnauthorized)
			return
		}

		session, err := models.GetSessionByToken(token)
		if err != nil {
			c.Set("error", "Your session has expired. Please sign in again.")
			c.Status(http.StatusUnauthorized)
			return
		}

		// Preload roles and permissions
		user, err := models.GetUserByID(session.User.ID)
		if err != nil {
			c.Set("error", "Unable to retrieve user information. Please try again.")
			c.Status(http.StatusNotFound)
			return
		}

		if !user.Approved {
			c.Set("error", "Your account is pending approval. Please wait for an administrator to approve your registration.")
			c.Status(http.StatusForbidden)
			return
		}

		c.Set("data", gin.H{
			"user": user,
		})
		c.Status(http.StatusOK)
	}
}
