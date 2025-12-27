package main

import (
	"os"

	"github.com/eogo-dev/eogo/internal/infra/console"
	"github.com/eogo-dev/eogo/internal/infra/console/commands"
)

const Version = "1.0.0"

func main() {
	// Initialize Console Application
	cli := console.New("eogo", Version)

	// Register Commands
	registerCommands(cli)

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

	// Register database commands
	dbMigrate := commands.NewDBMigrateCommand()
	app.Register(dbMigrate)
	app.RegisterAs("migrate", dbMigrate)

	dbFresh := commands.NewDBFreshCommand()
	app.Register(dbFresh)
	app.RegisterAs("migrate:fresh", dbFresh)

	dbRollback := commands.NewDBRollbackCommand()
	app.Register(dbRollback)
	app.RegisterAs("migrate:rollback", dbRollback)

	dbStatus := commands.NewDBStatusCommand()
	app.Register(dbStatus)
	app.RegisterAs("migrate:status", dbStatus)

	dbSeed := commands.NewDBSeedCommand()
	app.Register(dbSeed)
	app.RegisterAs("seed", dbSeed)

	// Register other commands
	app.Register(commands.NewServeCommand())
	app.Register(commands.NewEnvCommand())
	app.Register(commands.NewVersionCommand(Version))
	app.Register(commands.NewRouteListCommand())
}
