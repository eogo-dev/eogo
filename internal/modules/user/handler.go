package user

import (
	"strconv"

	"github.com/eogo-dev/eogo/pkg/logger"
	"github.com/eogo-dev/eogo/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handler handles user-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register handles user registration
// @Summary User registration
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserRegisterRequest true "Registration info"
// @Success 200 {object} User
// @Router /users/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, "Registration failed", err)
		return
	}

	response.Success(c, user)
}

// Login handles user login
// @Summary User login
// @Description Login and get access token
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserLoginRequest true "Login credentials"
// @Success 200 {object} User
// @Router /users/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, "Login failed", err)
		return
	}

	response.Success(c, resp)
}

// UpdateProfile updates user profile
// @Summary Update user profile
// @Description Update current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserUpdateRequest true "User info"
// @Success 200 {object} User
// @Router /users/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uint)

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	user, err := h.service.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		response.BadRequest(c, "Failed to update profile", err)
		return
	}

	response.Success(c, user)
}

// ChangePassword changes user password
// @Summary Change password
// @Description Change current user's password
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserChangePasswordRequest true "Password info"
// @Success 200 {string} string "Password changed successfully"
// @Router /users/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uint)

	var req UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		response.BadRequest(c, "Failed to change password", err)
		return
	}

	response.Success(c, gin.H{"message": "Password changed successfully"})
}

// ResetPassword resets user password
// @Summary Reset password
// @Description Reset password via email
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserPasswordResetRequest true "Email info"
// @Success 200 {string} string "Password reset email sent"
// @Router /users/password/reset [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req UserPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), &req); err != nil {
		response.BadRequest(c, "Failed to reset password", err)
		return
	}

	response.Success(c, gin.H{"message": "Password reset email sent"})
}

// GetProfile gets user profile
// @Summary Get user profile
// @Description Get current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} User
// @Router /users/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		response.InternalServerError(c, "Invalid user ID type", nil)
		return
	}

	user, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "Failed to get profile", err)
		return
	}

	response.Success(c, user)
}

// DeleteAccount deletes user account
// @Summary Delete account
// @Description Delete current user's account
// @Tags users
// @Success 200 {string} string "Account deleted"
// @Router /users/account [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uint)

	if err := h.service.DeleteAccount(c.Request.Context(), userID); err != nil {
		response.BadRequest(c, "Failed to delete account", err)
		return
	}

	response.Success(c, gin.H{"message": "Account deleted"})
}

// Get gets user by ID
// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /users/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err)
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "User not found", err)
		return
	}

	response.Success(c, user)
}

// List gets user list
// @Summary Get user list
// @Description Get paginated user list
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {array} User
// @Router /users [get]
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		logger.Error("Failed to get user list", map[string]any{"error": err})
		response.InternalServerError(c, "Failed to get user list", err)
		return
	}

	response.Success(c, gin.H{"total": total, "list": users})
}

// GetUserInfo gets user info
// @Summary Get user info
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserInfo
// @Router /users/info/{id} [get]
func (h *Handler) GetUserInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", err)
		return
	}

	userInfo, err := h.service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "User not found", err)
		return
	}

	response.Success(c, userInfo)
}
