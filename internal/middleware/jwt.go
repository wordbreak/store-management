package middleware

import "github.com/gin-gonic/gin"

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO implement
		c.Next()
	}
}
