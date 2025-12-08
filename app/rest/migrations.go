package rest

import (
	"net/http"

	"github.com/ethanburkett/goadmin/app/database"
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

func RegisterMigrationRoutes(r *gin.Engine, api *Api) {
	migrations := r.Group("/migrations")
	migrations.Use(AuthMiddleware())
	migrations.Use(RequirePermission("rbac.manage")) // Only super admins
	{
		migrations.GET("", getMigrations(api))
		migrations.GET("/status", getMigrationStatus(api))
		migrations.GET("/current", getCurrentMigration(api))
		migrations.POST("/apply", applyMigrations(api))
		migrations.POST("/rollback", rollbackLastMigration(api))
	}
}

// getMigrations returns all migrations
func getMigrations(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		migrations, err := models.GetAllMigrations(api.DB)
		if err != nil {
			c.Set("error", "Failed to fetch migrations")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"migrations": migrations})
		c.Status(http.StatusOK)
	}
}

// getMigrationStatus returns migration system status
func getMigrationStatus(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		runner := database.NewMigrationRunner(api.DB)
		for _, migration := range api.migrations {
			runner.Register(migration)
		}

		status, err := runner.GetStatus()
		if err != nil {
			c.Set("error", "Failed to get migration status")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", status)
		c.Status(http.StatusOK)
	}
}

// getCurrentMigration returns the current migration version
func getCurrentMigration(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		version, err := models.GetCurrentMigrationVersion(api.DB)
		if err != nil {
			c.Set("error", "Failed to get current migration")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{
			"version": version,
		})
		c.Status(http.StatusOK)
	}
}

// applyMigrations applies all pending migrations
func applyMigrations(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		runner := database.NewMigrationRunner(api.DB)
		for _, migration := range api.migrations {
			runner.Register(migration)
		}

		if err := runner.ApplyAll(); err != nil {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, err.Error(), "migration", "", "",
				map[string]interface{}{"error": err.Error()},
				"Failed to apply migrations")
			c.Set("error", "Failed to apply migrations: "+err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionSystemChange, models.SourceWebUI,
			true, "", "migration", "", "",
			map[string]interface{}{},
			"Applied pending migrations")

		c.Set("data", gin.H{"message": "All migrations applied successfully"})
		c.Status(http.StatusOK)
	}
}

// rollbackLastMigration rolls back the most recent migration
func rollbackLastMigration(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		runner := database.NewMigrationRunner(api.DB)
		for _, migration := range api.migrations {
			runner.Register(migration)
		}

		if err := runner.RollbackLast(); err != nil {
			Audit.LogAction(c, models.ActionSecurityViolation, models.SourceWebUI,
				false, err.Error(), "migration", "", "",
				map[string]interface{}{"error": err.Error()},
				"Failed to rollback migration")
			c.Set("error", "Failed to rollback migration: "+err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		Audit.LogAction(c, models.ActionSystemChange, models.SourceWebUI,
			true, "", "migration", "", "",
			map[string]interface{}{},
			"Rolled back last migration")

		c.Set("data", gin.H{"message": "Migration rolled back successfully"})
		c.Status(http.StatusOK)
	}
}
