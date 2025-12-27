package seeders

import (
	"github.com/eogo-dev/eogo/internal/modules/user"
	"github.com/eogo-dev/eogo/internal/platform/database"
)

type UserSeeder struct{}

func (s *UserSeeder) Run() error {
	db := database.GetDB()

	users := []user.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: secret
			Nickname: "Administrator",
			Status:   1,
		},
		{
			Username: "user",
			Email:    "user@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: secret
			Nickname: "Regular User",
			Status:   1,
		},
	}

	for _, u := range users {
		if err := db.FirstOrCreate(&u, user.User{Email: u.Email}).Error; err != nil {
			return err
		}
	}

	return nil
}

func init() {
	register(&UserSeeder{})
}
