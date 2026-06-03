package user

import (
	"net/http"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/rpc"
	"example.com/nano_template/pkg/web/middleware"
	userpb "example.com/nano_template/proto/user"
	"github.com/gin-gonic/gin"
)

// UserHandler defines the interface for user operations.
type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Auth(c *gin.Context)
	UpdatePassword(c *gin.Context)
	Cancel(c *gin.Context)
	ListProfiles(c *gin.Context)
	CreateProfile(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	DeleteProfile(c *gin.Context)
}

type userHandler struct {
	userClient        userpb.UserServiceClient
	userProfileClient userpb.UserProfileServiceClient
	model             UserModel
}

type UserHandlerOption func(*userHandler)

func NewUserHandler(opts ...UserHandlerOption) UserHandler {
	h := &userHandler{
		userClient:        rpc.GetUserServiceClient(),
		userProfileClient: rpc.GetUserProfileServiceClient(),
		model:             NewUserModel(),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Register godoc
// @Summary register user
// @Schemes
// @Description register a new user account
// @Tags user
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /user/register [post]
func (h *userHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}

	_, err := h.userClient.Register(c.Request.Context(), &userpb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		RoleId:   req.RoleID,
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.SuccNoMore(c, "User registered successfully")
}

// Login godoc
// @Summary login user
// @Schemes
// @Description authenticate user and return jwt token
// @Tags user
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} middleware.Response{data=LoginResponse}
// @Router /user/login [post]
func (h *userHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}

	resp, err := h.userClient.Login(c.Request.Context(), &userpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
		RoleId:   req.RoleID,
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	setTokenCookie(c, resp.Token)
	middleware.Succ(c, "Login successful", h.model.ToLoginResponse(resp))
}

// Logout godoc
// @Summary logout user
// @Schemes
// @Description user logout
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /user/logout [get]
func (h *userHandler) Logout(c *gin.Context) {
	middleware.UnauthedToCookie(c)
	middleware.SuccNoMore(c, "Logout successful")
}

// Auth godoc
// @Summary user auth check
// @Schemes
// @Description user auth check
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=AuthResponse}
// @Router /user/auth [get]
func (h *userHandler) Auth(c *gin.Context) {
	middleware.Succ(c, "auth user successful", h.model.ToAuthResponse(
		middleware.GetUserID(c),
		middleware.GetUsername(c),
		middleware.GetUserRole(c),
		middleware.GetRoleLevel(c),
	))
}

// UpdatePassword godoc
// @Summary update user password
// @Schemes
// @Description update password for current authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdatePasswordRequest true "Update password request"
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /user/password [put]
func (h *userHandler) UpdatePassword(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}

	_, err := h.userClient.UpdatePassword(c.Request.Context(), &userpb.UpdatePasswordRequest{
		UserId:      int32(userID),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.SuccNoMore(c, "Password updated successfully")
}

// Cancel godoc
// @Summary cancel user
// @Schemes
// @Description cancel current authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /user [delete]
func (h *userHandler) Cancel(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	_, err := h.userClient.Cancel(c.Request.Context(), &userpb.CancelRequest{UserId: int32(userID)})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.UnauthedToCookie(c)
	middleware.SuccNoMore(c, "User canceled successfully")
}

// ListProfiles godoc
// @Summary list user profiles
// @Schemes
// @Description retrieve user profiles with pagination
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10)"
// @Param keywords query string false "Profile keyword"
// @Success 200 {object} middleware.Response{data=middleware.Pagination[GetProfileResponse]}
// @Router /user/profiles [get]
func (h *userHandler) ListProfiles(c *gin.Context) {
	page, pageSize := middleware.ParsePageParams(c)
	var req ListProfilesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}

	resp, err := h.userProfileClient.List(c.Request.Context(), &userpb.ListUserProfilesRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Keywords: middleware.OptionalZero(req.Keywords),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	r := h.model.ToListProfilesResponse(resp)
	r.Page = page
	r.PageSize = pageSize
	middleware.Succ(c, "User profiles retrieved successfully", r)
}

// CreateProfile godoc
// @Summary create user profile
// @Schemes
// @Description create profile for current authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProfileRequest true "Create profile request"
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /user/profile [post]
func (h *userHandler) CreateProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}

	_, err := h.userProfileClient.Create(c.Request.Context(), &userpb.CreateUserProfileRequest{
		UserId:    int32(userID),
		Avatar:    req.Avatar,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Phone:     req.Phone,
		Signature: req.Signature,
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.SuccNoMore(c, "User profile created successfully")
}

// GetProfile godoc
// @Summary get user profile
// @Schemes
// @Description get current authenticated user profile
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} middleware.Response{data=GetProfileResponse}
// @Router /user/profile [get]
func (h *userHandler) GetProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	resp, err := h.userProfileClient.Get(c.Request.Context(), &userpb.GetUserProfileRequest{
		UserId: int32(userID),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.Succ(c, "User profile retrieved successfully", h.model.ToGetProfileResponse(resp))
}

// UpdateProfile godoc
// @Summary update user profile
// @Schemes
// @Description update profile fields for current authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Update profile request"
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /user/profile [put]
func (h *userHandler) UpdateProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}

	_, err := h.userProfileClient.Update(c.Request.Context(), &userpb.UpdateUserProfileRequest{
		UserId:    int32(userID),
		Avatar:    req.Avatar,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Phone:     req.Phone,
		Signature: req.Signature,
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.SuccNoMore(c, "User profile updated successfully")
}

// DeleteProfile godoc
// @Summary delete user profile
// @Schemes
// @Description delete current authenticated user's profile
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /user/profile [delete]
func (h *userHandler) DeleteProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	_, err := h.userProfileClient.Delete(c.Request.Context(), &userpb.DeleteUserProfileRequest{
		UserId: int32(userID),
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}

	middleware.SuccNoMore(c, "User profile deleted successfully")
}

func currentUserID(c *gin.Context) (uint, bool) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		middleware.Erro(c, http.StatusUnauthorized, "Invalid user ID")
		return 0, false
	}
	return userID, true
}

func setTokenCookie(c *gin.Context, token string) {
	cookieOpt := config.GetJwtConfig().CookieOption
	c.SetCookie(
		cookieOpt.AccessKey,
		token,
		cookieOpt.MaxAge,
		cookieOpt.Path,
		cookieOpt.Domain,
		cookieOpt.Secure,
		cookieOpt.HttpOnly,
	)
}
