package ai

import (
	"context"

	"example.com/nano_template/pkg/config"
	chatpb "example.com/nano_template/proto/chat"
	"google.golang.org/grpc"
)

func MountAiChatServer(s *grpc.Server, cfg *config.Config) {
	aiChatService := NewAiChatService(config.GetGDB())
	chatpb.RegisterAiChatServiceServer(s, &aiChatHandler{
		cfg:     cfg,
		service: aiChatService,
	})
}

type aiChatHandler struct {
	cfg *config.Config
	chatpb.UnimplementedAiChatServiceServer
	service AiChatService
}

func (h *aiChatHandler) CreateSession(ctx context.Context, req *chatpb.CreateSessionRequest) (*chatpb.CreateSessionResponse, error) {
	return h.service.CreateSession(ctx, req)
}

func (h *aiChatHandler) UpdateSession(ctx context.Context, req *chatpb.UpdateSessionRequest) (*chatpb.UpdateSessionResponse, error) {
	return h.service.UpdateSession(ctx, req)
}

func (h *aiChatHandler) GetSession(ctx context.Context, req *chatpb.GetSessionRequest) (*chatpb.GetSessionResponse, error) {
	return h.service.GetSession(ctx, req)
}

func (h *aiChatHandler) ListSessions(ctx context.Context, req *chatpb.ListSessionsRequest) (*chatpb.ListSessionsResponse, error) {
	return h.service.ListSessions(ctx, req)
}

func (h *aiChatHandler) DeleteSession(ctx context.Context, req *chatpb.DeleteSessionRequest) (*chatpb.DeleteSessnoResponse, error) {
	return h.service.DeleteSession(ctx, req)
}

func (h *aiChatHandler) Generate(req *chatpb.GenerateRequest, stream grpc.ServerStreamingServer[chatpb.GenerateResponse]) error {
	return h.service.Generate(req, stream)
}

func (h *aiChatHandler) GenerateModal(req *chatpb.GenerateModalRequest, stream grpc.ServerStreamingServer[chatpb.GenerateModalResponse]) error {
	return h.service.GenerateModal(req, stream)
}

func (h *aiChatHandler) Record(ctx context.Context, req *chatpb.RecordRequest) (*chatpb.RecordResponse, error) {
	return h.service.Record(ctx, req)
}

func (h *aiChatHandler) UpdateRecord(ctx context.Context, req *chatpb.UpdateRecordRequest) (*chatpb.UpdateRecordResponse, error) {
	return h.service.UpdateRecord(ctx, req)
}

func (h *aiChatHandler) ListRecords(ctx context.Context, req *chatpb.ListRecordsRequest) (*chatpb.ListRecordsResponse, error) {
	return h.service.ListRecords(ctx, req)
}

func (h *aiChatHandler) DeleteRecord(ctx context.Context, req *chatpb.DeleteRecordRequest) (*chatpb.DeleteRecordResponse, error) {
	return h.service.DeleteRecord(ctx, req)
}

func (h *aiChatHandler) UploadAppendixs(ctx context.Context, req *chatpb.UploadAppendixRequest) (*chatpb.UploadAppendixResponse, error) {
	return h.service.UploadAppendixs(ctx, req)
}

func (h *aiChatHandler) UpdateAppendix(ctx context.Context, req *chatpb.UpdateAppendixRequest) (*chatpb.UpdateAppendixResponse, error) {
	return h.service.UpdateAppendix(ctx, req)
}

func (h *aiChatHandler) ListAppendixs(ctx context.Context, req *chatpb.ListAppendixsRequest) (*chatpb.ListAppendixsResponse, error) {
	return h.service.ListAppendixs(ctx, req)
}

func (h *aiChatHandler) DeleteAppendix(ctx context.Context, req *chatpb.DeleteAppendixRequest) (*chatpb.DeleteAppendixResponse, error) {
	return h.service.DeleteAppendix(ctx, req)
}
