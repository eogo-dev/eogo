package feature

import (
	"fmt"
	"testing"
	"time"

	"github.com/eogo-dev/eogo/internal/app"
	"github.com/eogo-dev/eogo/internal/bootstrap"
	"github.com/eogo-dev/eogo/internal/infra/config"
	"github.com/eogo-dev/eogo/internal/infra/database"
	"github.com/eogo-dev/eogo/internal/infra/email"
	"github.com/eogo-dev/eogo/internal/infra/jwt"
	"github.com/eogo-dev/eogo/internal/infra/middleware"
	test_platform "github.com/eogo-dev/eogo/internal/infra/testing"
	"github.com/eogo-dev/eogo/internal/modules/permission"
	"github.com/eogo-dev/eogo/internal/modules/user"
	"github.com/eogo-dev/eogo/routes"
	"github.com/gin-gonic/gin"
)

// SetupApp initializes the application for feature testing.
// Uses manual DI instead of Wire for test flexibility.
func SetupApp() *gin.Engine {
	// 1. Create Test Config
	cfg := &config.Config{}
	cfg.Server.Mode = "test"
	cfg.Database.Enabled = true
	cfg.Database.Driver = "sqlite"
	cfg.Database.Memory = true
	cfg.Database.MaxIdleConns = 1
	cfg.Database.MaxOpenConns = 1
	cfg.JWT.Secret = "testing-secret"
	cfg.JWT.Expire = time.Hour

	// 2. Initialize Database (In-Memory SQLite)
	db, err := database.NewDB(cfg)
	if err != nil {
		fmt.Printf("NewDB Error: %v\n", err)
		panic("failed to init test db: " + err.Error())
	}

	// 3. Run Migrations
	if err := bootstrap.RunMigrations(db); err != nil {
		fmt.Printf("RunMigrations Error: %v\n", err)
		panic("failed to run migrations for test db: " + err.Error())
	}

	// 4. Create Services via DI
	jwtService := jwt.NewService(cfg)
	emailService := email.NewService(cfg)

	// Set JWT service for middleware
	middleware.SetJWTService(jwtService)

	// 5. Create Repositories
	userRepo := user.NewRepository(db)
	permRepo := permission.NewRepository(db)

	// 6. Create Services
	userService := user.NewService(userRepo, jwtService)
	permService := permission.NewService(permRepo)

	// 7. Create Handlers
	handlers := &app.Handlers{
		User:       user.NewHandler(userService),
		Permission: permission.NewHandler(permService),
	}

	// 8. Build Application
	_ = &app.Application{
		Config:       cfg,
		DB:           db,
		JWTService:   jwtService,
		EmailService: emailService,
		Handlers:     handlers,
	}

	// 9. Setup Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	routes.Setup(r, handlers)

	return r
}

// NewTestCase is a shortcut to create a test case with the setup app
func NewTestCase(t *testing.T) *test_platform.TestCase {
	engine := SetupApp()
	return test_platform.NewTestCase(t, engine)
}
