package bootstrap

import (
	"log"

	"github.com/eogo-dev/eogo/database/seeders"
	_ "github.com/eogo-dev/eogo/database/seeders" // Import to trigger init()
)

// RunSeeders runs all registered database seeders
func RunSeeders() error {
	log.Println("Running database seeders")

	allSeeders := seeders.All()

	for _, seeder := range allSeeders {
		if err := seeder.Run(); err != nil {
			log.Printf("Seeder failed: %v", err)
			return err
		}
	}

	log.Printf("Successfully ran %d seeders", len(allSeeders))
	return nil
}
