package bootstrap

import (
	"log"
	"time"

	"github.com/eogo-dev/eogo/database/migrations"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// RunMigrations runs all migrations for the application
func RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migrations")
	startTime := time.Now()

	// Initialize the migrator with all migrations from database/migrations
	m := gormigrate.New(db, &gormigrate.Options{
		TableName:      "migrations",
		IDColumnName:   "id",
		IDColumnSize:   255,
		UseTransaction: true,
	}, migrations.All())

	// Execute migrations
	if err := m.Migrate(); err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Printf("Migration completed successfully in %v", time.Since(startTime))
	return nil
}
