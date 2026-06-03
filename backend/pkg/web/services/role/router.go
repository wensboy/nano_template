package role

import (
	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

func MountRoleRouter(Router *gin.RouterGroup, cfg *config.Config) {
	roleHandler := NewRoleHandler()

	rolePublicRouter := Router.Group("/role")
	{
		rolePublicRouter.GET("", roleHandler.List)
	}

	rolePrivateRouter := Router.Group("/role").Use(middleware.JWTAuth(cfg))
	{
		rolePrivateRouter.POST("", roleHandler.Create)
		rolePrivateRouter.GET("/:id", roleHandler.Get)
		rolePrivateRouter.PUT("/:id", roleHandler.Update)
		rolePrivateRouter.DELETE("/:id", roleHandler.Delete)
	}
}
