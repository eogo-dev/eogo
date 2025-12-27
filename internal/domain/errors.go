package domain

import "errors"

// Domain-specific errors
// These are business errors that can be returned by any layer
var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrAccountDisabled    = errors.New("account is disabled")

	// Permission errors
	ErrPermissionDenied = errors.New("permission denied")
	ErrRoleNotFound     = errors.New("role not found")

	// Generic errors
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource already exists")
	ErrInvalidInput = errors.New("invalid input")
)
