package service

import "store-management/internal/repository"

var service *Service

type Service struct {
	AuthService AuthService
}

func Init(repository repository.Repository) {
	if service != nil {
		return
	}

	service = &Service{
		AuthService: NewAuthService(repository),
	}
}

func Get() *Service {
	return service
}
