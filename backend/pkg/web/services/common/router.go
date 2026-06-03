package common

import (
	"example.com/nano_template/pkg/config"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MountCommonRouter(Router *gin.RouterGroup, cfg *config.Config) {
	commonHandler := NewCommonHandler()
	Router.GET("/ping", commonHandler.Ping)
	Router.GET("/inspect", commonHandler.Inspect)
	Router.GET("/template/:id", commonHandler.GetTemplate)
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
