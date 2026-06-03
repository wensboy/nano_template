package aliyun

import (
	"context"
	"time"

	"example.com/nano_template/pkg/config"
	aliyunoss "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"gorm.io/gorm"
)

type AliyunOssService interface {
	PresignPutObject(ctx context.Context, req *aliyunoss.PutObjectRequest, expires int) (*aliyunoss.PresignResult, error)
	PresignGetObject(ctx context.Context, req *aliyunoss.GetObjectRequest, expires int) (*aliyunoss.PresignResult, error)
	ListObjects(ctx context.Context, req *aliyunoss.ListObjectsRequest) (*aliyunoss.ListObjectsResult, error)
}

type aliyunOssService struct {
	db *gorm.DB
}

func NewAliyunOssService(db *gorm.DB) AliyunOssService {
	return &aliyunOssService{
		db: db,
	}
}

func (s *aliyunOssService) PresignPutObject(ctx context.Context, req *aliyunoss.PutObjectRequest, expires int) (*aliyunoss.PresignResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.presign(ctx, req, expires)
}

func (s *aliyunOssService) PresignGetObject(ctx context.Context, req *aliyunoss.GetObjectRequest, expires int) (*aliyunoss.PresignResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.presign(ctx, req, expires)
}

func (s *aliyunOssService) ListObjects(ctx context.Context, req *aliyunoss.ListObjectsRequest) (*aliyunoss.ListObjectsResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return config.GetAliyunOss().ListObjects(ctx, req)
}

func (s *aliyunOssService) presign(ctx context.Context, req any, expires int) (*aliyunoss.PresignResult, error) {
	return config.GetAliyunOss().Presign(ctx, req, aliyunoss.PresignExpires(time.Duration(expires)*time.Second))
}
