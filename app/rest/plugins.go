package rest

import (
	"net/http"

	"github.com/ethanburkett/goadmin/app/plugins"
	"github.com/gin-gonic/gin"
)

// getAllPlugins returns all loaded plugins
func getAllPlugins(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		if plugins.GlobalPluginManager == nil {
			c.Set("data", gin.H{"plugins": []interface{}{}})
			c.Status(http.StatusOK)
			return
		}

		statuses := plugins.GlobalPluginManager.GetStatus()
		c.Set("data", gin.H{"plugins": statuses})
		c.Status(http.StatusOK)
	}
}

// getPluginStatus returns the status of a specific plugin
func getPluginStatus(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		pluginID := c.Param("id")

		if plugins.GlobalPluginManager == nil {
			c.Set("error", "Plugin manager not initialized")
			c.Status(http.StatusNotFound)
			return
		}

		status, err := plugins.GlobalPluginManager.GetPluginStatus(pluginID)
		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusNotFound)
			return
		}

		c.Set("data", status)
		c.Status(http.StatusOK)
	}
}

// startPlugin starts a plugin
func startPlugin(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		pluginID := c.Param("id")

		if plugins.GlobalPluginManager == nil {
			c.Set("error", "Plugin manager not initialized")
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := plugins.GlobalPluginManager.Start(pluginID); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		helper := &AuditHelper{}
		helper.LogAction(c, "plugin.started", "web_ui", true, "", "plugin", pluginID, pluginID, map[string]interface{}{
			"plugin_id": pluginID,
		}, "")

		c.Set("data", gin.H{"message": "Plugin started successfully"})
		c.Status(http.StatusOK)
	}
}

// stopPlugin stops a plugin
func stopPlugin(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		pluginID := c.Param("id")

		if plugins.GlobalPluginManager == nil {
			c.Set("error", "Plugin manager not initialized")
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := plugins.GlobalPluginManager.Stop(pluginID); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		helper := &AuditHelper{}
		helper.LogAction(c, "plugin.stopped", "web_ui", true, "", "plugin", pluginID, pluginID, map[string]interface{}{
			"plugin_id": pluginID,
		}, "")

		c.Set("data", gin.H{"message": "Plugin stopped successfully"})
		c.Status(http.StatusOK)
	}
}

// reloadPlugin reloads a plugin
func reloadPlugin(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		pluginID := c.Param("id")

		if plugins.GlobalPluginManager == nil {
			c.Set("error", "Plugin manager not initialized")
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := plugins.GlobalPluginManager.Reload(pluginID); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		helper := &AuditHelper{}
		helper.LogAction(c, "plugin.reloaded", "web_ui", true, "", "plugin", pluginID, pluginID, map[string]interface{}{
			"plugin_id": pluginID,
		}, "")

		c.Set("data", gin.H{"message": "Plugin reloaded successfully"})
		c.Status(http.StatusOK)
	}
}

// RegisterPluginRoutes registers plugin management routes
func RegisterPluginRoutes(r *gin.Engine, api *Api) {
	plugins := r.Group("/plugins")
	plugins.Use(AuthMiddleware())
	{
		// List all plugins (requires plugins.view)
		plugins.GET("", RequirePermission("plugins.view"), getAllPlugins(api))

		// Get specific plugin status (requires plugins.view)
		plugins.GET("/:id", RequirePermission("plugins.view"), getPluginStatus(api))

		// Start plugin (requires plugins.manage)
		plugins.POST("/:id/start", RequirePermission("plugins.manage"), startPlugin(api))

		// Stop plugin (requires plugins.manage)
		plugins.POST("/:id/stop", RequirePermission("plugins.manage"), stopPlugin(api))

		// Reload plugin (requires plugins.manage)
		plugins.POST("/:id/reload", RequirePermission("plugins.manage"), reloadPlugin(api))
	}
}
