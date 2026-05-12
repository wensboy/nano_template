package aliyun

import (
	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func MountAliyunRouter(Router *gin.RouterGroup, cfg *config.Config) {
	aliyunService := NewAliyunService(config.GetGDB())
	aliyunHandler := NewAliyunHandler(aliyunService, cfg.AliyunOssConfig.ValidMimes)

	aliyunPrivateRouter := Router.Group("/aliyun").Use(middleware.AliyunOssHandler(cfg.AliyunOssConfig))
	{
		aliyunPrivateRouter.POST("/presign", aliyunHandler.Presign)
	}
}
