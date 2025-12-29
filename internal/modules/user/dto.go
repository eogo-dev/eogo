package user

import (
	"github.com/eogo-dev/eogo/internal/domain"
)

// ============================================================================
// Request DTOs (Input)
// ============================================================================

// UserRegisterRequest represents the registration request
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"max=50"`
	Phone    string `json:"phone" binding:"max=20"`
}

// UserLoginRequest represents the login request
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest represents the profile update request
type UserUpdateRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
	Phone    string `json:"phone" binding:"max=20"`
	Bio      string `json:"bio" binding:"max=500"`
}

// UserChangePasswordRequest represents the password change request
type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// UserPasswordResetRequest represents the password reset request
type UserPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ============================================================================
// Response DTOs (Output)
// ============================================================================

// UserLoginResponse represents the login response
type UserLoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        *domain.User `json:"user"` // Domain直接输出，Password自动隐藏
}

// ============================================================================
// Model Mappers (UserPO <-> domain.User)
// ============================================================================

// toDomain converts UserPO to domain.User
func (po *UserPO) toDomain() *domain.User {
	if po == nil {
		return nil
	}
	return &domain.User{
		ID:        po.ID,
		Username:  po.Username,
		Email:     po.Email,
		Password:  po.Password,
		Nickname:  po.Nickname,
		Avatar:    po.Avatar,
		Phone:     po.Phone,
		Bio:       po.Bio,
		Status:    po.Status,
		LastLogin: po.LastLogin,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

// toUserPO converts domain.User to UserPO for database operations
func toUserPO(u *domain.User) *UserPO {
	if u == nil {
		return nil
	}
	return &UserPO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		Bio:       u.Bio,
		Status:    u.Status,
		LastLogin: u.LastLogin,
	}
}

// toDomainList converts a slice of UserPO to domain.User slice
func toDomainList(poList []*UserPO) []*domain.User {
	result := make([]*domain.User, len(poList))
	for i, po := range poList {
		result[i] = po.toDomain()
	}
	return result
}
