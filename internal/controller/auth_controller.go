package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"store-management/internal/response"
	"store-management/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return authController{
		authService: authService,
	}
}

type registerInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func (c authController) Register(ctx *gin.Context) {
	var input registerInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.New(http.StatusBadRequest, response.MessageInvalidInput, nil))
		return
	}

	err := c.authService.Register(input.PhoneNumber, input.Password)

	if err != nil {
		if errors.Is(err, service.ErrDuplicateUser) {
			ctx.JSON(http.StatusConflict, response.New(http.StatusConflict, err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return
	}

	ctx.JSON(http.StatusCreated, response.New(http.StatusCreated, response.MessageOK, nil))
}

type loginInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func (c authController) Login(ctx *gin.Context) {
	var input loginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.New(http.StatusBadRequest, response.MessageInvalidInput, nil))
		return
	}

	user, err := c.authService.Login(input.PhoneNumber, input.Password)

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, response.New(http.StatusNotFound, err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("auth_token", tokenString, int(time.Hour.Seconds()), "/", "", false, true)
	ctx.JSON(http.StatusOK, response.New(http.StatusOK, response.MessageOK, nil))
}
