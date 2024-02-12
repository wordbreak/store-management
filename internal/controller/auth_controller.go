package controller

import (
	"errors"
	"net/http"
	"store-management/internal/response"
	"store-management/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return authController{
		authService: authService,
	}
}

type RegisterInput struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func (c authController) Register(ctx *gin.Context) {
	var input RegisterInput
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
