package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

type CreateGroupRequest struct {
	Name        string   `json:"name" binding:"required"`
	Power       int      `json:"power" binding:"required,min=0,max=100"`
	Permissions []string `json:"permissions"`
	Description string   `json:"description"`
}

type UpdateGroupRequest struct {
	Name        string   `json:"name"`
	Power       int      `json:"power" binding:"min=0,max=100"`
	Permissions []string `json:"permissions"`
	Description string   `json:"description"`
}

type AssignPlayerRequest struct {
	PlayerID uint  `json:"playerId" binding:"required"`
	GroupID  *uint `json:"groupId"` // nil to remove from group
}

type CreateInGamePlayerRequest struct {
	GUID    string `json:"guid" binding:"required"`
	Name    string `json:"name" binding:"required"`
	GroupID *uint  `json:"groupId"`
}

type InGameIamGodRequest struct {
	GUID string `json:"guid" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func RegisterGroupRoutes(r *gin.Engine, api *Api) {
	// In-game iamgod endpoint - no auth required
	r.POST("/ingame/iamgod", inGameIamGod(api))

	groups := r.Group("/groups")
	groups.Use(AuthMiddleware())
	{
		// Group management
		groups.GET("", RequirePermission("groups.manage"), getAllGroups(api))
		groups.POST("", RequirePermission("groups.manage"), createGroup(api))
		groups.GET("/:id", RequirePermission("groups.manage"), getGroup(api))
		groups.PUT("/:id", RequirePermission("groups.manage"), updateGroup(api))
		groups.DELETE("/:id", RequirePermission("groups.manage"), deleteGroup(api))

		// In-game player management
		groups.GET("/players", RequirePermission("players.view"), getAllInGamePlayers(api))
		groups.POST("/players", RequirePermission("groups.manage"), createInGamePlayer(api))
		groups.PUT("/players/:id/assign", RequirePermission("groups.manage"), assignPlayerToGroup(api))
		groups.DELETE("/players/:id/group", RequirePermission("groups.manage"), removePlayerFromGroup(api))
	}
}

func getAllGroups(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		groups, err := models.GetAllGroups()
		if err != nil {
			c.Set("error", "Failed to retrieve groups")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", groups)
		c.Status(http.StatusOK)
	}
}

func getGroup(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid group ID")
			c.Status(http.StatusBadRequest)
			return
		}

		group, err := models.GetGroupByID(uint(id))
		if err != nil {
			c.Set("error", "Group not found")
			c.Status(http.StatusNotFound)
			return
		}

		c.Set("data", group)
		c.Status(http.StatusOK)
	}
}

func createGroup(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		// Convert permissions array to JSON string
		permissionsJSON, err := json.Marshal(req.Permissions)
		if err != nil {
			c.Set("error", "Invalid permissions format")
			c.Status(http.StatusBadRequest)
			return
		}

		err = models.CreateGroup(req.Name, req.Power, string(permissionsJSON), req.Description)
		if err != nil {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, err.Error(), "group", "", req.Name,
				map[string]interface{}{
					"power":       req.Power,
					"permissions": req.Permissions,
				},
				"Failed to create group")
			c.Set("error", "Failed to create group")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, "group_created", models.SourceWebUI,
			true, "", "group", "", req.Name,
			map[string]interface{}{
				"power":       req.Power,
				"permissions": req.Permissions,
				"description": req.Description,
			},
			"Group created successfully")

		c.Set("data", gin.H{"message": "Group created successfully"})
		c.Status(http.StatusCreated)
	}
}

func updateGroup(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid group ID")
			c.Status(http.StatusBadRequest)
			return
		}

		// Get existing group for audit trail
		group, err := models.GetGroupByID(uint(id))
		if err != nil {
			c.Set("error", "Group not found")
			c.Status(http.StatusNotFound)
			return
		}

		var req UpdateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		updates := make(map[string]interface{})
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.Power >= 0 {
			updates["power"] = req.Power
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
		if req.Description != "" {
			updates["description"] = req.Description
		}

		err = models.UpdateGroup(uint(id), updates)
		if err != nil {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, err.Error(), "group", strconv.FormatUint(uint64(id), 10), group.Name,
				updates,
				"Failed to update group")
			c.Set("error", "Failed to update group")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, "group_updated", models.SourceWebUI,
			true, "", "group", strconv.FormatUint(uint64(id), 10), group.Name,
			updates,
			"Group updated successfully")

		c.Set("data", gin.H{"message": "Group updated successfully"})
		c.Status(http.StatusOK)
	}
}

func deleteGroup(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid group ID")
			c.Status(http.StatusBadRequest)
			return
		}

		// Get group before deletion for audit trail
		group, err := models.GetGroupByID(uint(id))
		if err != nil {
			c.Set("error", "Group not found")
			c.Status(http.StatusNotFound)
			return
		}

		err = models.DeleteGroup(uint(id))
		if err != nil {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, err.Error(), "group", strconv.FormatUint(uint64(id), 10), group.Name,
				nil,
				"Failed to delete group")
			c.Set("error", "Failed to delete group")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, "group_deleted", models.SourceWebUI,
			true, "", "group", strconv.FormatUint(uint64(id), 10), group.Name,
			map[string]interface{}{"power": group.Power},
			"Group deleted successfully")

		c.Set("data", gin.H{"message": "Group deleted successfully"})
		c.Status(http.StatusOK)
	}
}

func getAllInGamePlayers(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Optional server ID filter
		var serverID *uint
		if serverIDStr := c.Query("server_id"); serverIDStr != "" {
			id, err := strconv.ParseUint(serverIDStr, 10, 32)
			if err != nil {
				c.Set("error", "Invalid server ID")
				c.Status(http.StatusBadRequest)
				return
			}
			sid := uint(id)
			serverID = &sid
		}

		players, err := models.GetAllInGamePlayers(serverID)
		if err != nil {
			c.Set("error", "Failed to retrieve players")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", players)
		c.Status(http.StatusOK)
	}
}

func createInGamePlayer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateInGamePlayerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		player, err := models.CreateOrUpdateInGamePlayer(req.GUID, req.Name)
		if err != nil {
			c.Set("error", "Failed to create player")
			c.Status(http.StatusInternalServerError)
			return
		}

		// Assign to group if provided
		if req.GroupID != nil {
			err = models.AssignPlayerToGroup(player.ID, *req.GroupID)
			if err != nil {
				c.Set("error", "Failed to assign player to group")
				c.Status(http.StatusInternalServerError)
				return
			}
		}

		c.Set("data", player)
		c.Status(http.StatusCreated)
	}
}

func assignPlayerToGroup(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid player ID")
			c.Status(http.StatusBadRequest)
			return
		}

		// Get player for audit trail
		player, err := models.GetInGamePlayerByID(uint(id))
		if err != nil {
			c.Set("error", "Player not found")
			c.Status(http.StatusNotFound)
			return
		}

		var req AssignPlayerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		var action models.ActionType
		var groupName string
		if req.GroupID == nil {
			err = models.RemovePlayerFromGroup(uint(id))
			action = "player_removed_from_group"
			groupName = "none"
		} else {
			err = models.AssignPlayerToGroup(uint(id), *req.GroupID)
			action = "player_assigned_to_group"
			group, _ := models.GetGroupByID(*req.GroupID)
			if group != nil {
				groupName = group.Name
			}
		}

		if err != nil {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, err.Error(), "player", player.GUID, player.Name,
				map[string]interface{}{"group_id": req.GroupID},
				"Failed to update player group")
			c.Set("error", "Failed to update player group")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, action, models.SourceWebUI,
			true, "", "player", player.GUID, player.Name,
			map[string]interface{}{"group": groupName},
			"Player group updated successfully")

		c.Set("data", gin.H{"message": "Player group updated successfully"})
		c.Status(http.StatusOK)
	}
}

func removePlayerFromGroup(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid player ID")
			c.Status(http.StatusBadRequest)
			return
		}

		err = models.RemovePlayerFromGroup(uint(id))
		if err != nil {
			c.Set("error", "Failed to remove player from group")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Player removed from group successfully"})
		c.Status(http.StatusOK)
	}
}

func inGameIamGod(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if in-game iamgod has already been used
		if models.HasSetting("ingame_iamgod_used", "true") {
			c.Set("error", "This command has already been used and cannot be re-enabled")
			c.Status(http.StatusForbidden)
			return
		}

		var req InGameIamGodRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		// Get or create the Owner group
		groups, err := models.GetAllGroups()
		if err != nil {
			c.Set("error", "Failed to retrieve groups")
			c.Status(http.StatusInternalServerError)
			return
		}

		var ownerGroup *models.Group
		for _, g := range groups {
			if g.Name == "Owner" {
				ownerGroup = &g
				break
			}
		}

		if ownerGroup == nil {
			c.Set("error", "Owner group not found")
			c.Status(http.StatusInternalServerError)
			return
		}

		// Create or update the in-game player
		player, err := models.GetInGamePlayerByGUID(req.GUID)
		if err != nil {
			// Player doesn't exist, create them
			player, err = models.CreateOrUpdateInGamePlayer(req.GUID, req.Name)
			if err != nil {
				c.Set("error", "Failed to create player")
				c.Status(http.StatusInternalServerError)
				return
			}
		}

		// Assign player to Owner group
		if err := models.AssignPlayerToGroup(player.ID, ownerGroup.ID); err != nil {
			c.Set("error", "Failed to assign player to Owner group")
			c.Status(http.StatusInternalServerError)
			return
		}

		// Mark as used in database - persist across restarts
		if err := models.SetSetting("ingame_iamgod_used", "true"); err != nil {
			c.Set("error", "Failed to mark iamgod as used")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{
			"message": "Owner privileges granted successfully. This command has been disabled.",
			"player": gin.H{
				"id":   player.ID,
				"guid": player.GUID,
				"name": player.Name,
			},
		})
		c.Status(http.StatusOK)
	}
}
