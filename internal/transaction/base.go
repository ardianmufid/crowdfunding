package transaction

import (
	"crowdfunding/internal/campaign"
	"crowdfunding/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Init(router *gin.RouterGroup, db *sqlx.DB) {

	// campaign repository
	campaignRepo := campaign.NewRepository(db)

	// user repository
	repoUser := user.NewRepository(db)
	svcUser := user.NewService(repoUser)

	repo := NewRepository(db)
	svc := NewService(repo, campaignRepo)
	handler := NewHandler(svc, svcUser)

	router.GET("/campaigns/:id/transactions", handler.GetCampaignTransactions)

	// campaignRouter := router.Group("campaigns")
	// {
	// 	campaignRouter.GET("", handler.GetCampaigns)
	// 	campaignRouter.GET("/:id", handler.GetCampaign)
	// }

	// protectedRouter := router.Group("campaigns")
	// protectedRouter.Use(middleware.AuthMiddleware())
	// {
	// 	protectedRouter.POST("", handler.CreateCampaign)
	// 	protectedRouter.PUT("/:id", handler.UpdateCampaign)
	// }
	// router.Use(middleware.AuthMiddleware()).POST("/campaign-images", handler.UploadImage)

}
