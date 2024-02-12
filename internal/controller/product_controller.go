package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"store-management/internal/model"
	"store-management/internal/response"
	"store-management/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductController interface {
	Get(ctx *AuthContext)
	Create(ctx *AuthContext)
	Delete(ctx *AuthContext)
	Update(ctx *AuthContext)
}

type productController struct {
	storeService service.StoreService
}

func NewProductController(storeService service.StoreService) ProductController {
	return &productController{
		storeService: storeService,
	}
}

type createInput struct {
	Category    string  `json:"category" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Cost        float64 `json:"cost" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Barcode     string  `json:"barcode" binding:"required"`
	ExpiryDate  int64   `json:"expiry_date" binding:"required"`
	Size        string  `json:"size" binding:"required"`
}

type getUriInput struct {
	ID int64 `uri:"id" binding:"required"`
}

func (c *productController) Get(ctx *AuthContext) {
	var input getUriInput

	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.New(http.StatusBadRequest, response.MessageInvalidInput, nil))
		return
	}

	store, err := c.storeService.GetStoreByUserID(ctx.User.ID)
	if err != nil {
		if errors.Is(err, service.ErrStoreNotFound) {
			ctx.JSON(http.StatusNotFound, response.New(http.StatusNotFound, response.MessageNotFound, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return
	}

	product, err := c.storeService.GetProduct(store.ID, input.ID)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) || errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, response.New(http.StatusNotFound, response.MessageNotFound, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return
	}

	ctx.JSON(http.StatusOK, response.New(http.StatusOK, response.MessageOK, product))

}

func (c *productController) Create(ctx *AuthContext) {
	var input createInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.New(http.StatusBadRequest, response.MessageInvalidInput, nil))
		return
	}

	store, err := c.storeService.GetStoreByUserID(ctx.User.ID)
	if err != nil {
		if errors.Is(err, service.ErrStoreNotFound) {
			ctx.JSON(http.StatusNotFound, response.New(http.StatusNotFound, err.Error(), nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return
	}

	product := &model.Product{
		Category:    input.Category,
		Name:        input.Name,
		Price:       input.Price,
		Cost:        input.Cost,
		Description: input.Description,
		Barcode:     input.Barcode,
		ExpiryDate:  time.Unix(input.ExpiryDate, 0),
		Size:        input.Size,
	}

	if productId, err := c.storeService.CreateProduct(store.ID, product); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
	} else {
		ctx.JSON(http.StatusCreated, response.New(http.StatusCreated, response.MessageOK, gin.H{
			"id": productId,
		}))
	}
}

type deleteInput struct {
	ID int64 `uri:"id" binding:"required"`
}

func (c *productController) Delete(ctx *AuthContext) {
	var input deleteInput
	if err := ctx.ShouldBindUri(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.New(http.StatusBadRequest, response.MessageInvalidInput, nil))
		return
	}

	store, err := c.storeService.GetStoreByUserID(ctx.User.ID)
	if err != nil {
		if errors.Is(err, service.ErrStoreNotFound) {
			ctx.JSON(http.StatusNotFound, response.New(http.StatusNotFound, response.MessageNotFound, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return

	}

	if err := c.storeService.DeleteProduct(store.ID, input.ID); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, response.New(http.StatusNotFound, response.MessageNotFound, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return
	}

	ctx.JSON(http.StatusOK, response.New(http.StatusOK, response.MessageOK, nil))
}

type updateUriInput struct {
	ID int64 `uri:"id" binding:"required"`
}
type updateInput struct {
	Category    string  `json:"category"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Cost        float64 `json:"cost"`
	Description string  `json:"description"`
	Barcode     string  `json:"barcode"`
	ExpiryDate  int64   `json:"expiry_date"`
	Size        string  `json:"size"`
}

func (c *productController) Update(ctx *AuthContext) {
	var uriInput updateUriInput
	var input updateInput

	if err := ctx.ShouldBindUri(&uriInput); err != nil {
		ctx.JSON(http.StatusBadRequest, response.New(http.StatusBadRequest, response.MessageInvalidInput, nil))
		return
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.New(http.StatusBadRequest, response.MessageInvalidInput, nil))
		return
	}

	store, err := c.storeService.GetStoreByUserID(ctx.User.ID)
	if err != nil {
		if errors.Is(err, service.ErrStoreNotFound) {
			ctx.JSON(http.StatusNotFound, response.New(http.StatusNotFound, response.MessageNotFound, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return

	}

	product := &model.Product{
		ID:          uriInput.ID,
		Category:    input.Category,
		Name:        input.Name,
		Price:       input.Price,
		Cost:        input.Cost,
		Description: input.Description,
		Barcode:     input.Barcode,
		ExpiryDate:  time.Unix(input.ExpiryDate, 0),
		Size:        input.Size,
	}

	if err := c.storeService.UpdateProduct(store.ID, product); err != nil {
		if errors.Is(err, service.ErrProductNotFound) || errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, response.New(http.StatusNotFound, response.MessageNotFound, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.New(http.StatusInternalServerError, response.MessageInternalError, nil))
		return
	}

	ctx.JSON(http.StatusOK, response.New(http.StatusOK, response.MessageOK, nil))
}
