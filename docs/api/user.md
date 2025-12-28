// @Summary User registration
// @Description Create a new user account with the provided registration information
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserRegisterRequest true "Registration info"
// @Success 200 {object} User "User registered successfully"
// @Failure 400 {object} response.Response "Invalid request parameters or registration failed"
// @Router /users/register [post]

// @Summary User login
// @Description Authenticate user and return access token upon successful login
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserLoginRequest true "Login credentials"
// @Success 200 {object} User "Login successful, returns user data and token"
// @Failure 400 {object} response.Response "Invalid request or login failed"
// @Router /users/login [post]

// @Summary Update user profile
// @Description Update the currently authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserUpdateRequest true "User profile update data"
// @Success 200 {object} User "Profile updated successfully"
// @Failure 400 {object} response.Response "Invalid request or update failed"
// @Failure 401 {object} response.Response "Unauthorized, user not authenticated"
// @Router /users/profile [put]

// @Summary Change password
// @Description Change the currently authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserChangePasswordRequest true "Current and new password"
// @Success 200 {object} response.Response "Password changed successfully"
// @Failure 400 {object} response.Response "Invalid request or password change failed"
// @Failure 401 {object} response.Response "Unauthorized, user not authenticated"
// @Router /users/password [put]

// @Summary Reset password
// @Description Send a password reset email to the user's registered email address
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserPasswordResetRequest true "Email address for reset"
// @Success 200 {object} response.Response "Password reset email sent"
// @Failure 400 {object} response.Response "Invalid request or reset failed"
// @Router /users/password/reset [post]

// @Summary Get user profile
// @Description Retrieve the currently authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} User "User profile retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized, missing or invalid authentication"
// @Failure 500 {object} response.Response "Internal server error while retrieving profile"
// @Router /users/profile [get]

// @Summary Delete account
// @Description Permanently delete the currently authenticated user's account
// @Tags users
// @Success 200 {object} response.Response "Account deleted successfully"
// @Failure 401 {object} response.Response "Unauthorized, user not authenticated"
// @Failure 400 {object} response.Response "Failed to delete account"
// @Router /users/account [delete]

// @Summary Get user by ID
// @Description Retrieve user information by unique user ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User "User found and returned successfully"
// @Failure 400 {object} response.Response "Invalid user ID format"
// @Failure 404 {object} response.Response "User not found"
// @Router /users/{id} [get]

// @Summary Get user list
// @Description Retrieve a paginated list of users
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of users per page" default(10)
// @Success 200 {object} response.Response "Returns paginated list of users with total count"
// @Failure 500 {object} response.Response "Internal server error while fetching user list"
// @Router /users [get]

// @Summary Get user info
// @Description Retrieve detailed user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserInfo "User info retrieved successfully"
// @Failure 400 {object} response.Response "Invalid user ID format"
// @Failure 404 {object} response.Response "User not found"
// @Router /users/info/{id} [get]