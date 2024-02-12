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
	CreateProduct(storeId int64, product *model.Product) (int64, error)
	DeleteProduct(storeId, productId int64) error
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

func (s *storeServiceImpl) CreateProduct(storeId int64, product *model.Product) (int64, error) {
	return s.repo.product.CreateProduct(storeId, product)
}

func (s *storeServiceImpl) DeleteProduct(userId, productId int64) error {
	store, err := s.repo.store.FindStoreByUserID(userId)
	if err != nil {
		return ErrProductNotFound
	}
	if err = s.repo.product.DeleteProduct(store.ID, productId); err != nil {
		return ErrProductNotFound
	}
	return nil
}

func NewStoreService(repo repository.Repository) StoreService {
	service := &storeServiceImpl{}
	service.repo.product = repo.ProductRepository()
	service.repo.store = repo.StoreRepository()
	return service
}
