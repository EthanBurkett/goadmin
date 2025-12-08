package rest

import (
	"net/http"
	"strconv"

	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
)

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AssignPermissionRequest struct {
	PermissionID uint `json:"permissionId" binding:"required"`
}

type AssignRoleRequest struct {
	RoleID uint `json:"roleId" binding:"required"`
}

type ApproveUserRequest struct {
	RoleID uint `json:"roleId" binding:"required"`
}

func RegisterRBACRoutes(r *gin.Engine, api *Api) {
	rbac := r.Group("/rbac")
	rbac.Use(AuthMiddleware())
	{
		// Roles
		roles := rbac.Group("/roles")
		roles.Use(RequirePermission("rbac.manage"))
		{
			roles.GET("", getAllRoles(api))
			roles.POST("", createRole(api))
			roles.GET("/:id", getRole(api))
			roles.DELETE("/:id", deleteRole(api))
			roles.POST("/:id/permissions", assignPermissionToRole(api))
			roles.DELETE("/:id/permissions/:permissionId", removePermissionFromRole(api))
		}

		// Permissions
		permissions := rbac.Group("/permissions")
		permissions.Use(RequirePermission("rbac.manage"))
		{
			permissions.GET("", getAllPermissions(api))
			permissions.POST("", createPermission(api))
			permissions.GET("/:id", getPermission(api))
			permissions.DELETE("/:id", deletePermission(api))
		}

		// User role management
		users := rbac.Group("/users")
		users.Use(RequirePermission("rbac.manage"))
		{
			users.GET("", getAllUsers(api))
			users.GET("/pending", getPendingUsers(api))
			users.POST("/:id/approve", approveUser(api))
			users.POST("/:id/deny", denyUser(api))
			users.DELETE("/:id", RequirePermission("users.delete"), deleteUser(api))
			users.GET("/:id", getUser(api))
			users.POST("/:id/roles", assignRoleToUser(api))
			users.DELETE("/:id/roles/:roleId", removeRoleFromUser(api))
		}
	}
}

// Role handlers
func getAllRoles(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := models.GetAllRoles()
		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Set("data", roles)
		c.Status(http.StatusOK)
	}
}

func createRole(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid role data. Please check your input and try again.")
			c.Status(http.StatusBadRequest)
			return
		}

		role, err := models.CreateRole(req.Name, req.Description)
		if err != nil {
			c.Set("error", "Failed to create role. The role name may already exist.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", role)
		c.Status(http.StatusCreated)
	}
}

func getRole(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		role, err := models.GetRoleByID(uint(id))
		if err != nil {
			c.Set("error", "Role not found")
			c.Status(http.StatusNotFound)
			return
		}

		c.Set("data", role)
		c.Status(http.StatusOK)
	}
}

func deleteRole(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.DeleteRole(uint(id)); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Role deleted successfully"})
		c.Status(http.StatusOK)
	}
}

func assignPermissionToRole(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		var req AssignPermissionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.AddPermissionToRole(uint(id), req.PermissionID); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Permission assigned to role"})
		c.Status(http.StatusOK)
	}
}

func removePermissionFromRole(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		permissionId, err := strconv.ParseUint(c.Param("permissionId"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid permission ID")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.RemovePermissionFromRole(uint(id), uint(permissionId)); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Permission removed from role"})
		c.Status(http.StatusOK)
	}
}

// Permission handlers
func getAllPermissions(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, err := models.GetAllPermissions()
		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Set("data", permissions)
		c.Status(http.StatusOK)
	}
}

func createPermission(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreatePermissionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid permission data. Please check your input and try again.")
			c.Status(http.StatusBadRequest)
			return
		}

		permission, err := models.CreatePermission(req.Name, req.Description)
		if err != nil {
			c.Set("error", "Failed to create permission. The permission name may already exist.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", permission)
		c.Status(http.StatusCreated)
	}
}

func getPermission(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		permission, err := models.GetPermissionByID(uint(id))
		if err != nil {
			c.Set("error", "Permission not found")
			c.Status(http.StatusNotFound)
			return
		}

		c.Set("data", permission)
		c.Status(http.StatusOK)
	}
}

func deletePermission(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.DeletePermission(uint(id)); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Permission deleted successfully"})
		c.Status(http.StatusOK)
	}
}

// User handlers
func getAllUsers(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := models.GetAllUsers()
		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Set("data", users)
		c.Status(http.StatusOK)
	}
}

func getUser(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		user, err := models.GetUserByID(uint(id))
		if err != nil {
			c.Set("error", "User not found")
			c.Status(http.StatusNotFound)
			return
		}

		c.Set("data", gin.H{
			"id":       user.ID,
			"username": user.Username,
			"roles":    user.Roles,
		})
		c.Status(http.StatusOK)
	}
}

func assignRoleToUser(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		var req AssignRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.AddRoleToUser(uint(id), req.RoleID); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Role assigned to user"})
		c.Status(http.StatusOK)
	}
}

func removeRoleFromUser(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		roleId, err := strconv.ParseUint(c.Param("roleId"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid role ID")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.RemoveRoleFromUser(uint(id), uint(roleId)); err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "Role removed from user"})
		c.Status(http.StatusOK)
	}
}

func getPendingUsers(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := models.GetPendingUsers()
		if err != nil {
			c.Set("error", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Set("data", users)
		c.Status(http.StatusOK)
	}
}

func approveUser(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		var req ApproveUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("error", "Invalid request. Please specify a role for the user.")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.ApproveUser(uint(id), req.RoleID); err != nil {
			c.Set("error", "Failed to approve user. The user or role may not exist.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "User has been approved and can now log in with their assigned role."})
		c.Status(http.StatusOK)
	}
}

func denyUser(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid ID")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.DenyUser(uint(id)); err != nil {
			c.Set("error", "Failed to deny user registration. The user may not exist.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "User registration has been denied and the account has been removed."})
		c.Status(http.StatusOK)
	}
}

func deleteUser(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Set("error", "Invalid user ID")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := models.DeleteUser(uint(id)); err != nil {
			c.Set("error", "Failed to delete user. The user may not exist or there was a database error.")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("data", gin.H{"message": "User account has been permanently deleted."})
		c.Status(http.StatusOK)
	}
}
