package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

type CreateCommandRequest struct {
	Name            string   `json:"name" binding:"required"`
	Usage           string   `json:"usage" binding:"required"`
	Description     string   `json:"description"`
	RconCommand     string   `json:"rconCommand" binding:"required"`
	MinArgs         int      `json:"minArgs"`
	MaxArgs         int      `json:"maxArgs"`
	MinPower        int      `json:"minPower" binding:"min=0,max=100"`
	Permissions     []string `json:"permissions"`
	RequirementType string   `json:"requirementType"` // "permission", "power", or "both"
}

type UpdateCommandRequest struct {
	Name            string   `json:"name"`
	Usage           string   `json:"usage"`
	Description     string   `json:"description"`
	RconCommand     string   `json:"rconCommand"`
	MinArgs         *int     `json:"minArgs"`
	MaxArgs         *int     `json:"maxArgs"`
	MinPower        *int     `json:"minPower"`
	Permissions     []string `json:"permissions"`
	RequirementType *string  `json:"requirementType"`
	Enabled         *bool    `json:"enabled"`
}

func RegisterCommandRoutes(r *gin.Engine, api *Api) {
	commands := r.Group("/commands")
	commands.Use(AuthMiddleware())
	{
		commands.GET("", RequirePermission("rbac.manage"), getAllCommands(api))
		commands.POST("", RequirePermission("rbac.manage"), createCommand(api))
		commands.GET("/:id", RequirePermission("rbac.manage"), getCommand(api))
		commands.PUT("/:id", RequirePermission("rbac.manage"), updateCommand(api))
		commands.DELETE("/:id", RequirePermission("rbac.manage"), deleteCommand(api))
	}
}

func getAllCommands(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		commands, err := models.GetAllCustomCommands()
		if err != nil {
			c.Set("error", "Failed to retrieve commands")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", commands)
		c.Status(http.StatusOK)
	}
}

func getCommand(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid command ID")
			c.Status(http.StatusBadRequest)
			return
		}

		commands, err := models.GetAllCustomCommands()
		if err != nil {
			c.Set("error", "Command not found")
			c.Status(http.StatusNotFound)
			return
		}

		for _, cmd := range commands {
			if cmd.ID == uint(id) {
				c.Set("data", cmd)
				c.Status(http.StatusOK)
				return
			}
		}

		c.Set("error", "Command not found")
		c.Status(http.StatusNotFound)
	}
}

func createCommand(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateCommandRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		// Convert permissions array to JSON string
		permissionsJSON := "[]"
		if req.Permissions != nil && len(req.Permissions) > 0 {
			data, err := json.Marshal(req.Permissions)
			if err != nil {
				c.Set("error", "Invalid permissions format")
				c.Status(http.StatusBadRequest)
				return
			}
			permissionsJSON = string(data)
		}

		// Default to "both" if not specified
		requirementType := req.RequirementType
		if requirementType == "" {
			requirementType = "both"
		}

		err := models.CreateCustomCommand(
			req.Name,
			req.Usage,
			req.Description,
			req.RconCommand,
			permissionsJSON,
			requirementType,
			req.MinArgs,
			req.MaxArgs,
			req.MinPower,
			false, // User-created commands are not built-in
		)
		if err != nil {
			c.Set("error", "Failed to create command")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Command created successfully"})
		c.Status(http.StatusCreated)
	}
}

func updateCommand(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid command ID")
			c.Status(http.StatusBadRequest)
			return
		}

		// Check if command is built-in
		cmd, err := models.GetCustomCommandByID(uint(id))
		if err != nil {
			c.Set("error", "Command not found")
			c.Status(http.StatusNotFound)
			return
		}
		if cmd.IsBuiltIn {
			c.Set("error", "Cannot modify built-in commands")
			c.Status(http.StatusForbidden)
			return
		}

		var req UpdateCommandRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		updates := make(map[string]interface{})
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.Usage != "" {
			updates["usage"] = req.Usage
		}
		if req.Description != "" {
			updates["description"] = req.Description
		}
		if req.RconCommand != "" {
			updates["rcon_command"] = req.RconCommand
		}
		if req.MinArgs != nil {
			updates["min_args"] = *req.MinArgs
		}
		if req.MaxArgs != nil {
			updates["max_args"] = *req.MaxArgs
		}
		if req.MinPower != nil {
			updates["min_power"] = *req.MinPower
		}
		if req.Permissions != nil {
			permissionsJSON, err := json.Marshal(req.Permissions)
			if err != nil {
				c.Set("error", "Invalid permissions format")
				c.Status(http.StatusBadRequest)
				return
			}
			updates["permissions"] = string(permissionsJSON)
		}
		if req.RequirementType != nil {
			updates["requirement_type"] = *req.RequirementType
		}
		if req.Enabled != nil {
			updates["enabled"] = *req.Enabled
		}

		err = models.UpdateCustomCommand(uint(id), updates)
		if err != nil {
			c.Set("error", "Failed to update command")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Command updated successfully"})
		c.Status(http.StatusOK)
	}
}

func deleteCommand(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid command ID")
			c.Status(http.StatusBadRequest)
			return
		}

		// Check if command is built-in
		cmd, err := models.GetCustomCommandByID(uint(id))
		if err != nil {
			c.Set("error", "Command not found")
			c.Status(http.StatusNotFound)
			return
		}
		if cmd.IsBuiltIn {
			c.Set("error", "Cannot delete built-in commands")
			c.Status(http.StatusForbidden)
			return
		}

		err = models.DeleteCustomCommand(uint(id))
		if err != nil {
			c.Set("error", "Failed to delete command")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Command deleted successfully"})
		c.Status(http.StatusOK)
	}
}
