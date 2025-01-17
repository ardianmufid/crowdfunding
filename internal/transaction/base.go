package transaction

import (
	"crowdfunding/internal/campaign"
	"crowdfunding/internal/middleware"
	"crowdfunding/internal/payment"
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

	// payment service
	paymentService := payment.NewService()

	repo := NewRepository(db)
	svc := NewService(repo, campaignRepo, paymentService)
	handler := NewHandler(svc, svcUser)

	router.Use(middleware.AuthMiddleware()).GET("/campaigns/:id/transactions", handler.GetCampaignTransactions)
	router.Use(middleware.AuthMiddleware()).GET("/transactions", handler.GetUserTransactions)
	router.Use(middleware.AuthMiddleware()).POST("/transactions", handler.CreateTransaction)
	router.POST("/transactions/notification", handler.GetNotification)

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
