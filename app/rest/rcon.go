package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

type RconCommandRequest struct {
	Command string `json:"command" binding:"required"`
}

type KickRequest struct {
	PlayerID string `json:"playerId" binding:"required"`
	Reason   string `json:"reason"`
}

type BanRequest struct {
	PlayerID string `json:"playerId" binding:"required"`
	Reason   string `json:"reason"`
}

type SayRequest struct {
	Message string `json:"message" binding:"required"`
}

type MapChangeRequest struct {
	MapName string `json:"mapName" binding:"required"`
}

type GametypeChangeRequest struct {
	Gametype string `json:"gametype" binding:"required"`
}

type TellRequest struct {
	PlayerID string `json:"playerId" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

type ExecRequest struct {
	Filename string `json:"filename" binding:"required"`
}

type UnbanRequest struct {
	PlayerName string `json:"playerName" binding:"required"`
}

type DumpUserRequest struct {
	PlayerName string `json:"playerName" binding:"required"`
}

type SetCvarRequest struct {
	Cvar  string `json:"cvar" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func RegisterRconRoutes(r *gin.Engine, api *Api) {
	rcon := r.Group("/rcon")
	rcon.Use(AuthMiddleware())
	rcon.Use(RateLimitByUser(RconRateLimiter)) // Add rate limiting for RCON commands
	{
		// Basic commands
		rcon.POST("/command", RequirePermission("rcon.command"), sendCommand(api))
		rcon.GET("/history", RequirePermission("rcon.command"), getCommandHistory(api))

		// Player management
		rcon.POST("/kick", RequirePermission("rcon.kick"), kickPlayer(api))
		rcon.POST("/ban", RequirePermission("rcon.ban"), banPlayer(api))
		rcon.POST("/unban", RequirePermission("rcon.ban"), unbanPlayer(api))
		rcon.POST("/dumpuser", RequirePermission("rcon.command"), dumpUser(api))
		rcon.POST("/tell", RequirePermission("rcon.say"), tellPlayer(api))

		// Server communication
		rcon.POST("/say", RequirePermission("rcon.say"), sayMessage(api))

		// Map/Game control
		rcon.POST("/map", RequirePermission("rcon.map"), changeMap(api))
		rcon.POST("/map-rotate", RequirePermission("rcon.map"), mapRotate(api))
		rcon.POST("/map-restart", RequirePermission("rcon.map"), mapRestart(api))
		rcon.POST("/fast-restart", RequirePermission("rcon.map"), fastRestart(api))
		rcon.POST("/gametype", RequirePermission("rcon.map"), changeGametype(api))

		// Server control
		rcon.POST("/exec", RequirePermission("rcon.admin"), execConfig(api))
		rcon.POST("/writeconfig", RequirePermission("rcon.admin"), writeConfig(api))
		rcon.POST("/set", RequirePermission("rcon.admin"), setCvar(api))
		rcon.GET("/serverinfo", RequirePermission("status.view"), getServerInfo(api))
		rcon.GET("/systeminfo", RequirePermission("status.view"), getSystemInfo(api))

		// Stats endpoints
		rcon.GET("/stats/server", RequirePermission("status.view"), getServerStats(api))
		rcon.GET("/stats/system", RequirePermission("status.view"), getSystemStats(api))
		rcon.GET("/stats/players", RequirePermission("status.view"), getPlayerStats(api))
	}
}

func sendCommand(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RconCommandRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		// Validate and sanitize command
		sanitizedCommand, err := ValidateRconCommand(req.Command)
		if err != nil {
			// Log security violation for invalid commands
			violationType := "invalid_command"
			if CommandValidatorInstance.IsRestrictedCommand(req.Command) {
				violationType = "restricted_command_attempt"
			}
			Audit.LogSecurityViolation(c, violationType, req.Command, err.Error())

			c.Set("error", fmt.Sprintf("Invalid command: %v", err))
			c.Status(http.StatusBadRequest)
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		response, err := api.rcon.SendCommand(sanitizedCommand)
		success := err == nil

		// Save command history
		if success {
			models.CreateCommandHistory(user.ID, sanitizedCommand, response, true)
		} else {
			models.CreateCommandHistory(user.ID, sanitizedCommand, err.Error(), false)
		}

		// Log to audit trail
		var errorMsg string
		if err != nil {
			errorMsg = err.Error()
		}
		Audit.LogRconCommand(c, sanitizedCommand, response, success, errorMsg)

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func kickPlayer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req KickRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "clientkick " + req.PlayerID
		if req.Reason != "" {
			command += " " + req.Reason
		}

		response, err := api.rcon.SendCommand(command)
		success := err == nil

		// Save command history
		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		// Log to audit trail
		var errorMsg string
		if err != nil {
			errorMsg = err.Error()
		}
		Audit.LogKick(c, "", req.PlayerID, req.Reason, success, errorMsg)

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func banPlayer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "banclient " + req.PlayerID
		if req.Reason != "" {
			command += " " + req.Reason
		}

		response, err := api.rcon.SendCommand(command)
		success := err == nil

		// Save command history
		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		// Log to audit trail
		var errorMsg string
		if err != nil {
			errorMsg = err.Error()
		}
		Audit.LogBan(c, "", req.PlayerID, req.Reason, success, errorMsg)

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func sayMessage(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SayRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "say " + req.Message
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		// Save command history
		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func getCommandHistory(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		// Get command history for this user (limit to last 50)
		history, err := models.GetCommandHistoryByUser(user.ID, 50)
		if err != nil {
			c.Set("error", "Failed to retrieve command history")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", history)
		c.Status(http.StatusOK)
	}
}

// Player Management Handlers

func unbanPlayer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UnbanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "unbanUser " + req.PlayerName
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func dumpUser(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DumpUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "dumpuser " + req.PlayerName
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func tellPlayer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TellRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "tell " + req.PlayerID + " " + req.Message
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

// Map/Game Control Handlers

func changeMap(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MapChangeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "map " + req.MapName
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func mapRotate(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "map_rotate"
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func mapRestart(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "map_restart"
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func fastRestart(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "fast_restart"
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func changeGametype(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GametypeChangeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "g_gametype " + req.Gametype
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

// Server Control Handlers

func execConfig(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ExecRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "exec " + req.Filename
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func writeConfig(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ExecRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "writeconfig " + req.Filename
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func setCvar(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SetCvarRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.Set("error", "User not found")
			c.Status(http.StatusUnauthorized)
			return
		}
		user := userVal.(*models.User)

		command := "set " + req.Cvar + " " + req.Value
		response, err := api.rcon.SendCommand(command)
		success := err == nil

		if success {
			models.CreateCommandHistory(user.ID, command, response, true)
		} else {
			models.CreateCommandHistory(user.ID, command, err.Error(), false)
		}

		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func getServerInfo(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := api.rcon.SendCommand("serverinfo")
		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

func getSystemInfo(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := api.rcon.SendCommand("systeminfo")
		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"response": response})
		c.Status(http.StatusOK)
	}
}

// Stats Handlers

func getServerStats(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Default to last 24 hours
		endTime := time.Now()
		startTime := endTime.Add(-24 * time.Hour)

		// Allow custom time range via query params
		if start := c.Query("start"); start != "" {
			if t, err := time.Parse(time.RFC3339, start); err == nil {
				startTime = t
			}
		}
		if end := c.Query("end"); end != "" {
			if t, err := time.Parse(time.RFC3339, end); err == nil {
				endTime = t
			}
		}

		stats, err := models.GetServerStatsRange(startTime, endTime)
		if err != nil {
			c.Set("error", "Failed to retrieve server stats")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", stats)
		c.Status(http.StatusOK)
	}
}

func getSystemStats(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		endTime := time.Now()
		startTime := endTime.Add(-24 * time.Hour)

		if start := c.Query("start"); start != "" {
			if t, err := time.Parse(time.RFC3339, start); err == nil {
				startTime = t
			}
		}
		if end := c.Query("end"); end != "" {
			if t, err := time.Parse(time.RFC3339, end); err == nil {
				endTime = t
			}
		}

		stats, err := models.GetSystemStatsRange(startTime, endTime)
		if err != nil {
			c.Set("error", "Failed to retrieve system stats")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", stats)
		c.Status(http.StatusOK)
	}
}

func getPlayerStats(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		endTime := time.Now()
		startTime := endTime.Add(-24 * time.Hour)

		if start := c.Query("start"); start != "" {
			if t, err := time.Parse(time.RFC3339, start); err == nil {
				startTime = t
			}
		}
		if end := c.Query("end"); end != "" {
			if t, err := time.Parse(time.RFC3339, end); err == nil {
				endTime = t
			}
		}

		stats, err := models.GetPlayerStatsRange(startTime, endTime)
		if err != nil {
			c.Set("error", "Failed to retrieve player stats")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", stats)
		c.Status(http.StatusOK)
	}
}
