package migrations

import (
	"github.com/eogo-dev/eogo/internal/modules/user"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(&gormigrate.Migration{
		ID: "2025_06_18_000001_seed_default_users",
		Migrate: func(db *gorm.DB) error {
			var count int64
			db.Model(&user.UserPO{}).Count(&count)

			if count == 0 {
				adminUser := &user.UserPO{
					Username: "admin",
					Email:    "admin@example.com",
					Password: "hashed_password_here",
					Nickname: "Admin User",
					Status:   1,
				}
				return db.Create(adminUser).Error
			}
			return nil
		},
		Rollback: func(db *gorm.DB) error {
			return db.Where("username = ?", "admin").Delete(&user.UserPO{}).Error
		},
	})
}
