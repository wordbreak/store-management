package repository

import "store-management/internal/datasource"

type Repository interface {
	UserRepository() UserRepository
	StoreRepository() StoreRepository
}

type repositoryImpl struct {
	writer datasource.SQL
	reader datasource.SQL

	user  UserRepository
	store StoreRepository
}

func (r *repositoryImpl) UserRepository() UserRepository {
	return r.user
}

func (r *repositoryImpl) StoreRepository() StoreRepository {
	return r.store
}

var repo Repository

func Init(writer, reader datasource.SQL) {
	if repo != nil {
		return
	}

	repo = &repositoryImpl{
		writer: writer,
		reader: reader,

		user:  NewUserRepository(writer, reader),
		store: NewStoreRepository(writer, reader),
	}
}

func Get() Repository {
	return repo
}
