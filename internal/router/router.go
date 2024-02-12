package router

import (
	"store-management/internal/controller"
	"store-management/internal/service"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine, srv *service.Service) {
	authController := controller.NewAuthController(srv.AuthService)

	v1 := router.Group("/v1")
	v1.POST("/auth/register", authController.Register)
	v1.POST("/auth/login", authController.Login)
}
