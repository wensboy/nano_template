package role

import (
	"context"
	"math"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/rpc/middleware"
	"example.com/nano_template/pkg/util"
	rolepb "example.com/nano_template/proto/role"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func MountRoleServer(s *grpc.Server, cfg *config.Config) {
	roleService := NewRoleService(config.GetGDB())
	roleHandler := &roleHandler{cfg: cfg, service: roleService, model: NewRoleModel()}
	rolepb.RegisterRoleServiceServer(s, roleHandler)
	for _, role := range cfg.SystemConfig.RoleConfig.Roles {
		_, err := roleHandler.service.Create(
			role.Name,
			role.Level,
			uint(role.State),
			middleware.OptionalString(role.Description),
		)
		if err != nil {
			util.Warn("[system]", zap.String("create-role", "fail"), zap.Any("role", role), zap.String("reason", err.Error()))
		}
	}
}

type roleHandler struct {
	cfg *config.Config
	rolepb.UnimplementedRoleServiceServer
	service RoleService
	model   RoleModel
}

func (h *roleHandler) Create(ctx context.Context, req *rolepb.CreateRequest) (*rolepb.CreateResponse, error) {
	role, err := h.service.Create(
		req.Name,
		int(req.Level),
		uint(req.State),
		middleware.OptionalString(req.Description),
	)
	if err != nil {
		return nil, err
	}
	return &rolepb.CreateResponse{
		Id: uint64(role.ID),
	}, nil
}

func (h *roleHandler) Get(ctx context.Context, req *rolepb.GetRequest) (*rolepb.GetResponse, error) {
	role, err := h.service.Get(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &rolepb.GetResponse{Role: h.model.ToPbRole(role)}, nil
}

func (h *roleHandler) List(ctx context.Context, req *rolepb.ListRequest) (*rolepb.ListResponse, error) {
	page := middleware.OptionalInt(int(req.Page))
	pageSize := middleware.OptionalInt(int(req.PageSize))

	level := middleware.Ptr(int(req.Level))
	if *level == math.MaxInt32 {
		level = nil
	}

	roles, total, err := h.service.List(page, pageSize, middleware.OptionalString(req.Keywords), level, middleware.Ptr(uint(req.State)))
	if err != nil {
		return nil, err
	}
	tRoles := make([]*rolepb.Role, 0, len(roles))
	for _, role := range roles {
		tRoles = append(tRoles, h.model.ToPbRole(&role))
	}
	return &rolepb.ListResponse{Roles: tRoles, Total: total}, nil
}

func (h *roleHandler) Update(ctx context.Context, req *rolepb.UpdateRequest) (*rolepb.UpdateResponse, error) {
	err := h.service.Update(
		uint(req.Id),
		middleware.OptionalString(req.Name),
		middleware.OptionalString(req.Description),
		middleware.OptionalUint(uint(req.Level)),
		middleware.OptionalUint(uint(req.State)),
	)
	if err != nil {
		return nil, err
	}
	return &rolepb.UpdateResponse{}, nil
}

func (h *roleHandler) Delete(ctx context.Context, req *rolepb.DeleteRequest) (*rolepb.DeleteResponse, error) {
	err := h.service.Delete(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &rolepb.DeleteResponse{}, nil
}
