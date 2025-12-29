package routes

import (
	"net/http"

	"github.com/zgiai/zgo/internal/app"
	"github.com/zgiai/zgo/internal/infra/router"
	"github.com/gin-gonic/gin"
)

// RegisterAPI registers all API routes using fluent router
func RegisterAPI(r *router.Router, handlers *app.Handlers) {
	// 1. Health Checks
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "v1"})
	}).Name("health")

	// 2. Register Module Routes
	handlers.User.RegisterRoutes(r)
	handlers.Permission.RegisterRoutes(r)
}
