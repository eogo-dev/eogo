package app

import (
	"github.com/eogo-dev/eogo/internal/infra/config"
	"github.com/eogo-dev/eogo/internal/infra/email"
	"github.com/eogo-dev/eogo/internal/infra/jwt"
	"github.com/eogo-dev/eogo/internal/modules/permission"
	"github.com/eogo-dev/eogo/internal/modules/user"
	"gorm.io/gorm"
)

// Application holds all application dependencies injected via Wire.
// This is the root container for the entire application.
type Application struct {
	Config       *config.Config
	DB           *gorm.DB
	JWTService   *jwt.Service
	EmailService *email.Service
	Handlers     *Handlers
}

// Handlers holds all HTTP handlers for modules.
type Handlers struct {
	User       *user.Handler
	Permission *permission.Handler
}
