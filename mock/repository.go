package mock

import (
	"store-management/internal/repository"

	"github.com/stretchr/testify/mock"
)

type MockedRepository struct {
	mock.Mock
	repository.Repository
}

func (m *MockedRepository) UserRepository() repository.UserRepository {
	args := m.Called()
	return args.Get(0).(repository.UserRepository)
}

func (m *MockedRepository) StoreRepository() repository.StoreRepository {
	args := m.Called()
	return args.Get(0).(repository.StoreRepository)
}

func NewMockedRepository(userRepoMock *UserRepositoryMock, storeRepoMock *StoreRepositoryMock) *MockedRepository {
	mockRepo := new(MockedRepository)
	mockRepo.On("UserRepository").Return(userRepoMock)
	mockRepo.On("StoreRepository").Return(storeRepoMock)
	return mockRepo
}
