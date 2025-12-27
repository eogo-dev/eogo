package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eogo-dev/eogo/internal/platform/email"
	"github.com/eogo-dev/eogo/internal/platform/jwt"
	"github.com/eogo-dev/eogo/internal/platform/logger"
	"github.com/eogo-dev/eogo/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the interface for user-related operations
type Service interface {
	Register(ctx context.Context, req *UserRegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *UserLoginRequest) (*UserLoginResponse, error)
	GetProfile(ctx context.Context, userID uint) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *UserUpdateRequest) (*UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *UserChangePasswordRequest) error
	ResetPassword(ctx context.Context, req *UserPasswordResetRequest) error
	DeleteAccount(ctx context.Context, userID uint) error
	GetByID(ctx context.Context, id uint) (*UserResponse, error)
	List(ctx context.Context, page, pageSize int) ([]*UserResponse, int64, error)
	GetUserByID(ctx context.Context, id uint) (*UserInfo, error)
}

// service implements the Service interface
type service struct {
	repo       Repository
	jwtService *jwt.Service
}

// NewService creates a new service instance
func NewService(repo Repository, jwtService *jwt.Service) Service {
	return &service{
		repo:       repo,
		jwtService: jwtService,
	}
}

// Create creates a new user
func (s *service) Create(ctx context.Context, model *User) error {
	return s.repo.Create(ctx, model)
}

// Update updates an existing user
func (s *service) Update(ctx context.Context, model *User) error {
	return s.repo.Update(ctx, model)
}

// Delete deletes a user by ID
func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// GetByID retrieves a user by ID
func (s *service) GetByID(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.ToResponse(user), nil
}

// List retrieves a paginated list of users
func (s *service) List(ctx context.Context, page, pageSize int) ([]*UserResponse, int64, error) {
	users, total, err := s.repo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var res []*UserResponse
	for _, u := range users {
		res = append(res, s.ToResponse(u))
	}
	return res, total, nil
}

// Register handles user registration
func (s *service) Register(ctx context.Context, req *UserRegisterRequest) (*UserResponse, error) {
	// Check if email already exists
	exists, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil && exists != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send welcome email
	if err := email.SendWelcomeEmail(user.Email, user.Username); err != nil {
		logger.Error("failed to send welcome email:", map[string]any{"error": err})
	}

	return s.ToResponse(user), nil
}

// Login handles user login
func (s *service) Login(ctx context.Context, req *UserLoginRequest) (*UserLoginResponse, error) {
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		user, err = s.repo.FindByEmail(ctx, req.Username)
		if err != nil {
			return nil, errors.New("invalid username or password")
		}
	}

	if user.Status == 0 {
		return nil, errors.New("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	now := time.Now()
	user.LastLogin = &now
	_ = s.repo.Update(ctx, user)

	return &UserLoginResponse{
		AccessToken: token,
		User:        user,
	}, nil
}

// GetProfile retrieves user profile
func (s *service) GetProfile(ctx context.Context, userID uint) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return s.ToResponse(user), nil
}

// UpdateProfile updates user profile
func (s *service) UpdateProfile(ctx context.Context, userID uint, req *UserUpdateRequest) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.ToResponse(user), nil
}

// ChangePassword changes user password
func (s *service) ChangePassword(ctx context.Context, userID uint, req *UserChangePasswordRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("incorrect old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	return s.repo.Update(ctx, user)
}

// ResetPassword resets user password via email
func (s *service) ResetPassword(ctx context.Context, req *UserPasswordResetRequest) error {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("email not found")
	}

	newPassword := utils.GenerateRandomString(12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return email.SendPasswordResetEmail(user.Email, newPassword)
}

// DeleteAccount deletes user account
func (s *service) DeleteAccount(ctx context.Context, userID uint) error {
	return s.repo.Delete(ctx, userID)
}

// GetUserByID retrieves user information for monitor/profile
func (s *service) GetUserByID(ctx context.Context, id uint) (*UserInfo, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Phone:     user.Phone,
		Bio:       user.Bio,
		Status:    user.Status,
		LastLogin: user.LastLogin,
	}, nil
}

func (s *service) ToResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Phone:     user.Phone,
		Bio:       user.Bio,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
