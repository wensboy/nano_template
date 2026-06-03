package user

import (
	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

func MountUserRouter(Router *gin.RouterGroup, cfg *config.Config) {
	userHandler := NewUserHandler()

	userPublicRouter := Router.Group("/user")
	{
		userPublicRouter.POST("/register", userHandler.Register)
		userPublicRouter.POST("/login", userHandler.Login)

	}

	userPrivateRouter := Router.Group("/user").Use(middleware.JWTAuth(cfg))
	{
		userPrivateRouter.GET("/logout", userHandler.Logout)
		userPrivateRouter.GET("/auth", userHandler.Auth)
		userPrivateRouter.DELETE("", userHandler.Cancel)
		userPrivateRouter.PUT("/password", userHandler.UpdatePassword)
		// 用户首选项服务
		userPrivateRouter.GET("/profiles", userHandler.ListProfiles)
		userPrivateRouter.POST("/profile", userHandler.CreateProfile)
		userPrivateRouter.GET("/profile", userHandler.GetProfile)
		userPrivateRouter.PUT("/profile", userHandler.UpdateProfile)
		userPrivateRouter.DELETE("/profile", userHandler.DeleteProfile)
	}
}
