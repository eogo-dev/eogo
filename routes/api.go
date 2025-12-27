package routes

import (
	"net/http"

	"github.com/eogo-dev/eogo/internal/modules/permission"
	"github.com/eogo-dev/eogo/internal/modules/user"
	pkgmiddleware "github.com/eogo-dev/eogo/internal/platform/middleware"
	"github.com/eogo-dev/eogo/internal/platform/router"
	"github.com/gin-gonic/gin"
)

// RegisterAPI registers all API routes using fluent router
func RegisterAPI(r *router.Router) {
	// 1. Health Checks
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "v1"})
	}).Name("health")

	// 2. Register Module Routes
	user.Register(r)
	permission.Register(r)
}

// SetupMiddleware configures middleware groups for the router
func SetupMiddleware(r *router.Router) {
	r.MiddlewareGroup("auth", pkgmiddleware.JWTAuth())
	r.AliasMiddleware("jwt", pkgmiddleware.JWTAuth())
}
