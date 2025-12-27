package jwt

import (
	"fmt"
	"time"

	"github.com/eogo-dev/eogo/internal/platform/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	defaultService *Service
)

// Service provides JWT helpers bound to a configuration instance.
type Service struct {
	cfg *config.Config
}

// NewService constructs a JWT service using the provided configuration.
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// Init initializes the JWT service
func Init(c *config.Config) {
	defaultService = NewService(c)
}

// SetDefaultService overrides the global JWT service used by package-level helpers.
func SetDefaultService(service *Service) {
	defaultService = service
}

// ServiceInstance returns the currently configured global JWT service.
func ServiceInstance() (*Service, error) {
	if defaultService == nil {
		return nil, fmt.Errorf("jwt service not initialized")
	}
	return defaultService, nil
}

// MustServiceInstance returns the global service or panics if not initialized.
func MustServiceInstance() *Service {
	svc, err := ServiceInstance()
	if err != nil {
		panic(err)
	}
	return svc
}

// Claims represents custom JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token
func (s *Service) GenerateToken(userID uint, username string) (string, error) {
	if s == nil || s.cfg == nil {
		return "", fmt.Errorf("jwt service not initialized")
	}

	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWT.ExpireDuration())),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

// ParseToken parses and validates a JWT token
func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	if s == nil || s.cfg == nil {
		return nil, fmt.Errorf("jwt service not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateToken generates a JWT token using the global service.
func GenerateToken(userID uint, username string) (string, error) {
	svc, err := ServiceInstance()
	if err != nil {
		return "", err
	}
	return svc.GenerateToken(userID, username)
}

// ParseToken parses a JWT token using the global service.
func ParseToken(tokenString string) (*Claims, error) {
	svc, err := ServiceInstance()
	if err != nil {
		return nil, err
	}
	return svc.ParseToken(tokenString)
}
