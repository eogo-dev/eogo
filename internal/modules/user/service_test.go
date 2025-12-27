package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eogo-dev/eogo/internal/platform/config"
	"github.com/eogo-dev/eogo/internal/platform/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// TestMain initializes dependencies for testing
func TestMain(m *testing.M) {
	cfg := &config.Config{}
	cfg.JWT.Secret = "test-secret"
	cfg.JWT.Expire = time.Hour
	jwt.Init(cfg)
	m.Run()
}

// MockUserRepository is a mock implementation of Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, u *User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, page, pageSize int) ([]*User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

// Test Cases

func TestUserService_GetProfile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, jwt.MustServiceInstance())
	ctx := context.Background()

	expectedUser := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	mockRepo.On("FindByID", ctx, uint(1)).Return(expectedUser, nil)

	result, err := service.GetProfile(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetProfile_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, jwt.MustServiceInstance())
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, errors.New("record not found"))

	result, err := service.GetProfile(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, jwt.MustServiceInstance())
	ctx := context.Background()

	req := &UserRegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	// Create should be called with matched user object. We use mock.AnythingOfType for *User
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	user, err := service.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	// Initialize JWT for test
	cfg := &config.Config{}
	cfg.JWT.Secret = "test-secret"
	cfg.JWT.Expire = time.Hour
	jwt.Init(cfg)

	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, jwt.MustServiceInstance())
	ctx := context.Background()

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &User{
		ID:       1,
		Username: "loginuser",
		Password: string(hashedPassword),
		Status:   1,
	}

	mockRepo.On("FindByUsername", ctx, "loginuser").Return(existingUser, nil)
	// Update last login
	mockRepo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	req := &UserLoginRequest{
		Username: "loginuser",
		Password: password,
	}

	resp, err := service.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, existingUser.ID, resp.User.ID)
	assert.NotEmpty(t, resp.AccessToken)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteAccount_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, jwt.MustServiceInstance())
	ctx := context.Background()

	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := service.DeleteAccount(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_List_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, jwt.MustServiceInstance())
	ctx := context.Background()

	users := []*User{
		{ID: 1, Username: "user1"},
		{ID: 2, Username: "user2"},
	}
	mockRepo.On("FindAll", ctx, 1, 10).Return(users, int64(2), nil)

	result, total, err := service.List(ctx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}
