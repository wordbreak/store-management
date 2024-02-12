package mock

import (
	"store-management/internal/model"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindUser(phoneNumber string) (*model.User, error) {
	args := m.Called(phoneNumber)
	user := args.Get(0)
	err := args.Error(1)
	if user == nil {
		return nil, err
	}
	return user.(*model.User), err
}

func (m *UserRepositoryMock) CreateUser(phoneNumber, password string) error {
	args := m.Called(phoneNumber, password)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindUserByID(id int64) (*model.User, error) {
	args := m.Called(id)
	user := args.Get(0)
	err := args.Error(1)
	if user == nil {
		return nil, err
	}
	return user.(*model.User), err
}
