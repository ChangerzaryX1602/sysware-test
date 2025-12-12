package todo

import (
	"github.com/ChangerzaryX1602/sysware-test/pkg/domain"
	"github.com/ChangerzaryX1602/sysware-test/pkg/models"

	helpers "github.com/zercle/gofiber-helpers"
)

type todoService struct {
	repository domain.TodoRepository
}

func NewTodoService(repository domain.TodoRepository) domain.TodoService {
	return &todoService{repository: repository}
}

func (s *todoService) CreateTodo(todo models.Todo) []helpers.ResponseError {
	err := s.repository.CreateTodo(todo)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *todoService) GetTodo(id uint) (*models.Todo, []helpers.ResponseError) {
	todo, err := s.repository.GetTodo(id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return todo, nil
}

func (s *todoService) GetTodos(pagination models.Pagination, search models.Search) ([]models.Todo, *models.Pagination, *models.Search, []helpers.ResponseError) {
	todos, paginated, searched, err := s.repository.GetTodos(pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return todos, paginated, searched, nil
}

func (s *todoService) GetMyTodos(userId uint, pagination models.Pagination, search models.Search) ([]models.Todo, *models.Pagination, *models.Search, []helpers.ResponseError) {
	todos, paginated, searched, err := s.repository.GetMyTodos(userId, pagination, search)
	if err != nil {
		return nil, nil, nil, []helpers.ResponseError{*err}
	}
	return todos, paginated, searched, nil
}

func (s *todoService) GetMyTodo(userId uint, id uint) (*models.Todo, []helpers.ResponseError) {
	todo, err := s.repository.GetMyTodo(userId, id)
	if err != nil {
		return nil, []helpers.ResponseError{*err}
	}
	return todo, nil
}

func (s *todoService) UpdateTodo(id uint, todo models.Todo) []helpers.ResponseError {
	err := s.repository.UpdateTodo(id, todo)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *todoService) UpdateMyTodo(userId uint, id uint, todo models.Todo) []helpers.ResponseError {
	err := s.repository.UpdateMyTodo(userId, id, todo)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *todoService) DeleteTodo(id uint) []helpers.ResponseError {
	err := s.repository.DeleteTodo(id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}

func (s *todoService) DeleteMyTodo(userId uint, id uint) []helpers.ResponseError {
	err := s.repository.DeleteMyTodo(userId, id)
	if err != nil {
		return []helpers.ResponseError{*err}
	}
	return nil
}
