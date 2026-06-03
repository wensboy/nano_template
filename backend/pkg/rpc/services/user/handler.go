package user

import (
	"context"
	"time"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/rpc/middleware"
	userpb "example.com/nano_template/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MountUserServer(s *grpc.Server, cfg *config.Config) {
	userService := NewUserService(config.GetGDB())
	handler := &userHandler{cfg: cfg, service: userService, model: NewUserModel()}
	userpb.RegisterUserServiceServer(s, handler)
	userpb.RegisterUserProfileServiceServer(s, handler)
}

type userHandler struct {
	cfg *config.Config
	userpb.UnimplementedUserServiceServer
	userpb.UnimplementedUserProfileServiceServer
	service UserService
	model   UserModel
}

/* 用户基本服务 */

func (h *userHandler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	_, _, err := h.service.Register(req.Username, req.Password, uint(req.RoleId))
	if err != nil {
		return nil, err
	}

	return &userpb.RegisterResponse{}, nil
}

func (h *userHandler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	user, role, err := h.service.Login(req.Username, req.Password, uint(req.RoleId))
	if err != nil {
		return nil, err
	}

	token, err := middleware.GenerateJWT(h.cfg.JwtConfig.Secret, middleware.UserClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserRole:  role.ID,
		RoleLevel: role.Level,
	}, time.Duration(h.cfg.JwtConfig.TTL)*time.Second)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &userpb.LoginResponse{
		Token: token,
	}, nil
}

func (h *userHandler) UpdatePassword(ctx context.Context, req *userpb.UpdatePasswordRequest) (*userpb.UpdatePasswordResponse, error) {
	err := h.service.UpdatePassword(uint(req.UserId), req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &userpb.UpdatePasswordResponse{}, nil
}

func (h *userHandler) Cancel(ctx context.Context, req *userpb.CancelRequest) (*userpb.CancelResponse, error) {
	return &userpb.CancelResponse{}, h.service.Cancel(uint(req.UserId))
}

/* 用户首选项服务 */

func (h *userHandler) Create(ctx context.Context, req *userpb.CreateUserProfileRequest) (*userpb.CreateUserProfileResponse, error) {
	_, err := h.service.CreateProfile(
		uint(req.UserId),
		middleware.OptionalString(req.Avatar),
		middleware.OptionalString(req.Nickname),
		middleware.OptionalString(req.Email),
		middleware.OptionalString(req.Phone),
		middleware.OptionalString(req.Signature),
	)
	if err != nil {
		return nil, err
	}
	return &userpb.CreateUserProfileResponse{}, nil
}

func (h *userHandler) Get(ctx context.Context, req *userpb.GetUserProfileRequest) (*userpb.GetUserProfileResponse, error) {
	userId := middleware.OptionalUint(uint(req.UserId))
	id := middleware.OptionalUint(uint(req.Id))

	profile, err := h.service.GetProfile(userId, id)
	if err != nil {
		return nil, err
	}
	return h.model.ToPbUserProfile(profile), nil
}

func (h *userHandler) List(ctx context.Context, req *userpb.ListUserProfilesRequest) (*userpb.ListUserProfilesResponse, error) {
	page := middleware.OptionalInt(int(req.Page))
	pageSize := middleware.OptionalInt(int(req.PageSize))

	profiles, total, err := h.service.ListProfiles(page, pageSize, middleware.OptionalString(req.Keywords))
	if err != nil {
		return nil, err
	}

	respProfiles := make([]*userpb.GetUserProfileResponse, 0, len(profiles))
	for _, profile := range profiles {
		respProfiles = append(respProfiles, h.model.ToPbUserProfile(profile))
	}
	return &userpb.ListUserProfilesResponse{
		UserProfiles: respProfiles,
		Total:        total,
	}, nil
}

func (h *userHandler) Update(ctx context.Context, req *userpb.UpdateUserProfileRequest) (*userpb.UpdateUserProfileResponse, error) {
	userId := middleware.OptionalUint(uint(req.UserId))
	id := middleware.OptionalUint(uint(req.Id))

	_, err := h.service.UpdateProfile(
		userId,
		id,
		middleware.OptionalString(req.Avatar),
		middleware.OptionalString(req.Nickname),
		middleware.OptionalString(req.Email),
		middleware.OptionalString(req.Phone),
		middleware.OptionalString(req.Signature),
	)
	if err != nil {
		return nil, err
	}
	return &userpb.UpdateUserProfileResponse{}, nil
}

func (h *userHandler) Delete(ctx context.Context, req *userpb.DeleteUserProfileRequest) (*userpb.DeleteUserProfileResponse, error) {
	userId := middleware.OptionalUint(uint(req.UserId))
	id := middleware.OptionalUint(uint(req.Id))

	if err := h.service.DeleteProfile(userId, id); err != nil {
		return nil, err
	}
	return &userpb.DeleteUserProfileResponse{}, nil
}
