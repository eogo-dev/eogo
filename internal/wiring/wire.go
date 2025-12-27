//go:build wireinject
// +build wireinject

package wiring

import (
	"github.com/eogo-dev/eogo/internal/app"
	"github.com/eogo-dev/eogo/internal/infra"
	"github.com/eogo-dev/eogo/internal/modules/permission"
	"github.com/eogo-dev/eogo/internal/modules/user"
	"github.com/google/wire"
)

// InitApplication initializes the entire application with all dependencies.
// This is the single entry point for Wire DI.
func InitApplication() (*app.Application, error) {
	wire.Build(
		// Infrastructure providers
		infra.ProviderSet,

		// Module providers
		user.ProviderSet,
		permission.ProviderSet,

		// Aggregate handlers
		wire.Struct(new(app.Handlers), "*"),

		// Build final application
		wire.Struct(new(app.Application), "*"),
	)
	return nil, nil
}
