package ai

import (
	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

func MountAiChatRouter(Router *gin.RouterGroup, cfg *config.Config) {
	aiChatHandler := NewAiChatHandler()

	aiChatPrivateRouter := Router.Group("/aichat").Use(middleware.JWTAuth(cfg))
	{
		aiChatPrivateRouter.POST("/sessions", aiChatHandler.CreateSession)
		aiChatPrivateRouter.GET("/sessions", aiChatHandler.ListSessions)
		aiChatPrivateRouter.GET("/sessions/:session_id", aiChatHandler.GetSession)
		aiChatPrivateRouter.PUT("/sessions/:session_id", aiChatHandler.UpdateSession)
		aiChatPrivateRouter.DELETE("/sessions/:session_id", aiChatHandler.DeleteSession)

		aiChatPrivateRouter.POST("/generate", middleware.SSEHandler(), aiChatHandler.Generate)
		aiChatPrivateRouter.POST("/generate-modal", middleware.SSEHandler(), aiChatHandler.GenerateModal)

		aiChatPrivateRouter.POST("/records", aiChatHandler.Record)
		aiChatPrivateRouter.GET("/records/:session_id", aiChatHandler.ListRecords)
		aiChatPrivateRouter.PUT("/records/:record_id", aiChatHandler.UpdateRecord)
		aiChatPrivateRouter.DELETE("/records/:record_id", aiChatHandler.DeleteRecord)

		aiChatPrivateRouter.POST("/appendixs", aiChatHandler.UploadAppendixs)
		aiChatPrivateRouter.GET("/appendixs/:message_id", aiChatHandler.ListAppendixs)
		aiChatPrivateRouter.PUT("/appendixs/:appendix_id", aiChatHandler.UpdateAppendix)
		aiChatPrivateRouter.DELETE("/appendixs/:appendix_id", aiChatHandler.DeleteAppendix)
	}
}
