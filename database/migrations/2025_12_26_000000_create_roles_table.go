package migrations

import (
	"github.com/eogo-dev/eogo/internal/modules/permission"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(&gormigrate.Migration{
		ID: "2025_12_26_000000_create_roles_table",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&permission.Role{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable("roles")
		},
	})
}
