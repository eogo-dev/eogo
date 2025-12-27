package user

import (
	"github.com/eogo-dev/eogo/internal/domain"
)

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
