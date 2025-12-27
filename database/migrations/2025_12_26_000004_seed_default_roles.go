package migrations

import (
	"github.com/eogo-dev/eogo/internal/modules/permission"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(&gormigrate.Migration{
		ID: "2025_12_26_000004_seed_default_roles",
		Migrate: func(db *gorm.DB) error {
			roles := []permission.Role{
				{Name: "admin", DisplayName: "Administrator", Description: "Full access to all resources", IsDefault: false},
				{Name: "user", DisplayName: "User", Description: "Standard user access", IsDefault: true},
				{Name: "guest", DisplayName: "Guest", Description: "Read-only access", IsDefault: false},
			}
			for _, role := range roles {
				if err := db.FirstOrCreate(&role, permission.Role{Name: role.Name}).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(db *gorm.DB) error {
			return db.Where("name IN ?", []string{"admin", "user", "guest"}).Delete(&permission.Role{}).Error
		},
	})
}
