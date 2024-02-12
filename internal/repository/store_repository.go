package repository

import (
	"database/sql"
	"errors"
	"store-management/internal/datasource"
	"store-management/internal/model"
)

type StoreRepository interface {
	CreateStore(userId int64) error
	FindStoreByUserID(userId int64) (*model.Store, error)
}

type storeRepositoryImpl struct {
	writer datasource.SQL
	reader datasource.SQL
}

func NewStoreRepository(writer, reader datasource.SQL) StoreRepository {
	return &storeRepositoryImpl{
		writer: writer,
		reader: reader,
	}
}

func (s *storeRepositoryImpl) CreateStore(userId int64) error {
	_, err := s.writer.Exec("INSERT INTO store (user_id) VALUES (?)", userId)
	if err != nil {
		panic(err)
	}
	return nil
}

func (s *storeRepositoryImpl) FindStoreByUserID(userId int64) (*model.Store, error) {
	var store model.Store
	err := s.reader.Get(&store, "SELECT id, user_id FROM store WHERE user_id = ?", userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datasource.ErrNoRows
		}
		panic(err)
	}
	return &store, nil
}
