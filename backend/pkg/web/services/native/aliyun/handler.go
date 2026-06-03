package aliyun

import (
	"net/http"

	"example.com/nano_template/pkg/rpc"
	"example.com/nano_template/pkg/web/middleware"
	osspb "example.com/nano_template/proto/oss"
	"github.com/gin-gonic/gin"
)

type AliyunHandler interface {
	PresignUpload(*gin.Context)
	PresignDownload(c *gin.Context)
	PresignList(c *gin.Context)
}

type aliyunHandler struct {
	client osspb.AliyunOssServiceClient
	model  AliyunModel
}

type AliyunHandlerOption func(*aliyunHandler)

func NewAliyunHandler(opts ...AliyunHandlerOption) AliyunHandler {
	handler := &aliyunHandler{
		client: rpc.GetAliyunOssServiceClient(),
		model:  NewAliyunModel(),
	}
	for _, opt := range opts {
		opt(handler)
	}
	return handler
}

// PresignUpload godoc
// @Summary presign upload aliyun oss object
// @Schemes
// @Description presign upload aliyun oss object
// @Tags aliyun
// @Accept json
// @Produce json
// @Param req body aliyun.PresignUploadRequest true "presign upload request"
// @Success 200 {object} middleware.Response{data=PresignResponse}
// @Router /native/aliyun/presign/upload [post]
func (h *aliyunHandler) PresignUpload(c *gin.Context) {
	var req PresignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}
	if req.ObjectKey == nil || *req.ObjectKey == "" {
		middleware.Fail(c, "invalid object key")
		return
	}
	resp, err := h.client.PresignPutObject(c.Request.Context(), &osspb.PresignPutObjectRequest{
		ObjectKey: *req.ObjectKey,
		Mime:      req.Mime,
		Size:      int64(req.Size),
		Sender:    req.Sender,
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "aliyun oss presign upload success", h.model.ToPresignResponse(resp.PresignResponse))
}

// PresignDownload godoc
// @Summary presign download aliyun oss object
// @Schemes
// @Description presign download aliyun oss object
// @Tags aliyun
// @Accept json
// @Produce json
// @Param req body aliyun.PresignDownloadRequest true "presign download request"
// @Success 200 {object} middleware.Response{data=PresignResponse}
// @Router /native/aliyun/presign/download [post]
func (h *aliyunHandler) PresignDownload(c *gin.Context) {
	var req PresignDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}
	resp, err := h.client.PresignGetObject(c.Request.Context(), &osspb.PresignGetObjectRequest{
		BucketPrefix: req.BucketPrefix,
		ObjectKey:    req.ObjectKey,
		Getter:       req.Getter,
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "aliyun oss presign download success", h.model.ToPresignResponse(resp.PresignResponse))
}

// PresignList godoc
// @Summary presign list aliyun oss objects
// @Schemes
// @Description presign list aliyun oss objects
// @Tags aliyun
// @Accept json
// @Produce json
// @Param bucket_prefix query string false "bucket prefix"
// @Param marker query string false "marker"
// @Param page_size query int false "page size"
// @Param req body aliyun.PresignListRequest true "presign list request"
// @Success 200 {object} middleware.Response{data=ListObjectsResponse}
// @Router /native/aliyun/presign/list [post]
func (h *aliyunHandler) PresignList(c *gin.Context) {
	var req PresignListRequest
	if err := c.ShouldBind(&req); err != nil {
		middleware.Erro(c, http.StatusBadRequest, middleware.ErrInvalidQueryBody.Error())
		return
	}
	bucketPrefix := c.Query("bucket_prefix")
	marker := c.Query("marker")
	_, pageSize := middleware.ParsePageParams(c)
	resp, err := h.client.ListObjects(c.Request.Context(), &osspb.ListObjectsRequest{
		BucketPrefix: bucketPrefix,
		Marker:       marker,
		MaxKeys:      int32(pageSize),
		Getter:       req.Getter,
	})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "aliyun list objects success", h.model.ToListObjectsResponse(resp))
}
