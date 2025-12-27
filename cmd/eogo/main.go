package main

import (
	"os"

	"github.com/eogo-dev/eogo/internal/bootstrap"
)

const Version = "1.0.0"

func main() {
	// Initialize Application
	app := bootstrap.NewApp()

	// Initialize Console Kernel
	kernel := bootstrap.NewConsoleKernel(app)

	// Handle Command
	if err := kernel.Handle(os.Args); err != nil {
		os.Exit(1)
	}
}
