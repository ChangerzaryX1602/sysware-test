package auth

import (
	"fmt"
	"github.com/ChangerzaryX1602/sysware-test/pkg/domain"
	"github.com/ChangerzaryX1602/sysware-test/pkg/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	helpers "github.com/zercle/gofiber-helpers"
)

type authRepository struct {
	resources models.Resources
}

func NewAuthRepository(resources models.Resources) domain.AuthRepository {
	return &authRepository{resources: resources}
}
func (r *authRepository) SignToken(user models.User, host string) (string, *helpers.ResponseError) {
	token := jwt.NewWithClaims(r.resources.JwtResources.JwtSigningMethod, &jwt.RegisteredClaims{})
	claims := token.Claims.(*jwt.RegisteredClaims)
	claims.Subject = fmt.Sprintf("%d", user.ID)
	claims.Issuer = host
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	signToken, err := token.SignedString(r.resources.JwtResources.JwtSignKey)
	if err != nil {
		return "", &helpers.ResponseError{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Failed to sign json web token",
		}
	}
	return signToken, nil
}
