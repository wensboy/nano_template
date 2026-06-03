package user

import (
	"example.com/nano_template/pkg/web/middleware"
	userpb "example.com/nano_template/proto/user"
)

type (
	RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		RoleID   int32  `json:"role_id"`
	}
	LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		RoleID   int32  `json:"role_id"`
	}
	LoginResponse struct {
		Token string `json:"token"`
	}
	AuthResponse struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		RoleID   uint   `json:"role_id"`
		Level    int    `json:"level"`
	}
	UpdatePasswordRequest struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	ListProfilesRequest struct {
		Keywords *string `form:"keywords"`
	}
	CreateProfileRequest struct {
		Avatar    string `json:"avatar"`
		Nickname  string `json:"nickname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Signature string `json:"signature"`
	}
	UpdateProfileRequest struct {
		Avatar    string `json:"avatar"`
		Nickname  string `json:"nickname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Signature string `json:"signature"`
	}
	GetProfileResponse struct {
		ID        uint   `json:"id"`
		UserID    uint   `json:"user_id"`
		Avatar    string `json:"avatar"`
		Nickname  string `json:"nickname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Signature string `json:"signature"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)

type UserModel interface {
	ToLoginResponse(*userpb.LoginResponse) LoginResponse
	ToAuthResponse(userID uint, username string, roleID uint, level int) AuthResponse
	ToGetProfileResponse(*userpb.GetUserProfileResponse) GetProfileResponse
	ToListProfilesResponse(*userpb.ListUserProfilesResponse) middleware.Pagination[GetProfileResponse]
}

type userModel struct{}

func NewUserModel() UserModel {
	return &userModel{}
}

func (*userModel) ToLoginResponse(resp *userpb.LoginResponse) LoginResponse {
	return LoginResponse{
		Token: resp.Token,
	}
}

func (*userModel) ToAuthResponse(userID uint, username string, roleID uint, level int) AuthResponse {
	return AuthResponse{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		Level:    level,
	}
}

func (*userModel) ToGetProfileResponse(resp *userpb.GetUserProfileResponse) GetProfileResponse {
	return GetProfileResponse{
		ID:        uint(resp.Id),
		UserID:    uint(resp.UserId),
		Avatar:    resp.Avatar,
		Nickname:  resp.Nickname,
		Email:     resp.Email,
		Phone:     resp.Phone,
		Signature: resp.Signature,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

func (m *userModel) ToListProfilesResponse(resp *userpb.ListUserProfilesResponse) middleware.Pagination[GetProfileResponse] {
	items := make([]GetProfileResponse, 0, len(resp.UserProfiles))
	for _, profile := range resp.UserProfiles {
		items = append(items, m.ToGetProfileResponse(profile))
	}
	return middleware.Pagination[GetProfileResponse]{
		Total: int(resp.Total),
		Items: items,
	}
}
