package user

import (
	"crowdfunding/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(router *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("auth")
	{
		authRouter.POST("signup", handler.register)
		authRouter.POST("signin", handler.login)
	}

	router.POST("email_checkers", handler.checkEmailAvailability)

	protectedRouter := router.Group("")
	protectedRouter.Use(middleware.AuthMiddleware()) // Terapkan middleware ke seluruh grup
	{
		protectedRouter.POST("avatars", handler.UploadAvatar)
	}

	// router.POST("avatars", handler.UploadAvatar)
	// router.GET("user/fetch", handler.FetchUser)

}
