package auth

import (
	"github.com/ChangerzaryX1602/sysware-test/pkg/domain"
	"github.com/ChangerzaryX1602/sysware-test/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type authService struct {
	repository domain.AuthRepository
	service    domain.UserService
}

func NewAuthService(repository domain.AuthRepository, service domain.UserService) domain.AuthService {
	return &authService{repository: repository, service: service}
}
func (s *authService) Login(user models.User, host string) (*string, []helpers.ResponseError) {
	userE, errorFormList := s.service.GetUserByEmail(user.Email)
	if errorFormList != nil {
		return nil, errorFormList
	}
	checked := models.CheckPasswordHash(user.PasswordTemp, userE.Password)
	if !checked {
		return nil, []helpers.ResponseError{
			{
				Code:    401,
				Source:  helpers.WhereAmI(),
				Title:   "Unauthorized",
				Message: "Invalid credentials",
			},
		}
	}
	token, errorForm := s.repository.SignToken(*userE, host)
	if errorForm != nil {
		return nil, []helpers.ResponseError{*errorForm}
	}
	return &token, nil
}
func (s *authService) Register(user models.User) []helpers.ResponseError {
	err := s.service.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) GetUserByID(userId uint) (*models.User, []helpers.ResponseError) {

	user, err := s.service.GetUser(userId)
	if err != nil {
		return nil, err
	}
	// Clear password before returning
	user.Password = ""
	user.PasswordTemp = ""
	return user, nil
}
