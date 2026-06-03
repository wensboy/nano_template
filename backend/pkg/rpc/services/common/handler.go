package common

import (
	"context"

	"example.com/nano_template/pkg/config"
	commonpb "example.com/nano_template/proto/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MountCommonServer(s *grpc.Server, cfg *config.Config) {
	commonpb.RegisterCommonServiceServer(s, &commonHandler{cfg: cfg})
}

type commonHandler struct {
	cfg *config.Config
	commonpb.UnimplementedCommonServiceServer
}

func (h *commonHandler) Ping(ctx context.Context, req *commonpb.PingRequest) (*commonpb.PingResponse, error) {
	if req.Message != "ping" {
		return &commonpb.PingResponse{Message: "bong"}, nil
	}
	return &commonpb.PingResponse{Message: "pong"}, nil
}

func (h *commonHandler) Inspect(ctx context.Context, req *commonpb.InspectRequest) (*commonpb.InspectResponse, error) {
	return &commonpb.InspectResponse{
		Version:     "v0.1.0",
		Author:      "just you",
		Description: "A simple nano template service for cloud native applications.(web + rpc)",
	}, nil
}

func (h *commonHandler) GetTemplate(ctx context.Context, req *commonpb.GetTemplateRequest) (*commonpb.GetTemplateResponse, error) {
	template, ok := config.GetTemplate(req.TemplateId)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "template not found: %s", req.TemplateId)
	}
	return &commonpb.GetTemplateResponse{
		FrontMatter: &commonpb.TemplateFrontMatter{
			Role: template.FrontMatter.Role,
		},
		Content: string(template.Content),
	}, nil
}
