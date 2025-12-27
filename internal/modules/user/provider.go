package user

import (
	"github.com/eogo-dev/eogo/internal/domain"
	"github.com/google/wire"
)

// ProviderSet is the provider set for this module
// It binds concrete implementations to domain interfaces
var ProviderSet = wire.NewSet(
	NewRepository,
	wire.Bind(new(domain.UserRepository), new(*repository)),
	NewService,
	wire.Bind(new(Service), new(*service)),
	NewHandler,
)
