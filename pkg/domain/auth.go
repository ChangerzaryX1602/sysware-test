package domain

import (
	"github.com/ChangerzaryX1602/sysware-test/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type AuthRepository interface {
	SignToken(user models.User, host string) (string, *helpers.ResponseError)
}
type AuthService interface {
	Register(user models.User) []helpers.ResponseError
	Login(user models.User, host string) (*string, []helpers.ResponseError)
	GetUserByID(userId uint) (*models.User, []helpers.ResponseError)
}
