package campaign

import (
	"crowdfunding/internal/middleware"
	"crowdfunding/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(router *gin.RouterGroup, db *gorm.DB) {

	// user
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

	// router.POST("email_checkers", handler.checkEmailAvailability)

	protectedRouter := router.Group("campaigns")
	protectedRouter.Use(middleware.AuthMiddleware()) // Terapkan middleware ke seluruh grup
	{
		protectedRouter.POST("", handler.CreateCampaign)
	}

}
