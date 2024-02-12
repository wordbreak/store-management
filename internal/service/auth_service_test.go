package service

import (
	"encoding/hex"
	"store-management/internal/datasource"
	"store-management/internal/model"
	mock2 "store-management/mock"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthServiceSuite struct {
	suite.Suite
	mockedUserRepository *mock2.UserRepositoryMock
	service              AuthService
}

func (s *AuthServiceSuite) SetupTest() {
	s.mockedUserRepository = &mock2.UserRepositoryMock{}
	mockRepo := mock2.NewMockedRepository(s.mockedUserRepository)
	s.service = NewAuthService(mockRepo)
}

func (s *AuthServiceSuite) TestRegister_DuplicateUser_Error() {
	phoneNumber := "1234567890"
	password := "password"

	s.mockedUserRepository.On("FindUser", phoneNumber).Return(&model.User{}, nil).Once()

	err := s.service.Register(phoneNumber, password)
	s.EqualError(err, ErrDuplicateUser.Error())

	s.mockedUserRepository.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestRegister_Success() {
	phoneNumber := "1234567890"
	password := "password"

	s.mockedUserRepository.On("FindUser", phoneNumber).Return(nil, datasource.ErrNoRows).Once()
	s.mockedUserRepository.On("CreateUser", phoneNumber, mock.Anything).Return(nil).Once()

	err := s.service.Register(phoneNumber, password)
	s.NoError(err)

	s.mockedUserRepository.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestLogin_UserNotFound_Error() {
	phoneNumber := "1234567890"
	password := "password"

	s.mockedUserRepository.On("FindUser", phoneNumber).Return(nil, datasource.ErrNoRows).Once()

	user, err := s.service.Login(phoneNumber, password)
	s.Nil(user)
	s.EqualError(err, ErrUserNotFound.Error())

	s.mockedUserRepository.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestLogin_Success() {
	phoneNumber := "1234567890"
	password := "password"
	encryptedPassword := argon2IDHash.GenerateHash([]byte(password), []byte(phoneNumber))

	user := &model.User{
		PhoneNumber: phoneNumber,
		Password:    hex.EncodeToString(encryptedPassword),
	}

	s.mockedUserRepository.On("FindUser", phoneNumber).Return(user, nil).Once()

	user, err := s.service.Login(phoneNumber, password)
	s.NoError(err)
	s.NotNil(user)
	s.Equal(user.PhoneNumber, phoneNumber)

	s.mockedUserRepository.AssertExpectations(s.T())
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
