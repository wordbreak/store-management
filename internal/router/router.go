package router

import (
	"net/http"
	"store-management/internal/controller"
	"store-management/internal/model"
	"store-management/internal/response"
	"store-management/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthRequiredHandler(handler func(ctx *controller.AuthContext)) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := c.Value("user")
		if user == nil {
			c.JSON(http.StatusUnauthorized, response.New(http.StatusUnauthorized, response.MessageUnauthorized, nil))
			return
		}

		authContext := controller.AuthContext{Context: c, User: user.(*model.User)}
		handler(&authContext)
		c.Next()
	}
}

func Init(router *gin.Engine, srv *service.Service) {
	authController := controller.NewAuthController(srv.AuthService)
	productController := controller.NewProductController(srv.StoreService)

	v1 := router.Group("/v1")
	v1.POST("/auth/register", authController.Register)
	v1.POST("/auth/login", authController.Login)
	v1.POST("/auth/logout", authController.Logout)

	v1.GET("/product/:id", AuthRequiredHandler(productController.Get))
	v1.GET("/products", AuthRequiredHandler(productController.List))
	v1.GET("/products/search", AuthRequiredHandler(productController.Search))
	v1.POST("/product", AuthRequiredHandler(productController.Create))
	v1.DELETE("/product/:id", AuthRequiredHandler(productController.Delete))
	v1.PATCH("/product/:id", AuthRequiredHandler(productController.Update))
}
