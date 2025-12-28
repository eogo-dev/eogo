package main

import (
	"os"
	"strings"

	"github.com/eogo-dev/eogo/internal/infra/console"
	"github.com/eogo-dev/eogo/internal/infra/console/commands"
	"github.com/eogo-dev/eogo/internal/infra/plugin"
)

const Version = "1.0.0"

func main() {
	// Initialize Console Application
	cli := console.New("eogo", Version)

	// Register Commands
	registerCommands(cli)

	// Check if first argument is a plugin command
	if len(os.Args) > 1 && isPluginCommand(os.Args[1]) {
		pluginName := os.Args[1]
		pluginArgs := os.Args[2:]

		// Execute plugin
		if err := plugin.Execute(pluginName, pluginArgs); err != nil {
			os.Exit(1)
		}
		return
	}

	// Handle Command
	if err := cli.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func registerCommands(app *console.Application) {
	// Register make commands
	app.Register(commands.NewMakeModelCommand())
	app.Register(commands.NewMakeServiceCommand())
	app.Register(commands.NewMakeHandlerCommand())
	app.Register(commands.NewMakeRepositoryCommand())
	app.Register(commands.NewMakeSeederCommand())
	app.Register(commands.NewMakeMigrationCommand())
	app.Register(commands.NewMakeModuleCommand())

	// Register database migration commands (new Migrator-based)
	dbMigrate := commands.NewMigrateCommand()
	app.Register(dbMigrate)
	app.RegisterAs("migrate", dbMigrate)

	dbFresh := commands.NewFreshCommand()
	app.Register(dbFresh)
	app.RegisterAs("migrate:fresh", dbFresh)

	dbRollback := commands.NewRollbackCommand()
	app.Register(dbRollback)
	app.RegisterAs("migrate:rollback", dbRollback)

	dbReset := commands.NewResetCommand()
	app.Register(dbReset)
	app.RegisterAs("migrate:reset", dbReset)

	dbStatus := commands.NewStatusCommand()
	app.Register(dbStatus)
	app.RegisterAs("migrate:status", dbStatus)

	// Register seed command
	dbSeed := commands.NewDBSeedCommand()
	app.Register(dbSeed)
	app.RegisterAs("seed", dbSeed)

	// Register other commands
	app.Register(commands.NewServeCommand())
	app.Register(commands.NewEnvCommand())
	app.Register(commands.NewVersionCommand(Version))
	app.Register(commands.NewRouteListCommand())

	// Register plugin commands
	app.Register(commands.NewPluginListCommand())
}

// isPluginCommand checks if a command is a plugin command
func isPluginCommand(cmd string) bool {
	// Skip if it's a known core command
	coreCommands := map[string]bool{
		"make:model":       true,
		"make:service":     true,
		"make:handler":     true,
		"make:repository":  true,
		"make:seeder":      true,
		"make:migration":   true,
		"make:module":      true,
		"migrate":          true,
		"migrate:fresh":    true,
		"migrate:rollback": true,
		"migrate:reset":    true,
		"migrate:status":   true,
		"db:migrate":       true,
		"db:fresh":         true,
		"db:rollback":      true,
		"db:reset":         true,
		"db:status":        true,
		"db:seed":          true,
		"seed":             true,
		"serve":            true,
		"env":              true,
		"version":          true,
		"route:list":       true,
		"plugin:list":      true,
		"help":             true,
	}

	if coreCommands[cmd] {
		return false
	}

	// Check if it starts with a known prefix
	if strings.HasPrefix(cmd, "make:") ||
		strings.HasPrefix(cmd, "migrate:") ||
		strings.HasPrefix(cmd, "db:") ||
		strings.HasPrefix(cmd, "route:") ||
		strings.HasPrefix(cmd, "plugin:") {
		return false
	}

	// Check if plugin exists
	return plugin.IsInstalled(cmd)
}
