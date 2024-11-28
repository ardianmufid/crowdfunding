package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(router *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := newHandler(svc)

	_ = handler

	authRouter := router.Group("auth")
	{
		authRouter.POST("register", handler.register)
		// authRouter.POST("login", handler.login)
	}
}
