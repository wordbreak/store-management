package service

import (
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

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
