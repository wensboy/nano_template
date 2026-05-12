package common

import (
	"context"

	commonPb "example.com/nano_template/proto"
	"google.golang.org/grpc"
)

func MountCommonServer(s *grpc.Server) {
	commonPb.RegisterCommonServiceServer(s, &commonServer{})
}

type commonServer struct {
	commonPb.UnimplementedCommonServiceServer
}

func (s *commonServer) Ping(ctx context.Context, req *commonPb.PingRequest) (*commonPb.PingResponse, error) {
	if req.Message != "ping" {
		return &commonPb.PingResponse{Message: "eonq"}, nil
	}
	return &commonPb.PingResponse{Message: "pong"}, nil
}
