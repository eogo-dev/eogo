package commands

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/eogo-dev/eogo/internal/platform/config"
	"github.com/eogo-dev/eogo/internal/platform/console"
	"github.com/eogo-dev/eogo/internal/platform/container"
	"github.com/eogo-dev/eogo/internal/platform/database"
	"github.com/eogo-dev/eogo/internal/platform/email"
	"github.com/eogo-dev/eogo/internal/platform/jwt"
	"github.com/eogo-dev/eogo/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ServeCommand starts the HTTP server
type ServeCommand struct {
	output *console.Output
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{output: console.NewOutput()}
}

func (c *ServeCommand) Name() string        { return "serve" }
func (c *ServeCommand) Description() string { return "Start the HTTP server" }
func (c *ServeCommand) Usage() string       { return "serve [--port=8080]" }

func (c *ServeCommand) Run(args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Parse port override
	for i, arg := range args {
		if arg == "--port" && i+1 < len(args) {
			fmt.Sscanf(args[i+1], "%d", &cfg.Server.Port)
		}
	}

	container.App().Set(container.ServiceConfig, cfg)

	// Initialize services
	jwt.Init(cfg)
	container.App().Set(container.ServiceJWT, jwt.MustServiceInstance())

	email.Init(cfg)
	container.App().Set(container.ServiceEmail, email.MustServiceInstance())

	if cfg.Database.Enabled {
		db, err := database.InitDB(cfg.Database)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		container.App().Set(container.ServiceDB, db)
		c.output.Success("Database connected")
	}

	// Set Gin mode
	switch strings.ToLower(cfg.Server.Mode) {
	case "release", "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	r := gin.Default()

	// CORS
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
	}
	r.Use(cors.New(corsConfig))

	// Register routes
	routes.Setup(r)

	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	c.output.Success("Server starting on http://localhost%s", serverAddr)

	go func() {
		if err := r.Run(serverAddr); err != nil {
			c.output.Error("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	c.output.Info("Shutting down server...")
	return nil
}

// EnvCommand shows environment information
type EnvCommand struct {
	output *console.Output
}

func NewEnvCommand() *EnvCommand {
	return &EnvCommand{output: console.NewOutput()}
}

func (c *EnvCommand) Name() string        { return "env" }
func (c *EnvCommand) Description() string { return "Display the current environment" }
func (c *EnvCommand) Usage() string       { return "env" }

func (c *EnvCommand) Run(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	c.output.Title("Environment Information")

	c.output.TwoColumn("Environment", cfg.Server.Mode)
	c.output.TwoColumn("Server Port", fmt.Sprintf("%d", cfg.Server.Port))
	c.output.TwoColumn("Database Enabled", fmt.Sprintf("%v", cfg.Database.Enabled))
	if cfg.Database.Enabled {
		c.output.TwoColumn("Database Host", cfg.Database.Host)
		c.output.TwoColumn("Database Name", cfg.Database.DBName())
	}

	return nil
}

// VersionCommand shows version information
type VersionCommand struct {
	output  *console.Output
	version string
}

func NewVersionCommand(version string) *VersionCommand {
	return &VersionCommand{output: console.NewOutput(), version: version}
}

func (c *VersionCommand) Name() string        { return "version" }
func (c *VersionCommand) Description() string { return "Display application version" }
func (c *VersionCommand) Usage() string       { return "version" }

func (c *VersionCommand) Run(args []string) error {
	c.output.Info("Llama Gin Kit v%s", c.version)
	return nil
}

// RouteListCommand lists all registered routes
type RouteListCommand struct {
	output *console.Output
}

func NewRouteListCommand() *RouteListCommand {
	return &RouteListCommand{output: console.NewOutput()}
}

func (c *RouteListCommand) Name() string        { return "route:list" }
func (c *RouteListCommand) Description() string { return "List all registered routes" }
func (c *RouteListCommand) Usage() string       { return "route:list" }

func (c *RouteListCommand) Run(args []string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	routes.Setup(r)

	c.output.Title("Registered Routes")

	headers := []string{"Method", "Path", "Handler"}
	rows := make([][]string, 0)

	for _, route := range r.Routes() {
		rows = append(rows, []string{route.Method, route.Path, route.Handler})
	}

	c.output.Table(headers, rows)
	return nil
}
