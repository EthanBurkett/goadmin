package rest

import "github.com/gin-gonic/gin"

func RegisterStatusRoute(r *gin.Engine, api *Api) {
	status := r.Group("/status")
	status.Use(AuthMiddleware())
	status.Use(RequirePermission("status.view"))
	{
		status.GET("", func(c *gin.Context) {
			status, err := api.rcon.Status()
			if err != nil {
				c.Set("error", err.Error())
				c.Status(500)
				return
			}
			c.Set("data", status)
			c.Status(200)
		})
	}
}
