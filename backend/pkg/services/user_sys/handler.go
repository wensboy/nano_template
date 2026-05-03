package userSys

import (
	"net/http"

	"example.com/nano_template/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// UserHandlerInterface defines the interface for user handler operations.
type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserDetails(c *gin.Context)
	ChangePassword(c *gin.Context)
	UpdateUserProfile(c *gin.Context)
	DeactivateUser(c *gin.Context)
}

// userHandler handles user-related HTTP requests.
type userHandler struct {
	userService UserService
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdateUserProfileRequest struct {
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Signature string `json:"signature"`
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userService UserService) UserHandler {
	return &userHandler{userService: userService}
}

// Register godoc
// @Summary register user
// @Schemes
// @Description register a new user account
// @Tags user
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 200 {object} middleware.Response
// @Router /user/register [post]
// Register handles user registration.
func (h *userHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}

	user, err := h.userService.RegisterUser(req.Username, req.Password)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.Succ(c, "User registered successfully", gin.H{"user_id": user.ID})
}

// Login godoc
// @Summary login user
// @Schemes
// @Description authenticate user and return jwt token
// @Tags user
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} middleware.Response
// @Router /user/login [post]
// Login handles user login.
func (h *userHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}

	user, err := h.userService.LoginUser(req.Username, req.Password)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	token, err := middleware.GenerateJWT("", user.ID, user.Username, 0)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.Succ(c, "Login successful", gin.H{"token": token})
}

// GetUserDetails godoc
// @Summary get user details
// @Schemes
// @Description get current authenticated user details
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response
// @Router /user/details [get]
// GetUserDetails handles retrieving user details.
func (h *userHandler) GetUserDetails(c *gin.Context) {
	// userIDStr := c.Param("user_id")
	// userID, err := strconv.ParseUint(userIDStr, 10, 32)
	username := middleware.GetUsernameFromContext(c)
	if username == "" {
		middleware.Fail(c, "Invalid user ID")
		return
	}

	user, profile, err := h.userService.GetUserDetails(username)
	if err != nil {
		middleware.Fail(c, "User not found")
		return
	}

	uwp := UserWithProfile{
		UserID:    user.ID,
		Username:  user.Username,
		Avatar:    profile.Avatar,
		Nickname:  profile.Nickname,
		Email:     profile.Email,
		Gender:    profile.Gender,
		Signature: profile.Signature,
	}

	middleware.Succ(c, "User details retrieved", uwp)
}

// ChangePassword godoc
// @Summary change user password
// @Schemes
// @Description change password for current authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Change password request"
// @Success 200 {object} middleware.Response
// @Router /user/update/password [put]
// ChangePassword handles password change.
func (h *userHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		middleware.Fail(c, "Invalid user ID")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}

	err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.Succ(c, "Password changed successfully", nil)
}

// UpdateUserProfile godoc
// @Summary update user profile
// @Schemes
// @Description update profile fields for current authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param request body UpdateUserProfileRequest true "Update user profile request"
// @Success 200 {object} middleware.Response
// @Router /user/update/profile [put]
// UpdateUserProfile handles profile updates.
func (h *userHandler) UpdateUserProfile(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		middleware.Fail(c, "Invalid user ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		middleware.Fail(c, "Invalid request data")
		return
	}

	err := h.userService.UpdateUserProfile(userID, updates)
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.Succ(c, "Profile updated successfully", nil)
}

// DeactivateUser godoc
// @Summary deactivate user
// @Schemes
// @Description deactivate current authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response
// @Router /user/delete [delete]
// DeactivateUser handles user deactivation.
func (h *userHandler) DeactivateUser(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		middleware.Fail(c, "Invalid user ID")
		return
	}

	err := h.userService.DeactivateUser(userID)
	if err != nil {
		middleware.Erro(c, http.StatusInternalServerError, "Failed to deactivate user")
		return
	}

	middleware.Succ(c, "User deactivated successfully", nil)
}
