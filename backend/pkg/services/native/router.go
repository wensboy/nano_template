package native

import (
	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/services/native/aliyun"
	"github.com/gin-gonic/gin"
)

func MountNativeRouter(Router *gin.RouterGroup, cfg *config.Config) {
	nativeRouter := Router.Group("/native")
	{
		if cfg.AliyunOssConfig.Enable {
			aliyun.MountAliyunRouter(nativeRouter, cfg)
		}
	}
}
