package domain

import (
	"github.com/ChangerzaryX1602/sysware-test/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type TodoRepository interface {
	Migrate() error
	CreateTodo(todo models.Todo) *helpers.ResponseError
	GetTodo(id uint) (*models.Todo, *helpers.ResponseError)
	GetTodos(pagination models.Pagination, search models.Search) ([]models.Todo, *models.Pagination, *models.Search, *helpers.ResponseError)
	GetMyTodos(userId uint, pagination models.Pagination, search models.Search) ([]models.Todo, *models.Pagination, *models.Search, *helpers.ResponseError)
	GetMyTodo(userId uint, id uint) (*models.Todo, *helpers.ResponseError)
	UpdateTodo(id uint, todo models.Todo) *helpers.ResponseError
	UpdateMyTodo(userId uint, id uint, todo models.Todo) *helpers.ResponseError
	DeleteTodo(id uint) *helpers.ResponseError
	DeleteMyTodo(userId uint, id uint) *helpers.ResponseError
}

type TodoService interface {
	CreateTodo(todo models.Todo) []helpers.ResponseError
	GetTodo(id uint) (*models.Todo, []helpers.ResponseError)
	GetTodos(pagination models.Pagination, search models.Search) ([]models.Todo, *models.Pagination, *models.Search, []helpers.ResponseError)
	GetMyTodo(userId uint, id uint) (*models.Todo, []helpers.ResponseError)
	GetMyTodos(userId uint, pagination models.Pagination, search models.Search) ([]models.Todo, *models.Pagination, *models.Search, []helpers.ResponseError)
	UpdateTodo(id uint, todo models.Todo) []helpers.ResponseError
	UpdateMyTodo(userId uint, id uint, todo models.Todo) []helpers.ResponseError
	DeleteTodo(id uint) []helpers.ResponseError
	DeleteMyTodo(userId uint, id uint) []helpers.ResponseError
}
