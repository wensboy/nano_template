package common

import (
	"net/http"
	"strings"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type CommonHandler interface {
	Ping(c *gin.Context)
	Inspect(c *gin.Context)
	GetTemplate(c *gin.Context)
}

type commonHandler struct {
}

func NewCommonHandler() CommonHandler {
	return &commonHandler{}
}

// Ping godoc
// @Summary ping health check
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=nil}
// @Router /ping [get]
func (h *commonHandler) Ping(c *gin.Context) {
	middleware.Succ(c, "pong", nil)
}

// Inspect godoc
// @Summary inspect server information
// @Schemes
// @Description inspect server information
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=common.InspectResponse}
// @Router /inspect [get]
func (h *commonHandler) Inspect(c *gin.Context) {
	middleware.Succ(c, "server is running", InspectResponse{
		Version:     "v0.1.0",
		Author:      "wendisx",
		Description: "",
	})
}

// GetTemplate godoc
// @Summary get template by id
// @Schemes
// @Description get template by id
// @Tags example
// @Accept json
// @Produce json
// @Param template_id path string true "Template ID"
// @Success 200 {object} middleware.Response{data=common.TemplateResponse}
// @Router /template/{template_id} [get]
func (h *commonHandler) GetTemplate(c *gin.Context) {
	tempId := c.Param("template_id")
	if tempId == "" {
		middleware.Erro(c, http.StatusBadRequest, "invalid template id")
		return
	}
	template, ok := config.GetTemplate(tempId)
	if !ok {
		middleware.Fail(c, "template not found")
		return
	}
	middleware.Succ(c, "template found", TemplateResponse{
		Id:          tempId,
		FrontMatter: template.FrontMatter,
		Content:     strings.ReplaceAll(string(template.Content), "\n", ""),
	})
}
