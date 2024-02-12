package controller

import (
	"store-management/internal/model"

	"github.com/gin-gonic/gin"
)

type AuthContext struct {
	*gin.Context
	User *model.User
}
