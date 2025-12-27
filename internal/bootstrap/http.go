package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/eogo-dev/eogo/internal/platform/logger"
	"github.com/eogo-dev/eogo/pkg/support"
	"github.com/eogo-dev/eogo/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// HttpKernel handles HTTP server lifecycle
type HttpKernel struct {
	App    *Application
	Engine *gin.Engine
}

// NewHttpKernel creates a new HTTP kernel
func NewHttpKernel(app *Application) *HttpKernel {
	// Set Mode
	setGinMode(app.Config.Server.Mode)

	// Create Engine
	r := gin.New()

	// Add custom logger and recovery middleware
	r.Use(logger.GinLogger())
	r.Use(gin.Recovery())

	// Apply Global Middleware (CORS mainly)
	applyGlobalMiddleware(r, app)

	// Register Routes
	// We temporarily silence Gin's default route logging to keep console clean
	gin.SetMode(gin.ReleaseMode) // Temporarily set to release to silence route logs
	routes.Setup(r)
	setGinMode(app.Config.Server.Mode) // Restore correct mode

	// Print Professional Banner
	support.PrintBanner("1.0.0")

	return &HttpKernel{
		App:    app,
		Engine: r,
	}
}

// Handle starts the HTTP server
func (k *HttpKernel) Handle() {
	cfg := k.App.Config
	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: k.Engine,
	}

	// Start Server
	go func() {
		// Build clickable URL
		host := cfg.Server.Host
		if host == "" {
			host = "localhost"
		}
		url := fmt.Sprintf("http://%s:%d", host, cfg.Server.Port)

		log.Printf("\n")
		log.Printf("  🚀 Eogo Server Started!")
		log.Printf("  ➜ Local:   \033[36m%s\033[0m", url)
		log.Printf("  ➜ Mode:    %s", cfg.Server.Mode)
		log.Printf("\n")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

func setGinMode(mode string) {
	switch strings.ToLower(mode) {
	case "release", "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func applyGlobalMiddleware(r *gin.Engine, app *Application) {
	cfg := app.Config
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
	}
	r.Use(cors.New(corsConfig))
}
