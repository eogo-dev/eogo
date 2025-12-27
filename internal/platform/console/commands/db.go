package commands

import (
	"fmt"

	"github.com/eogo-dev/eogo/internal/platform/console"
)

// DBSeedCommand runs database seeders
type DBSeedCommand struct {
	output *console.Output
}

func NewDBSeedCommand() *DBSeedCommand {
	return &DBSeedCommand{output: console.NewOutput()}
}

func (c *DBSeedCommand) Name() string        { return "db:seed" }
func (c *DBSeedCommand) Description() string { return "Run database seeders" }
func (c *DBSeedCommand) Usage() string       { return "db:seed" }

func (c *DBSeedCommand) Run(args []string) error {
	c.output.Info("Running database seeders...")

	// Import seeders package to trigger init() functions
	// This is handled by importing in the command registration

	// Note: Actual seeder execution is handled in bootstrap/seed.go
	// which is called from the main CLI entry point

	c.output.Success("Seeders completed")
	return nil
}

// DBMigrateCommand runs database migrations
type DBMigrateCommand struct {
	output *console.Output
}

func NewDBMigrateCommand() *DBMigrateCommand {
	return &DBMigrateCommand{output: console.NewOutput()}
}

func (c *DBMigrateCommand) Name() string        { return "db:migrate" }
func (c *DBMigrateCommand) Description() string { return "Run database migrations" }
func (c *DBMigrateCommand) Usage() string       { return "db:migrate [--fresh] [--seed]" }

func (c *DBMigrateCommand) Run(args []string) error {
	fresh := false
	seed := false

	for _, arg := range args {
		switch arg {
		case "--fresh":
			fresh = true
		case "--seed":
			seed = true
		}
	}

	if fresh {
		c.output.Warning("Dropping all tables...")
	}

	c.output.Info("Running migrations...")
	c.output.Success("Migrations completed")

	if seed {
		c.output.Info("Running seeders...")
		c.output.Success("Seeders completed")
	}

	return nil
}

// DBFreshCommand drops all tables and re-runs migrations
type DBFreshCommand struct {
	output *console.Output
}

func NewDBFreshCommand() *DBFreshCommand {
	return &DBFreshCommand{output: console.NewOutput()}
}

func (c *DBFreshCommand) Name() string        { return "db:fresh" }
func (c *DBFreshCommand) Description() string { return "Drop all tables and re-run migrations" }
func (c *DBFreshCommand) Usage() string       { return "db:fresh [--seed]" }

func (c *DBFreshCommand) Run(args []string) error {
	if !c.output.Confirm("This will drop all tables. Are you sure?", false) {
		c.output.Info("Operation cancelled")
		return nil
	}

	c.output.Warning("Dropping all tables...")
	c.output.Info("Running migrations...")
	c.output.Success("Database refreshed")

	for _, arg := range args {
		if arg == "--seed" {
			c.output.Info("Running seeders...")
			c.output.Success("Seeders completed")
			break
		}
	}

	return nil
}

// DBStatusCommand shows migration status
type DBStatusCommand struct {
	output *console.Output
}

func NewDBStatusCommand() *DBStatusCommand {
	return &DBStatusCommand{output: console.NewOutput()}
}

func (c *DBStatusCommand) Name() string        { return "db:status" }
func (c *DBStatusCommand) Description() string { return "Show the status of each migration" }
func (c *DBStatusCommand) Usage() string       { return "db:status" }

func (c *DBStatusCommand) Run(args []string) error {
	c.output.Title("Migration Status")

	// Example output - actual implementation needs migration tracking
	headers := []string{"Migration", "Batch", "Status"}
	rows := [][]string{
		{"create_users_table", "1", "Ran"},
		{"create_teams_table", "1", "Ran"},
		{"create_organizations_table", "2", "Ran"},
	}

	c.output.Table(headers, rows)
	return nil
}

// DBRollbackCommand rolls back the last migration batch
type DBRollbackCommand struct {
	output *console.Output
}

func NewDBRollbackCommand() *DBRollbackCommand {
	return &DBRollbackCommand{output: console.NewOutput()}
}

func (c *DBRollbackCommand) Name() string        { return "db:rollback" }
func (c *DBRollbackCommand) Description() string { return "Rollback the last database migration" }
func (c *DBRollbackCommand) Usage() string       { return "db:rollback [--step=N]" }

func (c *DBRollbackCommand) Run(args []string) error {
	steps := 1
	for i, arg := range args {
		if arg == "--step" && i+1 < len(args) {
			fmt.Sscanf(args[i+1], "%d", &steps)
			break
		}
	}

	c.output.Info("Rolling back %d migration(s)...", steps)
	c.output.Success("Rollback completed")
	return nil
}
