package campaign

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(router *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	campaignRouter := router.Group("")
	{
		campaignRouter.GET("signup", handler.GetCampaigns)
		// campaignRouter.POST("signin", handler.login)
	}

	// router.POST("email_checkers", handler.checkEmailAvailability)

	// protectedRouter := router.Group("")
	// protectedRouter.Use(middleware.AuthMiddleware()) // Terapkan middleware ke seluruh grup
	// {
	// 	protectedRouter.POST("avatars", handler.UploadAvatar)
	// }

}
