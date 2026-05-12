package aliyun

import (
	"net/http"
	"path/filepath"

	"example.com/nano_template/pkg/middleware"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	aliyunoss "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/gin-gonic/gin"
)

type AliyunHandler interface {
	Presign(*gin.Context)
}

type aliyunHandler struct {
	aliyunService AliyunService
	validMimes    map[string]bool
}

func NewAliyunHandler(aliyunService AliyunService, validMimes []string) AliyunHandler {
	handler := &aliyunHandler{aliyunService: aliyunService, validMimes: make(map[string]bool)}
	for _, v := range validMimes {
		handler.validMimes[v] = true
	}
	return handler
}

// Presign godoc
// @Summary presign aliyun oss object
// @Schemes
// @Description presign aliyun oss object
// @Tags aliyun
// @Accept json
// @Produce json
// @Param req body aliyun.PresignRequest true "presign request"
// @Success 200 {object} middleware.Response{data=aliyun.PresignResponse}
// @Router /native/aliyun/presign [post]
func (h *aliyunHandler) Presign(c *gin.Context) {
	var req PresignRequest
	if err := c.ShouldBind(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	bucket := middleware.GetBucket(c)
	bucketPrefix := middleware.GetBucketPrefix(c)
	expires := middleware.GetExpires(c)
	if !h.validMimes[req.Mime] {
		middleware.Fail(c, "invalid mime type")
		return
	}
	if *req.ObjectKey == "" {
		middleware.Fail(c, "invalid object key")
		return
	}
	key := filepath.Clean(bucketPrefix + "/" + *req.ObjectKey)
	result, err := h.aliyunService.Presign(&aliyunoss.PutObjectRequest{
		Bucket:      oss.Ptr(bucket),
		Key:         oss.Ptr(key),
		ContentType: oss.Ptr(req.Mime),
	}, expires)
	if err != nil {
		middleware.Fail(c, err.Error())
	}
	var signedHeaders []SignedHeader
	if len(result.SignedHeaders) > 0 {
		for k, v := range result.SignedHeaders {
			signedHeaders = append(signedHeaders, SignedHeader{Key: k, Value: v})
		}
	}
	middleware.Succ(c, "aliyun oss presign success", PresignResponse{
		SignedUrl:     result.URL,
		Method:        result.Method,
		Expiration:    result.Expiration,
		SignedHeaders: signedHeaders,
	})
}
