package todo

import (
	"github.com/ChangerzaryX1602/sysware-test/pkg/domain"
	"github.com/ChangerzaryX1602/sysware-test/pkg/models"
	"github.com/ChangerzaryX1602/sysware-test/pkg/utils"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"gorm.io/gorm"
)

type todoRepository struct {
	resources models.Resources
}

func NewTodoRepository(resources models.Resources) domain.TodoRepository {
	return &todoRepository{resources: resources}
}

func (r *todoRepository) Migrate() error {
	return r.resources.MainDbConn.AutoMigrate(&models.Todo{})
}

func (r *todoRepository) CreateTodo(todo models.Todo) *helpers.ResponseError {
	if err := r.resources.MainDbConn.Create(&todo).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *todoRepository) GetTodo(id uint) (*models.Todo, *helpers.ResponseError) {
	var todo models.Todo
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: err.Error(),
			}
		}
		return nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return &todo, nil
}

func (r *todoRepository) GetTodos(pagination models.Pagination, search models.Search) ([]models.Todo, *models.Pagination, *models.Search, *helpers.ResponseError) {
	var todos []models.Todo
	db := r.resources.MainDbConn.Model(&models.Todo{})
	db = utils.ApplySearch(db, search)
	db = utils.ApplyPagination(db, &pagination, models.Todo{})
	if err := db.Find(&todos).Error; err != nil {
		return nil, nil, nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return todos, &pagination, &search, nil
}

func (r *todoRepository) GetMyTodos(userId uint, pagination models.Pagination, search models.Search) ([]models.Todo, *models.Pagination, *models.Search, *helpers.ResponseError) {
	var todos []models.Todo
	db := r.resources.MainDbConn.Model(&models.Todo{}).Where("user_id = ?", userId)
	db = utils.ApplySearch(db, search)
	db = utils.ApplyPagination(db, &pagination, models.Todo{})
	if err := db.Find(&todos).Error; err != nil {
		return nil, nil, nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return todos, &pagination, &search, nil
}

func (r *todoRepository) GetMyTodo(userId uint, id uint) (*models.Todo, *helpers.ResponseError) {
	var todo models.Todo
	if err := r.resources.MainDbConn.Where("id = ? AND user_id = ?", id, userId).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: "Todo not found or you do not have permission to access it",
			}
		}
		return nil, &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return &todo, nil
}

func (r *todoRepository) UpdateTodo(id uint, todo models.Todo) *helpers.ResponseError {
	var existingTodo models.Todo
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&existingTodo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: err.Error(),
			}
		}
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}

	if err := r.resources.MainDbConn.Model(&existingTodo).Updates(todo).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *todoRepository) UpdateMyTodo(userId uint, id uint, todo models.Todo) *helpers.ResponseError {
	var existingTodo models.Todo
	if err := r.resources.MainDbConn.Where("id = ? AND user_id = ?", id, userId).First(&existingTodo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: err.Error(),
			}
		}
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}

	if err := r.resources.MainDbConn.Model(&existingTodo).Updates(todo).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *todoRepository) DeleteTodo(id uint) *helpers.ResponseError {
	var todo models.Todo
	if err := r.resources.MainDbConn.Where("id = ?", id).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: err.Error(),
			}
		}
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}

	if err := r.resources.MainDbConn.Delete(&todo).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}

func (r *todoRepository) DeleteMyTodo(userId uint, id uint) *helpers.ResponseError {
	var todo models.Todo
	if err := r.resources.MainDbConn.Where("id = ? AND user_id = ?", id, userId).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &helpers.ResponseError{
				Code:    fiber.StatusNotFound,
				Source:  helpers.WhereAmI(),
				Title:   "Not Found",
				Message: err.Error(),
			}
		}
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}

	if err := r.resources.MainDbConn.Delete(&todo).Error; err != nil {
		return &helpers.ResponseError{
			Code:    fiber.StatusInternalServerError,
			Source:  helpers.WhereAmI(),
			Title:   "Database Error",
			Message: err.Error(),
		}
	}
	return nil
}
