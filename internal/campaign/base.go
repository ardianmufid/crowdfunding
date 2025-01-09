package campaign

import (
	"crowdfunding/internal/middleware"
	"crowdfunding/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Init(router *gin.RouterGroup, db *sqlx.DB) {

	repoUser := user.NewRepository(db)
	svcUser := user.NewService(repoUser)

	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc, svcUser)

	campaignRouter := router.Group("campaigns")
	{
		campaignRouter.GET("", handler.GetCampaigns)
		campaignRouter.GET("/:id", handler.GetCampaign)
	}

	protectedRouter := router.Group("campaigns")
	protectedRouter.Use(middleware.AuthMiddleware())
	{
		protectedRouter.POST("", handler.CreateCampaign)
		protectedRouter.PUT("/:id", handler.UpdateCampaign)
	}
	router.Use(middleware.AuthMiddleware()).POST("/campaign-images", handler.UploadImage)

}
