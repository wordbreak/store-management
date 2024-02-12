package service

import "store-management/internal/repository"

var service *Service

type Service struct {
	AuthService  AuthService
	StoreService StoreService
}

func Init(repository repository.Repository) {
	if service != nil {
		return
	}

	service = &Service{
		AuthService:  NewAuthService(repository),
		StoreService: NewStoreService(repository),
	}
}

func Get() *Service {
	return service
}
