package rest

import (
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

func RegisterPlayerRoutes(r *gin.Engine, api *Api) {
	players := r.Group("/players")
	players.Use(AuthMiddleware())
	players.Use(RequirePermission("players.view"))
	{
		players.GET("", getPlayers(api))
		players.GET("/:playerId", getPlayer(api))
	}
}

func getPlayers(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, err := api.rcon.Status()
		if err != nil {
			c.Set("error", err.Error())
			c.Status(500)
			return
		}

		c.Set("data", status.Players)
		c.Status(200)
	}
}

func getPlayer(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.Param("playerId")

		offlinePlayer := models.GetOfflinePlayer(identifier)
		if offlinePlayer != nil {
			c.Set("data", offlinePlayer)
			c.Status(200)
			return
		}

		player, err := api.rcon.GetPlayer(identifier)
		if err != nil || player.PlayerID == "" {
			c.Set("error", "Player not found")
			c.Status(404)
			return
		}

		ip := player.IP
		models.UpdateOfflinePlayer(&models.OfflinePlayer{
			PlayerID:      player.PlayerID,
			PlayerSteamID: player.PlayerSteamID,
			Name:          player.Name,
			IP:            ip,
			PBGuid:        player.PBGuid,
		}, nil)

		c.Set("data", player)
		c.Status(200)
	}
}
