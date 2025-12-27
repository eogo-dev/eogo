package seeders

// Seeder interface defines the contract for database seeders
type Seeder interface {
	Run() error
}

var registry []Seeder

// register adds a seeder to the registry
func register(s Seeder) {
	registry = append(registry, s)
}

// All returns all registered seeders
func All() []Seeder {
	return registry
}
