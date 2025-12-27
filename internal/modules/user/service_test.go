package user

import (
	"context"
	"errors"
	"testing"

	"github.com/eogo-dev/eogo/internal/domain"
	"github.com/eogo-dev/eogo/internal/infra/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// testJWTService is a shared JWT service instance for all tests in this package.
// Uses jwt.NewTestService() which provides a pre-configured test service.
var testJWTService = jwt.NewTestService()

// MockUserRepository is a mock implementation of domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, page, pageSize int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Test Cases

func TestUserService_GetProfile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, testJWTService)
	ctx := context.Background()

	expectedUser := &domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Status:   1,
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
	service := NewService(mockRepo, testJWTService)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, uint(999)).Return(nil, errors.New("record not found"))

	result, err := service.GetProfile(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, testJWTService)
	ctx := context.Background()

	req := &UserRegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	// Create should be called with matched user object. We use mock.AnythingOfType for *domain.User
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, testJWTService)
	ctx := context.Background()

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &domain.User{
		ID:       1,
		Username: "loginuser",
		Password: string(hashedPassword),
		Status:   1,
	}

	mockRepo.On("FindByUsername", ctx, "loginuser").Return(existingUser, nil)
	// Update last login
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

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
	service := NewService(mockRepo, testJWTService)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := service.DeleteAccount(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_List_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, testJWTService)
	ctx := context.Background()

	users := []*domain.User{
		{ID: 1, Username: "user1", Status: 1},
		{ID: 2, Username: "user2", Status: 1},
	}
	mockRepo.On("FindAll", ctx, 1, 10).Return(users, int64(2), nil)

	result, total, err := service.List(ctx, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}
