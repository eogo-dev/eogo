package feature

import (
	"fmt"
	"testing"
	"time"

	"github.com/eogo-dev/eogo/internal/bootstrap"
	"github.com/eogo-dev/eogo/internal/platform/config"
	"github.com/eogo-dev/eogo/internal/platform/container"
	"github.com/eogo-dev/eogo/internal/platform/database"
	"github.com/eogo-dev/eogo/internal/platform/jwt"
	test_platform "github.com/eogo-dev/eogo/internal/platform/testing"
	"github.com/eogo-dev/eogo/routes"
	"github.com/gin-gonic/gin"
)

// SetupApp initializes the application for testing
func SetupApp() *gin.Engine {
	// 1. Create Test Config
	cfg := &config.Config{}
	cfg.Server.Mode = "test"
	cfg.Database.Enabled = true
	cfg.Database.Driver = "sqlite" // Use SQLite for speed
	cfg.Database.Memory = true     // In-memory
	cfg.Database.MaxIdleConns = 1
	cfg.Database.MaxOpenConns = 1
	// ... other configs if needed (JWT secret etc)
	cfg.JWT.Secret = "testing-secret"
	cfg.JWT.Expire = time.Hour

	container.App().Set(container.ServiceConfig, cfg)

	// 2. Initialize Dependencies
	jwt.Init(cfg)
	container.App().Set(container.ServiceJWT, jwt.MustServiceInstance())

	// 3. Initialize Database (In-Memory)
	// Note: We need to run Migrations here ideally.
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		fmt.Printf("InitDB Error: %v\n", err)
		panic("failed to init test db: " + err.Error())
	}
	// 4. Run Migrations
	if err := bootstrap.RunMigrations(db); err != nil {
		fmt.Printf("RunMigrations Error: %v\n", err)
		panic("failed to run migrations for test db: " + err.Error())
	}
	// db.AutoMigrate(&User{}) ... but User struct is in internal/modules/user.
	// Circular dependency risk if we import user module here?
	// Ideally, SetupApp calls a function that registers all migrations.

	container.App().Set(container.ServiceDB, db)

	// 4. Setup Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	routes.Setup(r)

	return r
}

// NewTestCase is a shortcut to create a test case with the setup app
func NewTestCase(t *testing.T) *test_platform.TestCase {
	app := SetupApp()
	return test_platform.NewTestCase(t, app)
}
