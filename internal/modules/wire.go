//go:build wireinject
// +build wireinject

package app

import (
	"github.com/eogo-dev/eogo/internal/modules/permission"
	"github.com/eogo-dev/eogo/internal/modules/user"
	"github.com/eogo-dev/eogo/internal/platform/config"
	"github.com/eogo-dev/eogo/internal/platform/jwt"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// App holds all initialized handlers
type App struct {
	User       *user.Handler
	Permission *permission.Handler
}

// InitApp creates App with all dependencies wired
func InitApp(db *gorm.DB) (*App, error) {
	wire.Build(
		config.MustLoad,
		jwt.NewService,
		user.ProviderSet,
		permission.ProviderSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
