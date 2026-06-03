package common

import (
	"net/http"

	"example.com/nano_template/pkg/rpc"
	"example.com/nano_template/pkg/web/middleware"
	commonpb "example.com/nano_template/proto/common"
	"github.com/gin-gonic/gin"
)

type CommonHandler interface {
	Ping(c *gin.Context)
	Inspect(c *gin.Context)
	GetTemplate(c *gin.Context)
}

type commonHandler struct {
	client commonpb.CommonServiceClient
	model  CommonModel
}

type commonHandlerOption func(*commonHandler)

func NewCommonHandler(opts ...commonHandlerOption) CommonHandler {
	h := &commonHandler{
		client: rpc.GetCommonServiceClient(),
		model:  NewCommonModel(),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Ping godoc
// @Summary ping health check
// @Schemes
// @Description do ping
// @Tags common
// @Produce json
// @Success 200 {object} middleware.Response{data=middleware.EmptyData}
// @Router /ping [get]
func (h *commonHandler) Ping(c *gin.Context) {
	middleware.SuccNoMore(c, "pong")
}

// Inspect godoc
// @Summary inspect server information
// @Schemes
// @Description inspect server information
// @Tags common
// @Produce json
// @Success 200 {object} middleware.Response{data=InspectResponse}
// @Router /inspect [get]
func (h *commonHandler) Inspect(c *gin.Context) {
	resp, err := h.client.Inspect(c.Request.Context(), &commonpb.InspectRequest{})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "inspect success", h.model.ToInspectResponse(resp))
}

// GetTemplate godoc
// @Summary get template by id
// @Schemes
// @Description get template by id
// @Tags common
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} middleware.Response{data=GetTemplateResponse}
// @Router /template/{id} [get]
func (h *commonHandler) GetTemplate(c *gin.Context) {
	tempId := c.Param("id")
	if tempId == "" {
		middleware.Erro(c, http.StatusBadRequest, "invalid template id")
		return
	}
	resp, err := h.client.GetTemplate(c.Request.Context(), &commonpb.GetTemplateRequest{TemplateId: tempId})
	if err != nil {
		middleware.Fail(c, err.Error())
		return
	}
	middleware.Succ(c, "template found", h.model.ToGetTemplateResponse(tempId, resp))
}
