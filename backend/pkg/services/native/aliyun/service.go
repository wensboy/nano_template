package aliyun

import (
	"context"
	"time"

	"example.com/nano_template/pkg/config"
	aliyunoss "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"gorm.io/gorm"
)

type AliyunService interface {
	Presign(*aliyunoss.PutObjectRequest, int) (*aliyunoss.PresignResult, error)
}

type aliyunService struct {
	db *gorm.DB
}

func NewAliyunService(db *gorm.DB) AliyunService {
	return &aliyunService{db: db}
}

func (s *aliyunService) Presign(req *aliyunoss.PutObjectRequest, expires int) (*aliyunoss.PresignResult, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()
	return config.GetAliyunOss().Presign(ctx, req, aliyunoss.PresignExpires(time.Duration(expires)*time.Second))
}
