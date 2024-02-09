package main

import (
	"net/http"
	"store-management/internal/response"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, response.New(http.StatusOK, response.ResponseMessageOK, nil))
	})
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
