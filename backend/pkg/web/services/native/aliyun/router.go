package aliyun

import (
	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

func MountAliyunRouter(Router *gin.RouterGroup, cfg *config.Config) {
	aliyunHandler := NewAliyunHandler()

	aliyunPrivateRouter := Router.Group("/aliyun").Use(middleware.AliyunOssHandler(cfg.AliyunOssConfig))
	{
		aliyunPrivateRouter.POST("/presign/upload", aliyunHandler.PresignUpload)
		aliyunPrivateRouter.POST("/presign/download", aliyunHandler.PresignDownload)
		aliyunPrivateRouter.POST("/presign/list", aliyunHandler.PresignList)
	}
}
