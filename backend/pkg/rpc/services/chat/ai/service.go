package ai

import (
	"context"

	chatpb "example.com/nano_template/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AiChatService interface {
	CreateSession(ctx context.Context, req *chatpb.CreateSessionRequest) (*chatpb.CreateSessionResponse, error)
	UpdateSession(ctx context.Context, req *chatpb.UpdateSessionRequest) (*chatpb.UpdateSessionResponse, error)
	GetSession(ctx context.Context, req *chatpb.GetSessionRequest) (*chatpb.GetSessionResponse, error)
	ListSessions(ctx context.Context, req *chatpb.ListSessionsRequest) (*chatpb.ListSessionsResponse, error)
	DeleteSession(ctx context.Context, req *chatpb.DeleteSessionRequest) (*chatpb.DeleteSessnoResponse, error)
	Generate(req *chatpb.GenerateRequest, stream grpc.ServerStreamingServer[chatpb.GenerateResponse]) error
	GenerateModal(req *chatpb.GenerateModalRequest, stream grpc.ServerStreamingServer[chatpb.GenerateModalResponse]) error
	Record(ctx context.Context, req *chatpb.RecordRequest) (*chatpb.RecordResponse, error)
	UpdateRecord(ctx context.Context, req *chatpb.UpdateRecordRequest) (*chatpb.UpdateRecordResponse, error)
	ListRecords(ctx context.Context, req *chatpb.ListRecordsRequest) (*chatpb.ListRecordsResponse, error)
	DeleteRecord(ctx context.Context, req *chatpb.DeleteRecordRequest) (*chatpb.DeleteRecordResponse, error)
	UploadAppendixs(ctx context.Context, req *chatpb.UploadAppendixRequest) (*chatpb.UploadAppendixResponse, error)
	UpdateAppendix(ctx context.Context, req *chatpb.UpdateAppendixRequest) (*chatpb.UpdateAppendixResponse, error)
	ListAppendixs(ctx context.Context, req *chatpb.ListAppendixsRequest) (*chatpb.ListAppendixsResponse, error)
	DeleteAppendix(ctx context.Context, req *chatpb.DeleteAppendixRequest) (*chatpb.DeleteAppendixResponse, error)
}

type aiChatService struct {
	db *gorm.DB
}

func NewAiChatService(db *gorm.DB) AiChatService {
	return &aiChatService{
		db: db,
	}
}

func (s *aiChatService) CreateSession(ctx context.Context, req *chatpb.CreateSessionRequest) (*chatpb.CreateSessionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateSession not implemented")
}

func (s *aiChatService) UpdateSession(ctx context.Context, req *chatpb.UpdateSessionRequest) (*chatpb.UpdateSessionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateSession not implemented")
}

func (s *aiChatService) GetSession(ctx context.Context, req *chatpb.GetSessionRequest) (*chatpb.GetSessionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetSession not implemented")
}

func (s *aiChatService) ListSessions(ctx context.Context, req *chatpb.ListSessionsRequest) (*chatpb.ListSessionsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListSessions not implemented")
}

func (s *aiChatService) DeleteSession(ctx context.Context, req *chatpb.DeleteSessionRequest) (*chatpb.DeleteSessnoResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteSession not implemented")
}

func (s *aiChatService) Generate(req *chatpb.GenerateRequest, stream grpc.ServerStreamingServer[chatpb.GenerateResponse]) error {
	return stream.Send(&chatpb.GenerateResponse{})
}

func (s *aiChatService) GenerateModal(req *chatpb.GenerateModalRequest, stream grpc.ServerStreamingServer[chatpb.GenerateModalResponse]) error {
	return stream.Send(&chatpb.GenerateModalResponse{})
}

func (s *aiChatService) Record(ctx context.Context, req *chatpb.RecordRequest) (*chatpb.RecordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method Record not implemented")
}

func (s *aiChatService) UpdateRecord(ctx context.Context, req *chatpb.UpdateRecordRequest) (*chatpb.UpdateRecordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateRecord not implemented")
}

func (s *aiChatService) ListRecords(ctx context.Context, req *chatpb.ListRecordsRequest) (*chatpb.ListRecordsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListRecords not implemented")
}

func (s *aiChatService) DeleteRecord(ctx context.Context, req *chatpb.DeleteRecordRequest) (*chatpb.DeleteRecordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteRecord not implemented")
}

func (s *aiChatService) UploadAppendixs(ctx context.Context, req *chatpb.UploadAppendixRequest) (*chatpb.UploadAppendixResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UploadAppendixs not implemented")
}

func (s *aiChatService) UpdateAppendix(ctx context.Context, req *chatpb.UpdateAppendixRequest) (*chatpb.UpdateAppendixResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateAppendix not implemented")
}

func (s *aiChatService) ListAppendixs(ctx context.Context, req *chatpb.ListAppendixsRequest) (*chatpb.ListAppendixsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListAppendixs not implemented")
}

func (s *aiChatService) DeleteAppendix(ctx context.Context, req *chatpb.DeleteAppendixRequest) (*chatpb.DeleteAppendixResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteAppendix not implemented")
}
