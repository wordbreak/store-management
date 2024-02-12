package mock

import (
	"store-management/internal/model"

	"github.com/stretchr/testify/mock"
)

type StoreRepositoryMock struct {
	mock.Mock
}

func (m *StoreRepositoryMock) FindStoreByUserID(userId int64) (*model.Store, error) {
	args := m.Called(userId)
	store := args.Get(0)
	err := args.Error(1)
	if store == nil {
		return nil, err
	}
	return store.(*model.Store), err
}

func (m *StoreRepositoryMock) CreateStore(userId int64) error {
	args := m.Called(userId)
	return args.Error(0)
}
