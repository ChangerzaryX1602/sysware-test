package auth

import (
	"github.com/ChangerzaryX1602/sysware-test/internal/handlers"
	"github.com/ChangerzaryX1602/sysware-test/pkg/domain"
	"github.com/ChangerzaryX1602/sysware-test/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type authHandler struct {
	service         domain.AuthService
	routerResources *handlers.RouterResources
}

func NewAuthHandler(router fiber.Router, routerResources *handlers.RouterResources, service domain.AuthService) {
	handler := &authHandler{service: service, routerResources: routerResources}
	router.Post("/login", handler.Login())
	router.Post("/register", handler.Register())
}

// Login godoc
// @Summary Login to the application
// @Description Login with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User Login Credentials"
// @Router /auth/login [post]
func (h *authHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		host := c.Hostname()
		token, errorForm := h.service.Login(user, host)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"token":   token,
		})
	}

}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User Registration Data"
// @Router /auth/register [post]
func (h *authHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid Request",
						Message: err.Error(),
					},
				},
			})
		}
		errorForm := h.service.Register(user)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.ResponseForm{
			Success: true,
		})
	}
}
