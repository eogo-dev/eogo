package bootstrap

import (
	"log"

	"github.com/eogo-dev/eogo/internal/platform/config"
	"github.com/eogo-dev/eogo/internal/platform/container"
	"github.com/eogo-dev/eogo/internal/platform/database"
	"github.com/eogo-dev/eogo/internal/platform/email"
	"github.com/eogo-dev/eogo/internal/platform/jwt"
	"github.com/eogo-dev/eogo/internal/platform/logger"
	"gorm.io/gorm"
)

// Application represents the bootstrapped application
type Application struct {
	Config *config.Config
	DB     *gorm.DB
}

// NewApp bootstraps the application by loading config and services
func NewApp() *Application {
	// 1. Initialize Logger
	logger.Boot()

	// 2. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	container.App().Set(container.ServiceConfig, cfg)

	// 2. Initialize Core Services
	// JWT
	jwt.Init(cfg)
	container.App().Set(container.ServiceJWT, jwt.MustServiceInstance())

	// Email
	email.Init(cfg)
	container.App().Set(container.ServiceEmail, email.MustServiceInstance())

	// 3. Initialize Database
	var db *gorm.DB
	if cfg.Database.Enabled {
		var err error
		db, err = database.InitDB(cfg.Database)
		if err != nil {
			log.Printf("Warning: Failed to initialize database: %v", err)
		} else {
			container.App().Set(container.ServiceDB, db)
		}
	} else {
		log.Println("Database initialization skipped (DB_ENABLED=false)")
	}

	return &Application{
		Config: cfg,
		DB:     db,
	}
}
