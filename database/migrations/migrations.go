package migrations

import "github.com/go-gormigrate/gormigrate/v2"

var registry []*gormigrate.Migration

// register adds a migration to the registry
// This is called by init() functions in migration files
func register(m *gormigrate.Migration) {
	registry = append(registry, m)
}

// All returns all registered migrations in order
func All() []*gormigrate.Migration {
	return registry
}
