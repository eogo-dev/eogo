package user

import (
	"strconv"

	"github.com/eogo-dev/eogo/pkg/pagination"
	"github.com/eogo-dev/eogo/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handler handles user-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ============================================================================
// Authentication
// ============================================================================

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, "Registration failed", err)
		return
	}

	response.Created(c, user) // Domain直接输出，Password自动隐藏
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		response.HandleError(c, "Login failed", err)
		return
	}

	response.Success(c, resp)
}

// ============================================================================
// Profile (Authenticated User)
// ============================================================================

// GetProfile gets current user's profile
func (h *Handler) GetProfile(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	user, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, "Failed to get profile", err)
		return
	}

	response.Success(c, user) // 直接输出
}

// UpdateProfile updates current user's profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	user, err := h.service.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		response.HandleError(c, "Failed to update profile", err)
		return
	}

	response.Success(c, user)
}

// ChangePassword changes current user's password
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		response.HandleError(c, "Failed to change password", err)
		return
	}

	response.Success(c, gin.H{"message": "Password changed successfully"})
}

// DeleteAccount deletes current user's account
func (h *Handler) DeleteAccount(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	if err := h.service.DeleteAccount(c.Request.Context(), userID); err != nil {
		response.HandleError(c, "Failed to delete account", err)
		return
	}

	response.NoContent(c)
}

// ============================================================================
// Public
// ============================================================================

// ResetPassword initiates password reset
func (h *Handler) ResetPassword(c *gin.Context) {
	var req UserPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters", err)
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), &req); err != nil {
		response.HandleError(c, "Failed to reset password", err)
		return
	}

	response.Success(c, gin.H{"message": "Password reset email sent"})
}

// ============================================================================
// Admin/Query
// ============================================================================

// Get gets user by ID
func (h *Handler) Get(c *gin.Context) {
	id, err := h.parseID(c, "id")
	if err != nil {
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, "User not found", err)
		return
	}

	response.Success(c, user) // 直接输出
}

// List gets paginated user list
func (h *Handler) List(c *gin.Context) {
	req := pagination.FromContext(c)

	users, total, err := h.service.List(c.Request.Context(), req.GetPage(), req.GetPerPage())
	if err != nil {
		response.HandleError(c, "Failed to get user list", err)
		return
	}

	paginator := pagination.NewPaginator(users, total, req.GetPage(), req.GetPerPage())
	paginator.SetPath(c.Request.URL.Path)

	response.Success(c, paginator) // 统一用 Success，自动检测分页！
}

// GetUserInfo gets detailed user info by ID (alias for Get)
func (h *Handler) GetUserInfo(c *gin.Context) {
	h.Get(c)
}

// ============================================================================
// Helpers
// ============================================================================

func (h *Handler) getUserID(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c)
		return 0, false
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		response.InternalServerError(c, "Invalid user ID type", nil)
		return 0, false
	}

	return userID, true
}

func (h *Handler) parseID(c *gin.Context, param string) (uint, error) {
	idStr := c.Param(param)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid ID", err)
		return 0, err
	}
	return uint(id), nil
}
