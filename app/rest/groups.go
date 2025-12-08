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
		groups.GET("", RequirePermission("rbac.manage"), getAllGroups(api))
		groups.POST("", RequirePermission("rbac.manage"), createGroup(api))
		groups.GET("/:id", RequirePermission("rbac.manage"), getGroup(api))
		groups.PUT("/:id", RequirePermission("rbac.manage"), updateGroup(api))
		groups.DELETE("/:id", RequirePermission("rbac.manage"), deleteGroup(api))

		// In-game player management
		groups.GET("/players", RequirePermission("players.view"), getAllInGamePlayers(api))
		groups.POST("/players", RequirePermission("rbac.manage"), createInGamePlayer(api))
		groups.PUT("/players/:id/assign", RequirePermission("rbac.manage"), assignPlayerToGroup(api))
		groups.DELETE("/players/:id/group", RequirePermission("rbac.manage"), removePlayerFromGroup(api))
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
			c.Set("error", "Failed to create group")
			c.Status(http.StatusInternalServerError)
			return
		}

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
			c.Set("error", "Failed to update group")
			c.Status(http.StatusInternalServerError)
			return
		}

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

		err = models.DeleteGroup(uint(id))
		if err != nil {
			c.Set("error", "Failed to delete group")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Group deleted successfully"})
		c.Status(http.StatusOK)
	}
}

func getAllInGamePlayers(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		players, err := models.GetAllInGamePlayers()
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

		var req AssignPlayerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		if req.GroupID == nil {
			err = models.RemovePlayerFromGroup(uint(id))
		} else {
			err = models.AssignPlayerToGroup(uint(id), *req.GroupID)
		}

		if err != nil {
			c.Set("error", "Failed to update player group")
			c.Status(http.StatusInternalServerError)
			return
		}

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
