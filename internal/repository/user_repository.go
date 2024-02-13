package repository

import (
	"database/sql"
	"errors"
	"store-management/internal/datasource"
	"store-management/internal/model"
	"time"

	"github.com/go-sql-driver/mysql"
)

type UserRepository interface {
	CreateUser(phoneNumber, password string) error
	FindUser(phoneNumber string) (*model.User, error)
	FindUserByID(id int64) (*model.User, error)
	BlockAuthToken(string) error
}

type userRepositoryImpl struct {
	writer datasource.SQL
	reader datasource.SQL
	cache  datasource.Cache
}

func NewUserRepository(writer, reader datasource.SQL, cache datasource.Cache) UserRepository {
	return &userRepositoryImpl{
		writer: writer,
		reader: reader,
		cache:  cache,
	}
}

func (u *userRepositoryImpl) CreateUser(phoneNumber, password string) error {
	_, err := u.writer.Exec("INSERT INTO user (phone_number, password) VALUES (?, ?)", phoneNumber, password)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			if mysqlError.Number == datasource.MySQLDuplicateEntry {
				return datasource.ErrDuplicateEntry
			}
		}
		panic(err)
	}
	return nil
}

func (u *userRepositoryImpl) FindUser(phoneNumber string) (*model.User, error) {
	var user model.User
	err := u.reader.Get(&user, "SELECT id, phone_number, password FROM user WHERE phone_number = ?", phoneNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datasource.ErrNoRows
		}
		panic(err)
	}
	return &user, err
}

func (u *userRepositoryImpl) FindUserByID(id int64) (*model.User, error) {
	var user model.User
	err := u.reader.Get(&user, "SELECT id, phone_number, password FROM user WHERE id = ?", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datasource.ErrNoRows
		}
		panic(err)
	}
	return &user, err
}

func (u *userRepositoryImpl) BlockAuthToken(authToken string) error {
	_ = u.cache.Set(authToken, struct{}{}, time.Hour)
	return nil
}
