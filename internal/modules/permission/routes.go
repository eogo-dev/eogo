package permission

import (
	"github.com/eogo-dev/eogo/internal/platform/database"
	"github.com/eogo-dev/eogo/internal/platform/router"
)

// Register registers permission module routes
func Register(r *router.Router) {
	db := database.GetDB()
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	// Role routes (admin only)
	r.Group("", func(auth *router.Router) {
		auth.WithMiddleware("auth")

		// Role management
		auth.POST("/roles", handler.CreateRole).Name("roles.store")
		auth.GET("/roles", handler.ListRoles).Name("roles.index")
		auth.GET("/roles/:id", handler.GetRole).Name("roles.show").WhereNumber("id")
		auth.PUT("/roles/:id", handler.UpdateRole).Name("roles.update").WhereNumber("id")
		auth.DELETE("/roles/:id", handler.DeleteRole).Name("roles.destroy").WhereNumber("id")

		// Role assignment
		auth.POST("/roles/assign", handler.AssignRole).Name("roles.assign")
		auth.POST("/roles/remove", handler.RemoveRole).Name("roles.remove")

		// User roles
		auth.GET("/users/:id/roles", handler.GetUserRoles).Name("users.roles").WhereNumber("id")

		// Permissions
		auth.GET("/permissions", handler.ListPermissions).Name("permissions.index")
	})
}
