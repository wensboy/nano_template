package common

import (
	"example.com/nano_template/pkg/config"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MountCommonRouter(Router *gin.RouterGroup, cfg *config.Config) {
	commonHandler := NewCommonHandler()
	// 这里可以添加公共的路由，例如健康检查、版本信息、接口文档等
	Router.GET("/ping", commonHandler.Ping)
	Router.GET("/inspect", commonHandler.Inspect)
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	Router.GET("/template/:template_id", commonHandler.GetTemplate)
}
