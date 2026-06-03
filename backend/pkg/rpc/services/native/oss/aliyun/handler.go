package aliyun

import (
	"context"
	"path"
	"slices"
	"strings"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/util"
	osspb "example.com/nano_template/proto/oss"
	aliyunoss "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MountAliyunOssServer(s *grpc.Server, cfg *config.Config) {
	if !cfg.AliyunOssConfig.Enable {
		util.Warn("[modules]", zap.String("mount-service", "skipped"), zap.String("service", "aliyun oss service"))
		return
	}
	aliyunOssService := NewAliyunOssService(config.GetGDB())
	aliyunOssHandler := &aliyunOssHandler{
		cfg:     cfg,
		service: aliyunOssService,
		model:   NewAliyunOssModel(),
	}
	osspb.RegisterAliyunOssServiceServer(s, aliyunOssHandler)
}

type aliyunOssHandler struct {
	cfg *config.Config
	osspb.UnimplementedAliyunOssServiceServer
	service AliyunOssService
	model   AliyunOssModel
}

func (h *aliyunOssHandler) PresignPutObject(ctx context.Context, req *osspb.PresignPutObjectRequest) (*osspb.PresignPutObjectResponse, error) {
	key, err := h.objectKey(h.cfg.AliyunOssConfig.BucketPrefix, req.ObjectKey)
	if err != nil {
		return nil, err
	}
	if h.cfg.AliyunOssConfig.MaxSize > 0 && req.Size > int64(h.cfg.AliyunOssConfig.MaxSize) {
		return nil, status.Error(codes.InvalidArgument, "object size exceeds limit")
	}
	if !h.validMime(req.Mime) {
		return nil, status.Error(codes.InvalidArgument, "invalid mime type")
	}

	result, err := h.service.PresignPutObject(ctx, &aliyunoss.PutObjectRequest{
		Bucket:      aliyunoss.Ptr(h.cfg.AliyunOssConfig.Bucket),
		Key:         aliyunoss.Ptr(key),
		ContentType: optionalStringPtr(req.Mime),
	}, h.cfg.AliyunOssConfig.Expires)
	if err != nil {
		return nil, err
	}

	return &osspb.PresignPutObjectResponse{PresignResponse: h.model.ToPresignResponse(result)}, nil
}

func (h *aliyunOssHandler) PresignGetObject(ctx context.Context, req *osspb.PresignGetObjectRequest) (*osspb.PresignGetObjectResponse, error) {
	key, err := h.objectKey(req.BucketPrefix, req.ObjectKey)
	if err != nil {
		return nil, err
	}

	result, err := h.service.PresignGetObject(ctx, &aliyunoss.GetObjectRequest{
		Bucket: aliyunoss.Ptr(h.cfg.AliyunOssConfig.Bucket),
		Key:    aliyunoss.Ptr(key),
	}, h.cfg.AliyunOssConfig.Expires)
	if err != nil {
		return nil, err
	}

	return &osspb.PresignGetObjectResponse{PresignResponse: h.model.ToPresignResponse(result)}, nil
}

func (h *aliyunOssHandler) ListObjects(ctx context.Context, req *osspb.ListObjectsRequest) (*osspb.ListObjectsResponse, error) {
	prefix, err := h.objectPrefix(req.BucketPrefix)
	if err != nil {
		return nil, err
	}

	result, err := h.service.ListObjects(ctx, &aliyunoss.ListObjectsRequest{
		Bucket:  aliyunoss.Ptr(h.cfg.AliyunOssConfig.Bucket),
		Prefix:  optionalStringPtr(prefix),
		Marker:  optionalStringPtr(req.Marker),
		MaxKeys: req.MaxKeys,
	})
	if err != nil {
		return nil, err
	}

	return h.model.ToListObjectsResponse(result), nil
}

func (h *aliyunOssHandler) objectKey(bucketPrefix, objectKey string) (string, error) {
	key, err := cleanOssPath(objectKey)
	if err != nil || key == "" {
		return "", status.Error(codes.InvalidArgument, "invalid object key")
	}

	prefix, err := h.objectPrefix(bucketPrefix)
	if err != nil {
		return "", err
	}
	if prefix == "" {
		return key, nil
	}
	return path.Join(prefix, key), nil
}

func (h *aliyunOssHandler) objectPrefix(bucketPrefix string) (string, error) {
	rootPrefix, err := cleanOssPath(h.cfg.AliyunOssConfig.BucketPrefix)
	if err != nil {
		return "", err
	}
	requestPrefix, err := cleanOssPath(bucketPrefix)
	if err != nil {
		return "", err
	}
	if rootPrefix == "" {
		return requestPrefix, nil
	}
	if requestPrefix == "" {
		return rootPrefix, nil
	}
	return path.Join(rootPrefix, requestPrefix), nil
}

func (h *aliyunOssHandler) validMime(mime string) bool {
	if len(h.cfg.AliyunOssConfig.ValidMimes) == 0 || mime == "" {
		return true
	}
	return slices.Contains(h.cfg.AliyunOssConfig.ValidMimes, mime)
}

func cleanOssPath(value string) (string, error) {
	value = strings.Trim(strings.TrimSpace(value), "/")
	if value == "" {
		return "", nil
	}

	cleaned := path.Clean(value)
	if cleaned == "." {
		return "", nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", status.Error(codes.InvalidArgument, "invalid object path")
	}
	return cleaned, nil
}

func optionalStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return aliyunoss.Ptr(value)
}
