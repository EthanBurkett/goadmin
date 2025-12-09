package rest

import (
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
		commands.GET("", RequirePermission("commands.manage"), getAllCommands(api))
		commands.POST("", RequirePermission("commands.manage"), createCommand(api))
		commands.GET("/:id", RequirePermission("commands.manage"), getCommand(api))
		commands.PUT("/:id", RequirePermission("commands.manage"), updateCommand(api))
		commands.DELETE("/:id", RequirePermission("commands.manage"), deleteCommand(api))
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

		// Convert permission names to IDs
		var permissionIDs []uint
		for _, permName := range req.Permissions {
			perm, err := models.GetPermissionByName(permName)
			if err != nil {
				c.Set("error", "Invalid permission: "+permName)
				c.Status(http.StatusBadRequest)
				return
			}
			permissionIDs = append(permissionIDs, perm.ID)
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
			requirementType,
			req.MinArgs,
			req.MaxArgs,
			req.MinPower,
			false, // User-created commands are not built-in
			permissionIDs,
		)
		if err != nil {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, err.Error(), "command", "", req.Name,
				map[string]interface{}{
					"rcon":      req.RconCommand,
					"min_power": req.MinPower,
				},
				"Failed to create command")
			c.Set("error", "Failed to create command")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionCommandCreate, models.SourceWebUI,
			true, "", "command", "", req.Name,
			map[string]interface{}{
				"rcon":        req.RconCommand,
				"min_power":   req.MinPower,
				"permissions": req.Permissions,
			},
			"Command created successfully")

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

		// Check if command is built-in and pre-fetch for audit trail
		cmd, err := models.GetCustomCommandByID(uint(id))
		if err != nil {
			c.Set("error", "Command not found")
			c.Status(http.StatusNotFound)
			return
		}
		if cmd.IsBuiltIn {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, "Attempted to modify built-in command", "command", "", cmd.Name,
				map[string]interface{}{
					"reason": "cannot_modify_builtin",
				},
				"Cannot modify built-in commands")
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
		if req.RequirementType != nil {
			updates["requirement_type"] = *req.RequirementType
		}
		if req.Enabled != nil {
			updates["enabled"] = *req.Enabled
		}

		// Handle permissions separately using the association API
		if req.Permissions != nil {
			var permissionIDs []uint
			for _, permName := range req.Permissions {
				perm, err := models.GetPermissionByName(permName)
				if err != nil {
					c.Set("error", "Invalid permission: "+permName)
					c.Status(http.StatusBadRequest)
					return
				}
				permissionIDs = append(permissionIDs, perm.ID)
			}
			if err := cmd.SetCommandPermissions(permissionIDs); err != nil {
				c.Set("error", "Failed to update permissions")
				c.Status(http.StatusInternalServerError)
				return
			}
		}

		err = models.UpdateCustomCommand(uint(id), updates)
		if err != nil {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, err.Error(), "command", "", cmd.Name,
				map[string]interface{}{
					"updates": updates,
				},
				"Failed to update command")
			c.Set("error", "Failed to update command")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionCommandUpdate, models.SourceWebUI,
			true, "", "command", "", cmd.Name,
			map[string]interface{}{
				"updates": updates,
			},
			"Command updated successfully")

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

		// Check if command is built-in and pre-fetch for audit trail
		cmd, err := models.GetCustomCommandByID(uint(id))
		if err != nil {
			c.Set("error", "Command not found")
			c.Status(http.StatusNotFound)
			return
		}
		if cmd.IsBuiltIn {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, "Attempted to delete built-in command", "command", "", cmd.Name,
				map[string]interface{}{
					"reason": "cannot_delete_builtin",
				},
				"Cannot delete built-in commands")
			c.Set("error", "Cannot delete built-in commands")
			c.Status(http.StatusForbidden)
			return
		}

		err = models.DeleteCustomCommand(uint(id))
		if err != nil {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, err.Error(), "command", "", cmd.Name,
				map[string]interface{}{},
				"Failed to delete command")
			c.Set("error", "Failed to delete command")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionCommandDelete, models.SourceWebUI,
			true, "", "command", "", cmd.Name,
			map[string]interface{}{
				"rcon":        cmd.RconCommand,
				"min_power":   cmd.MinPower,
				"permissions": cmd.Permissions,
			},
			"Command deleted successfully")

		c.Set("data", gin.H{"message": "Command deleted successfully"})
		c.Status(http.StatusOK)
	}
}
