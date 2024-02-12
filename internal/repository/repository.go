package repository

import "store-management/internal/datasource"

type Repository interface {
	UserRepository() UserRepository
}

type repositoryImpl struct {
	writer datasource.SQL
	reader datasource.SQL

	user UserRepository
}

func (r *repositoryImpl) UserRepository() UserRepository {
	return r.user
}

var repo Repository

func Init(writer, reader datasource.SQL) {
	if repo != nil {
		return
	}

	repo = &repositoryImpl{
		writer: writer,
		reader: reader,

		user: NewUserRepository(writer, reader),
	}
}

func Get() Repository {
	return repo
}
