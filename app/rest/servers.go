package rest

import (
	"net/http"
	"strconv"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

type CreateServerRequest struct {
	Name         string `json:"name" binding:"required"`
	Host         string `json:"host" binding:"required"`
	Port         int    `json:"port" binding:"required"`
	RconPort     int    `json:"rconPort" binding:"required"`
	RconPassword string `json:"rconPassword" binding:"required"`
	GamesMpPath  string `json:"gamesMpPath"`
	Description  string `json:"description"`
	Region       string `json:"region"`
	MaxPlayers   int    `json:"maxPlayers"`
	IsDefault    bool   `json:"isDefault"`
}

type UpdateServerRequest struct {
	Name         *string `json:"name"`
	Host         *string `json:"host"`
	Port         *int    `json:"port"`
	RconPort     *int    `json:"rconPort"`
	RconPassword *string `json:"rconPassword"`
	GamesMpPath  *string `json:"gamesMpPath"`
	Description  *string `json:"description"`
	Region       *string `json:"region"`
	MaxPlayers   *int    `json:"maxPlayers"`
	IsActive     *bool   `json:"isActive"`
	IsDefault    *bool   `json:"isDefault"`
}

func RegisterServerRoutes(r *gin.Engine, api *Api) {
	servers := r.Group("/servers")
	servers.Use(AuthMiddleware())
	{
		servers.GET("", getAllServers(api))
		servers.GET("/active", getActiveServers(api))
		servers.GET("/default", getDefaultServer(api))
		servers.POST("", RequirePermission("servers.manage"), createServer(api))
		servers.GET("/:id", getServer(api))
		servers.PUT("/:id", RequirePermission("servers.manage"), updateServer(api))
		servers.DELETE("/:id", RequirePermission("servers.manage"), deleteServer(api))
		servers.POST("/:id/default", RequirePermission("servers.manage"), setDefaultServer(api))
		servers.POST("/:id/activate", RequirePermission("servers.manage"), activateServer(api))
		servers.POST("/:id/deactivate", RequirePermission("servers.manage"), deactivateServer(api))
	}
}

func getAllServers(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		servers, err := models.GetAllServers()
		if err != nil {
			c.Set("error", "Failed to retrieve servers")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", servers)
		c.Status(http.StatusOK)
	}
}

func getActiveServers(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		servers, err := models.GetActiveServers()
		if err != nil {
			c.Set("error", "Failed to retrieve active servers")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", servers)
		c.Status(http.StatusOK)
	}
}

func getDefaultServer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		server, err := models.GetDefaultServer()
		if err != nil {
			c.Set("error", "No default server configured")
			c.Status(http.StatusNotFound)
			return
		}

		c.Set("data", server)
		c.Status(http.StatusOK)
	}
}

func getServer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid server ID")
			c.Status(http.StatusBadRequest)
			return
		}

		server, err := models.GetServerByID(uint(id))
		if err != nil {
			c.Set("error", "Server not found")
			c.Status(http.StatusNotFound)
			return
		}

		c.Set("data", server)
		c.Status(http.StatusOK)
	}
}

func createServer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateServerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		server, err := models.CreateServer(
			req.Name,
			req.Host,
			req.RconPassword,
			req.GamesMpPath,
			req.Description,
			req.Region,
			req.Port,
			req.RconPort,
			req.MaxPlayers,
			req.IsDefault,
		)

		if err != nil {
			c.Set("error", "Failed to create server")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionCommandCreate, models.SourceWebUI,
			true, "", "server", "", server.Name,
			map[string]interface{}{
				"host":   server.Host,
				"port":   server.Port,
				"region": server.Region,
			},
			"Server created successfully")

		c.Set("data", server)
		c.Status(http.StatusCreated)
	}
}

