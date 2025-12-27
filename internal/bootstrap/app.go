package bootstrap

import (
	"github.com/eogo-dev/eogo/pkg/logger"
)

// InitLogger initializes the logger.
// Called before Wire initialization since logger is used during startup.
func InitLogger() {
	logger.Boot()
}
