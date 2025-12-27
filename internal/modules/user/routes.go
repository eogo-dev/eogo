package user

import (
	"github.com/eogo-dev/eogo/internal/platform/database"
	"github.com/eogo-dev/eogo/internal/platform/jwt"
	"github.com/eogo-dev/eogo/internal/platform/router"
)

// Register registers user module routes
func Register(r *router.Router) {
	db := database.GetDB()
	repo := NewRepository(db)
	jwtSvc := jwt.MustServiceInstance()
	service := NewService(repo, jwtSvc)
	handler := NewHandler(service)

	// Public routes
	r.POST("/register", handler.Register).Name("auth.register")
	r.POST("/login", handler.Login).Name("auth.login")
	r.POST("/password/reset", handler.ResetPassword).Name("auth.password.reset")

	// Protected routes
	r.Group("", func(auth *router.Router) {
		auth.WithMiddleware("auth")

		// Profile
		auth.GET("/users/profile", handler.GetProfile).Name("users.profile")
		auth.PUT("/users/profile", handler.UpdateProfile).Name("users.profile.update")
		auth.PUT("/users/password", handler.ChangePassword).Name("users.password.update")
		auth.DELETE("/users/account", handler.DeleteAccount).Name("users.account.delete")

		// User management
		auth.GET("/users", handler.List).Name("users.index")
		auth.GET("/users/:id", handler.Get).Name("users.show").WhereNumber("id")
		auth.GET("/users/:id/info", handler.GetUserInfo).Name("users.info").WhereNumber("id")
	})
}
