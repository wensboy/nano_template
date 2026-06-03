package middleware

import (
	"context"
	"strings"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UnaryInterceptor struct{}

type (
	AuthedUser struct {
		UserID    uint
		Username  string
		RoleID    uint
		RoleLevel uint
	}
)

const (
	RpcUserIDKey    = "auth_rpc_user_id"
	RpcUsernameKey  = "auth_rpc_username"
	RpcUserRoleKey  = "auth_rpc_user_role"
	RpcRoleLevelKey = "auth_rpc_role_level"
)

func (i *UnaryInterceptor) LoggerInterceptor(cfg *config.RpcConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !cfg.MiddlewareConfig.LoggerConfig.Enable {
			return handler(ctx, req)
		}
		util.Info("[rpc request]", zap.String("method", info.FullMethod), zap.Any("req", req))
		resp, err := handler(ctx, req)
		if err != nil {
			util.Error("[rpc error]", zap.String("method", info.FullMethod), zap.Error(err))
			return nil, err
		}
		util.Info("[rpc response]", zap.String("method", info.FullMethod), zap.Any("resp", resp))
		return resp, err
	}
}

func (i *UnaryInterceptor) AuthInterceptor(cfg *config.RpcConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx, err = extractAuth(ctx)
		if err != nil {
			if !cfg.MiddlewareConfig.AuthConfig.Enable || authPass(cfg.MiddlewareConfig.AuthConfig.PassFilter, info.FullMethod) {
				util.Info("[req passed]", zap.String("method", info.FullMethod))
				return handler(ctx, req)
			}
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (i *UnaryInterceptor) RecoveryInterceptor(cfg *config.RpcConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if !cfg.MiddlewareConfig.RecoveryConfig.Enable {
			return handler(ctx, req)
		}

		defer func() {
			if r := recover(); r != nil {
				util.Error("[PANIC]", zap.String("method", info.FullMethod), zap.Any("panic", r))
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

func GetAuthedUser(ctx context.Context) (AuthedUser, error) {
	if ctx == nil {
		return AuthedUser{}, status.Errorf(codes.Unauthenticated, "missing authed context")
	}

	userIDStr, ok := ctx.Value(RpcUserIDKey).(string)
	if !ok {
		return AuthedUser{}, status.Errorf(codes.Unauthenticated, "missing authed user info")
	}

	userName, ok := ctx.Value(RpcUsernameKey).(string)
	if !ok {
		return AuthedUser{}, status.Errorf(codes.Unauthenticated, "missing authed user info")
	}

	userRoleStr, ok := ctx.Value(RpcUserRoleKey).(string)
	if !ok {
		return AuthedUser{}, status.Errorf(codes.Unauthenticated, "missing authed user info")
	}

	roleLevelStr, ok := ctx.Value(RpcRoleLevelKey).(string)
	if !ok {
		return AuthedUser{}, status.Errorf(codes.Unauthenticated, "missing authed user info")
	}

	userID, ok := util.StrToUint(userIDStr)
	if !ok {
		return AuthedUser{}, status.Errorf(codes.Unauthenticated, "invalid user ID")
	}

	userRole, ok := util.StrToUint(userRoleStr)
	if !ok {
		return AuthedUser{}, status.Errorf(codes.Unauthenticated, "invalid user role")
	}

	roleLevel, ok := util.StrToUint(roleLevelStr)
	if !ok {
		return AuthedUser{}, status.Errorf(codes.Unauthenticated, "invalid role level")
	}

	return AuthedUser{
		UserID:    userID,
		Username:  userName,
		RoleID:    userRole,
		RoleLevel: roleLevel,
	}, nil
}

func extractAuth(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	userID := md.Get(RpcUserIDKey)
	userName := md.Get(RpcUsernameKey)
	userRole := md.Get(RpcUserRoleKey)
	roleLevel := md.Get(RpcRoleLevelKey)

	if len(userID) == 0 || len(userName) == 0 || len(userRole) == 0 || len(roleLevel) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user info")
	}
	util.Info("[authed user] - ", zap.String("user_id", userID[0]), zap.String("user_name", userName[0]), zap.String("user_role", userRole[0]), zap.String("role_level", roleLevel[0]))
	ctx = context.WithValue(ctx, RpcUserIDKey, userID[0])
	ctx = context.WithValue(ctx, RpcUsernameKey, userName[0])
	ctx = context.WithValue(ctx, RpcUserRoleKey, userRole[0])
	ctx = context.WithValue(ctx, RpcRoleLevelKey, roleLevel[0])
	return ctx, nil
}

func authPass(whiteList []string, target string) bool {
	for _, prefix := range whiteList {
		if strings.HasPrefix(target, prefix) {
			return true
		}
	}
	return false
}
