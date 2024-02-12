package service

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"errors"
	"store-management/internal/datasource"
	"store-management/internal/model"
	"store-management/internal/repository"

	"golang.org/x/crypto/argon2"
)

var (
	ErrDuplicateUser = errors.New("duplicate user")
	ErrUserNotFound  = errors.New("user not found")
)

type AuthService interface {
	Register(phoneNumber string, password string) error
	Login(phoneNumber string, password string) (*model.User, error)
}

type authServiceImpl struct {
	repo struct {
		user repository.UserRepository
	}
}

type Argon2idHash struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

var argon2IDHash = &Argon2idHash{
	time:    1,
	saltLen: 32,
	memory:  64 * 1024,
	threads: 32,
	keyLen:  256,
}

func (a *Argon2idHash) GenerateHash(password, salt []byte) []byte {
	if len(salt) == 0 {
		panic("salt cannot be empty")
	}
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)
	return hash
}

func (a *Argon2idHash) Compare(hash, salt, password []byte) bool {
	generatedHash := a.GenerateHash(password, salt)
	return bytes.Equal(hash, generatedHash)
}

func (s authServiceImpl) Register(phoneNumber string, password string) error {
	user, err := s.repo.user.FindUser(phoneNumber)
	if err != nil && !errors.Is(err, datasource.ErrNoRows) {
		panic(err)
	}
	if user != nil {
		return ErrDuplicateUser
	}

	hash := argon2IDHash.GenerateHash([]byte(password), []byte(phoneNumber))
	encryptedPassword := hex.EncodeToString(hash)

	err = s.repo.user.CreateUser(phoneNumber, encryptedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDuplicateUser
		}
		panic(err)
	}

	return nil
}

func (s authServiceImpl) Login(phoneNumber string, password string) (*model.User, error) {
	user, err := s.repo.user.FindUser(phoneNumber)
	if err != nil {
		if errors.Is(err, datasource.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		panic(err)
	}

	storedPassword, err := hex.DecodeString(user.Password)
	if err != nil {
		panic(err)
	}
	salt := []byte(user.PhoneNumber)
	bytedPassword := []byte(password)

	if passwordMatched := argon2IDHash.Compare(storedPassword, salt, bytedPassword); passwordMatched {
		return user, nil
	} else {
		return nil, ErrUserNotFound
	}
}

func NewAuthService(repo repository.Repository) AuthService {
	service := authServiceImpl{}
	service.repo.user = repo.UserRepository()
	return service
}
