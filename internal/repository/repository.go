package repository

import "store-management/internal/datasource"

type Repository interface {
	UserRepository() UserRepository
	StoreRepository() StoreRepository
	ProductRepository() ProductRepository
}

type repositoryImpl struct {
	writer      datasource.SQL
	reader      datasource.SQL
	transaction datasource.Transaction

	user    UserRepository
	store   StoreRepository
	product ProductRepository
}

func (r *repositoryImpl) UserRepository() UserRepository {
	return r.user
}

func (r *repositoryImpl) StoreRepository() StoreRepository {
	return r.store
}

func (r *repositoryImpl) ProductRepository() ProductRepository {
	return r.product
}

var repo Repository

func Init(writer, reader datasource.SQL, transaction datasource.Transaction) {
	if repo != nil {
		return
	}

	repo = &repositoryImpl{
		writer:      writer,
		reader:      reader,
		transaction: transaction,

		user:    NewUserRepository(writer, reader),
		store:   NewStoreRepository(writer, reader),
		product: NewProductRepository(writer, reader, transaction),
	}
}

func Get() Repository {
	return repo
}
