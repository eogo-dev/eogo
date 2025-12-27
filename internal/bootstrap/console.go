package bootstrap

import (
	"github.com/eogo-dev/eogo/internal/platform/console"
	"github.com/eogo-dev/eogo/internal/platform/console/commands"
)

// ConsoleKernel handles CLI commands
type ConsoleKernel struct {
	App *Application
	Cli *console.Application
}

// NewConsoleKernel creates a new Console kernel
func NewConsoleKernel(app *Application) *ConsoleKernel {
	cli := console.New("eogo", "1.0.0")

	// Register Commands
	registerCommands(cli)

	return &ConsoleKernel{
		App: app,
		Cli: cli,
	}
}

// Handle executes the console application
func (k *ConsoleKernel) Handle(args []string) error {
	return k.Cli.Run(args)
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
	app.Register(commands.NewVersionCommand("1.0.0"))
	app.Register(commands.NewRouteListCommand())
}
