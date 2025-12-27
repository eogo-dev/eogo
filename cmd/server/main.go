package main

import (
	"github.com/eogo-dev/eogo/internal/bootstrap"
)

// @title Eogo API
// @version 1.0
// @description An elegant Go web framework for production-ready development
// @host localhost:8025
// @BasePath /v1

func main() {
	// Initialize Application
	app := bootstrap.NewApp()

	// Initialize HTTP Kernel
	kernel := bootstrap.NewHttpKernel(app)

	// Handle Request
	kernel.Handle()
}
