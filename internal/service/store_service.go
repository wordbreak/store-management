package service

import (
	"errors"
	"store-management/internal/model"
	"store-management/internal/repository"
)

var (
	ErrStoreNotFound   = errors.New("store not found")
	ErrProductNotFound = errors.New("product not found")
)

type StoreService interface {
	GetStoreByUserID(userId int64) (*model.Store, error)
	GetProduct(storeId, productId int64) (*model.Product, error)
	GetProductsWithPagination(storeId, page, size int64) ([]*model.Product, error)
	CreateProduct(storeId int64, product *model.Product) (int64, error)
	DeleteProduct(storeId, productId int64) error
	UpdateProduct(storeId int64, product *model.Product) error
}

type storeServiceImpl struct {
	repo struct {
		store   repository.StoreRepository
		product repository.ProductRepository
	}
}

func (s *storeServiceImpl) GetStoreByUserID(userId int64) (*model.Store, error) {
	if store, err := s.repo.store.FindStoreByUserID(userId); err != nil {
		return nil, ErrStoreNotFound
	} else {
		return store, nil
	}
}

func (s *storeServiceImpl) GetProduct(storeId, productId int64) (*model.Product, error) {
	product, err := s.repo.product.FindProduct(storeId, productId)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return product, nil

}

func (s *storeServiceImpl) CreateProduct(storeId int64, product *model.Product) (int64, error) {
	return s.repo.product.CreateProduct(storeId, product)
}

func (s *storeServiceImpl) DeleteProduct(storeId, productId int64) error {
	if err := s.repo.product.DeleteProduct(storeId, productId); err != nil {
		return ErrProductNotFound
	}
	return nil
}

func (s *storeServiceImpl) UpdateProduct(userId int64, product *model.Product) error {
	store, err := s.repo.store.FindStoreByUserID(userId)
	if err != nil {
		return ErrProductNotFound
	}
	return s.repo.product.UpdateProduct(store.ID, product)
}

func (s *storeServiceImpl) GetProductsWithPagination(storeId int64, cursor int64, limit int64) ([]*model.Product, error) {
	return s.repo.product.FindProductsWithPagination(storeId, cursor, limit)
}

func NewStoreService(repo repository.Repository) StoreService {
	service := &storeServiceImpl{}
	service.repo.product = repo.ProductRepository()
	service.repo.store = repo.StoreRepository()
	return service
}
