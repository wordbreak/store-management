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

func NewMockedRepository(userRepoMock *UserRepositoryMock) *MockedRepository {
	mockRepo := new(MockedRepository)
	mockRepo.On("UserRepository").Return(userRepoMock)
	return mockRepo
}
