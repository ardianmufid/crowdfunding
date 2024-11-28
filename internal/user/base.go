package user

import (
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
}
