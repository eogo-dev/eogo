package infra

import (
	"github.com/eogo-dev/eogo/internal/infra/config"
	"github.com/eogo-dev/eogo/internal/infra/database"
	"github.com/eogo-dev/eogo/internal/infra/email"
	"github.com/eogo-dev/eogo/internal/infra/jwt"
	"github.com/google/wire"
)

// ProviderSet aggregates all infrastructure providers for Wire DI.
// This is the single source of truth for infrastructure dependencies.
var ProviderSet = wire.NewSet(
	// Config - loaded from environment
	config.Load,

	// Database - depends on Config
	database.NewDB,

	// JWT Service - depends on Config
	jwt.NewService,

	// Email Service - depends on Config
	email.NewService,
)
