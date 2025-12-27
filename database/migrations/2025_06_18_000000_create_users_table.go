package migrations

import (
	"github.com/eogo-dev/eogo/internal/modules/user"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(&gormigrate.Migration{
		ID: "2025_06_18_000000_create_users_table",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&user.UserPO{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable("users")
		},
	})
}
