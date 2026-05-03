package userSys

import (
	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func MountUserSysRouter(Router *gin.RouterGroup, cfg *config.Config) {
	userService := NewUserService(config.GDB)
	userHandler := NewUserHandler(userService)

	userPublicRouter := Router.Group("/user")
	{
		userPublicRouter.POST("/register", userHandler.Register)
		userPublicRouter.POST("/login", userHandler.Login)
	}

	userPrivateRouter := Router.Group("/user").Use(middleware.JWTAuth())
	{
		userPrivateRouter.GET("/details", userHandler.GetUserDetails)
		userPrivateRouter.PUT("/update/profile", userHandler.UpdateUserProfile)
		userPrivateRouter.PUT("/update/password", userHandler.ChangePassword)
		userPrivateRouter.DELETE("/delete", userHandler.DeactivateUser)
	}
}
