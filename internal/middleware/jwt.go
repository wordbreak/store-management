package middleware

import (
	"fmt"
	"net/http"
	"os"
	"store-management/internal/repository"
	"store-management/internal/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware(userRepository repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")

		if err != nil {
			c.Next()
			return
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.JSON(http.StatusUnauthorized, response.New(http.StatusUnauthorized, response.MessageUnauthorized, nil))
				return
			}

			if sub, err := claims.GetSubject(); err == nil {
				if id, err := strconv.ParseInt(sub, 10, 64); err == nil {
					if user, err := userRepository.FindUserByID(id); err == nil {
						c.Set("user", user)
					}
				}
			}
		}
		c.Next()
	}
}