func updateServer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid server ID")
			c.Status(http.StatusBadRequest)
			return
		}

		var req UpdateServerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}

		// Get existing server for audit trail
		server, err := models.GetServerByID(uint(id))
		if err != nil {
			c.Set("error", "Server not found")
			c.Status(http.StatusNotFound)
			return
		}

		updates := make(map[string]interface{})
		if req.Name != nil {
			updates["name"] = *req.Name
		}
		if req.Host != nil {
			updates["host"] = *req.Host
		}
		if req.Port != nil {
			updates["port"] = *req.Port
		}
		if req.RconPort != nil {
			updates["rcon_port"] = *req.RconPort
		}
		if req.RconPassword != nil {
			updates["rcon_password"] = *req.RconPassword
		}
		if req.GamesMpPath != nil {
			updates["games_mp_path"] = *req.GamesMpPath
		}
		if req.Description != nil {
			updates["description"] = *req.Description
		}
		if req.Region != nil {
			updates["region"] = *req.Region
		}
		if req.MaxPlayers != nil {
			updates["max_players"] = *req.MaxPlayers
		}
		if req.IsActive != nil {
			updates["is_active"] = *req.IsActive
		}
		if req.IsDefault != nil {
			updates["is_default"] = *req.IsDefault
		}

		err = models.UpdateServer(uint(id), updates)
		if err != nil {
			c.Set("error", "Failed to update server")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionCommandUpdate, models.SourceWebUI,
			true, "", "server", "", server.Name,
			map[string]interface{}{
				"updates": updates,
			},
			"Server updated successfully")

		c.Set("data", gin.H{"message": "Server updated successfully"})
		c.Status(http.StatusOK)
	}
}

func deleteServer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid server ID")
			c.Status(http.StatusBadRequest)
			return
		}

		// Get server for audit trail before deletion
		server, err := models.GetServerByID(uint(id))
		if err != nil {
			c.Set("error", "Server not found")
			c.Status(http.StatusNotFound)
			return
		}

		// Prevent deletion of default server
		if server.IsDefault {
			c.Set("error", "Cannot delete default server. Set another server as default first.")
			c.Status(http.StatusBadRequest)
			return
		}

		err = models.DeleteServer(uint(id))
		if err != nil {
			c.Set("error", "Failed to delete server")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionCommandDelete, models.SourceWebUI,
			true, "", "server", "", server.Name,
			nil,
			"Server deleted successfully")

		c.Set("data", gin.H{"message": "Server deleted successfully"})
		c.Status(http.StatusOK)
	}
}

func setDefaultServer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid server ID")
			c.Status(http.StatusBadRequest)
			return
		}

		server, err := models.GetServerByID(uint(id))
		if err != nil {
			c.Set("error", "Server not found")
			c.Status(http.StatusNotFound)
			return
		}

		err = server.SetAsDefault()
		if err != nil {
			c.Set("error", "Failed to set server as default")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionCommandUpdate, models.SourceWebUI,
			true, "", "server", "", server.Name,
			map[string]interface{}{
				"action": "set_default",
			},
			"Server set as default")

		c.Set("data", gin.H{"message": "Server set as default"})
		c.Status(http.StatusOK)
	}
}

func activateServer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid server ID")
			c.Status(http.StatusBadRequest)
			return
		}

		server, err := models.GetServerByID(uint(id))
		if err != nil {
			c.Set("error", "Server not found")
			c.Status(http.StatusNotFound)
			return
		}

		err = server.Activate()
		if err != nil {
			c.Set("error", "Failed to activate server")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionCommandUpdate, models.SourceWebUI,
			true, "", "server", "", server.Name,
			map[string]interface{}{
				"action": "activate",
			},
			"Server activated")

		c.Set("data", gin.H{"message": "Server activated"})
		c.Status(http.StatusOK)
	}
}

func deactivateServer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid server ID")
			c.Status(http.StatusBadRequest)
			return
		}

		server, err := models.GetServerByID(uint(id))
		if err != nil {
			c.Set("error", "Server not found")
			c.Status(http.StatusNotFound)
			return
		}

		// Prevent deactivation of default server
		if server.IsDefault {
			c.Set("error", "Cannot deactivate default server")
			c.Status(http.StatusBadRequest)
			return
		}

		err = server.Deactivate()
		if err != nil {
			c.Set("error", "Failed to deactivate server")
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionCommandUpdate, models.SourceWebUI,
			true, "", "server", "", server.Name,
			map[string]interface{}{
				"action": "deactivate",
			},
			"Server deactivated")

		c.Set("data", gin.H{"message": "Server deactivated"})
		c.Status(http.StatusOK)
	}
}
