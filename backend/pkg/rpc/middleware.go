package rpc

import (
	"context"

	"example.com/nano_template/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func loggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	util.Info("[rpc request]", zap.String("method", info.FullMethod), zap.Any("req", req))
	resp, err := handler(ctx, req)
	if err != nil {
		util.Error("[rpc error]", zap.String("method", info.FullMethod), zap.Error(err))
		return nil, err
	}
	util.Info("[rpc response]", zap.String("method", info.FullMethod), zap.Any("resp", resp))
	return resp, err
}
