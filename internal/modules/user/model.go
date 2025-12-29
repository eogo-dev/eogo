package user

import (
	"time"

	"gorm.io/gorm"
)

// UserPO is the persistent object for database operations (internal to repository)
// This struct contains GORM tags and is NOT exposed outside the repository layer
type UserPO struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username  string         `gorm:"size:50;not null" json:"username"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	Email     string         `gorm:"size:100;not null;unique" json:"email"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Bio       string         `gorm:"size:500" json:"bio"`
	Status    int            `gorm:"default:1" json:"status"` // 1: active, 0: disabled
	LastLogin *time.Time     `json:"last_login"`
}

// TableName specifies the database table name
func (UserPO) TableName() string {
	return "users"
}
