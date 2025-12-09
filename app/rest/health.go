package rest

import (
	"net/http"
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"github.com/ethanburkett/goadmin/app/rcon"
	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]Health `json:"checks"`
}

// Health represents an individual health check
type Health struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func RegisterHealthRoutes(router *gin.Engine, api *Api) {
	health := router.Group("/health")
	{
		health.GET("", getHealth(api))
		health.GET("/ready", getReadiness(api))
		health.GET("/live", getLiveness())
	}
}

// getHealth returns detailed health status of all components
func getHealth(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := HealthResponse{
			Timestamp: time.Now().Format(time.RFC3339),
			Checks:    make(map[string]Health),
		}

		// Check database
		dbHealth := checkDatabase()
		response.Checks["database"] = dbHealth

		// Check RCON
		rconHealth := checkRCON(api.rcon)
		response.Checks["rcon"] = rconHealth

		// Determine overall status
		allHealthy := dbHealth.Status == "healthy" && rconHealth.Status == "healthy"
		if allHealthy {
			response.Status = "healthy"
		} else if dbHealth.Status == "unhealthy" {
			response.Status = "unhealthy"
		} else {
			response.Status = "degraded"
		}

		statusCode := http.StatusOK
		if response.Status == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, response)
	}
}

// getReadiness returns whether the service is ready to accept traffic
func getReadiness(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if database is accessible
		dbHealth := checkDatabase()
		if dbHealth.Status != "healthy" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "not ready",
				"reason":  "database unavailable",
				"message": dbHealth.Message,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	}
}

// getLiveness returns whether the service is alive (for k8s liveness probe)
func getLiveness() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	}
}

// checkDatabase checks database connectivity and connection pool status
func checkDatabase() Health {
	if database.DB == nil {
		return Health{
			Status:  "unhealthy",
			Message: "database not initialized",
		}
	}

	sqlDB, err := database.DB.DB()
	if err != nil {
		return Health{
			Status:  "unhealthy",
			Message: "failed to get underlying database",
		}
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		return Health{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	// Get connection pool stats
	stats := sqlDB.Stats()
	details := map[string]interface{}{
		"open_connections":   stats.OpenConnections,
		"in_use":             stats.InUse,
		"idle":               stats.Idle,
		"wait_count":         stats.WaitCount,
		"wait_duration":      stats.WaitDuration.String(),
		"max_idle_closed":    stats.MaxIdleClosed,
		"max_lifetime_close": stats.MaxLifetimeClosed,
	}

	return Health{
		Status:  "healthy",
		Details: details,
	}
}

// checkRCON checks RCON server connectivity
func checkRCON(rconClient *rcon.Client) Health {
	if rconClient == nil {
		return Health{
			Status:  "unhealthy",
			Message: "RCON client not initialized",
		}
	}

	// Try a simple command
	_, err := rconClient.SendCommand("status")
	if err != nil {
		return Health{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	return Health{
		Status: "healthy",
	}
}
