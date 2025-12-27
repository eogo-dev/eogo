package domain

import (
	"context"
	"time"
)

// User represents the pure domain entity (no ORM tags, no JSON tags)
// This is the core business object that all layers work with
type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string // hashed password
	Nickname  string
	Avatar    string
	Phone     string
	Bio       string
	Status    int // 1: active, 0: disabled
	LastLogin *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsActive returns whether the user account is active
func (u *User) IsActive() bool {
	return u.Status == 1
}

// UserRepository defines the contract for user data operations
// Implementations live in modules/user/repository.go
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*User, int64, error)
}
