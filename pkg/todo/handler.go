package todo

import (
	"strconv"

	"github.com/ChangerzaryX1602/sysware-test/internal/handlers"
	"github.com/ChangerzaryX1602/sysware-test/pkg/domain"
	"github.com/ChangerzaryX1602/sysware-test/pkg/models"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type todoHandler struct {
	service domain.TodoService
}

func NewTodoHandler(router fiber.Router, resource *handlers.RouterResources, service domain.TodoService) {
	handler := &todoHandler{service: service}
	router.Post("/me", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.Me)), handler.CreateMyTodo())
	router.Get("/me", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.Me)), handler.GetMyTodos())
	router.Get("/me/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.Me)), handler.GetMyTodo())
	router.Patch("/me/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.Me)), handler.UpdateMyTodo())
	router.Delete("/me/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.Me)), handler.DeleteMyTodo())
	router.Post("/", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.Create)), handler.CreateTodo())
	router.Get("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.Read)), handler.GetTodo())
	router.Get("/", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.List)), handler.GetTodos())
	router.Patch("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.Update)), handler.UpdateTodo())
	router.Delete("/:id", resource.ReqAuthPerms(models.PermissionGroupName(models.TodoGroup, models.Delete)), handler.DeleteTodo())
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param todo body models.Todo true "Todo Data"
// @Router /todos [post]
func (h *todoHandler) CreateTodo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo models.Todo
		if err := c.BodyParser(&todo); err != nil {
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
		if errorForm := h.service.CreateTodo(todo); errorForm != nil {
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

// CreateMyTodo godoc
// @Summary Create a new todo for me
// @Description Create a new todo for the authenticated user
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param todo body models.Todo true "Todo Data"
// @Router /todos/me [post]
func (h *todoHandler) CreateMyTodo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id").(string)
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid User ID",
						Message: err.Error(),
					},
				},
			})
		}
		var todo models.Todo
		if err := c.BodyParser(&todo); err != nil {
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
		todo.UserID = uint(userIdInt)
		if errorForm := h.service.CreateTodo(todo); errorForm != nil {
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

// GetTodo godoc
// @Summary Get a todo by ID
// @Description Get a todo by ID
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Router /todos/{id} [get]
func (h *todoHandler) GetTodo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
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
		todo, errorForm := h.service.GetTodo(uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    todo,
		})
	}
}

// GetTodos godoc
// @Summary Get all todos
// @Description Get all todos with pagination and search
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort order (asc/desc)"
// @Param keyword query string false "Search keyword"
// @Param column query string false "Search column"
// @Router /todos [get]
func (h *todoHandler) GetTodos() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pagination := models.Pagination{}
		search := models.Search{}
		if err := c.QueryParser(&pagination); err != nil {
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
		if err := c.QueryParser(&search); err != nil {
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
		todos, paginated, searched, errorForm := h.service.GetTodos(pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data: map[string]interface{}{
				"items":      todos,
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

// GetMyTodos godoc
// @Summary Get my todos
// @Description Get todos for the authenticated user
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort order (asc/desc)"
// @Param keyword query string false "Search keyword"
// @Param column query string false "Search column"
// @Router /todos/me [get]
func (h *todoHandler) GetMyTodos() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id").(string)
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid User ID",
						Message: err.Error(),
					},
				},
			})
		}
		pagination := models.Pagination{}
		search := models.Search{}
		if err := c.QueryParser(&pagination); err != nil {
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
		if err := c.QueryParser(&search); err != nil {
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
		todos, paginated, searched, errorForm := h.service.GetMyTodos(uint(userIdInt), pagination, search)
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data: map[string]interface{}{
				"items":      todos,
				"pagination": paginated,
				"search":     searched,
			},
		})
	}
}

// GetMyTodo godoc
// @Summary Get my todo by ID
// @Description Get a specific todo for the authenticated user
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Router /todos/me/{id} [get]
func (h *todoHandler) GetMyTodo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id").(string)
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid User ID",
						Message: err.Error(),
					},
				},
			})
		}
		id, err := c.ParamsInt("id")
		if err != nil {
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
		todo, errorForm := h.service.GetMyTodo(uint(userIdInt), uint(id))
		if errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    todo,
		})
	}
}

// UpdateMyTodo godoc
// @Summary Update my todo
// @Description Update a todo for the authenticated user
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Param todo body models.Todo true "Todo Data"
// @Router /todos/me/{id} [patch]
func (h *todoHandler) UpdateMyTodo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id").(string)
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid User ID",
						Message: err.Error(),
					},
				},
			})
		}
		id, err := c.ParamsInt("id")
		if err != nil {
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
		var todo models.Todo
		if err := c.BodyParser(&todo); err != nil {
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
		if errorForm := h.service.UpdateMyTodo(uint(userIdInt), uint(id), todo); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    todo,
		})
	}
}

// DeleteMyTodo godoc
// @Summary Delete my todo
// @Description Delete a todo for the authenticated user
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Router /todos/me/{id} [delete]
func (h *todoHandler) DeleteMyTodo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user_id").(string)
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Success: false,
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Source:  helpers.WhereAmI(),
						Title:   "Invalid User ID",
						Message: err.Error(),
					},
				},
			})
		}
		id, err := c.ParamsInt("id")
		if err != nil {
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
		if errorForm := h.service.DeleteMyTodo(uint(userIdInt), uint(id)); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
		})
	}
}

// UpdateTodo godoc
// @Summary Update a todo (Admin)
// @Description Update any todo
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Param todo body models.Todo true "Todo Data"
// @Router /todos/{id} [patch]
func (h *todoHandler) UpdateTodo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
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
		var todo models.Todo
		if err := c.BodyParser(&todo); err != nil {
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
		if errorForm := h.service.UpdateTodo(uint(id), todo); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
			Data:    todo,
		})
	}
}

// DeleteTodo godoc
// @Summary Delete a todo (Admin)
// @Description Delete any todo
// @Tags Todo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Router /todos/{id} [delete]
func (h *todoHandler) DeleteTodo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
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
		if errorForm := h.service.DeleteTodo(uint(id)); errorForm != nil {
			return c.Status(errorForm[0].Code).JSON(helpers.ResponseForm{
				Success: false,
				Errors:  errorForm,
			})
		}
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseForm{
			Success: true,
		})
	}
}
