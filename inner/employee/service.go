package employee

import (
	"fmt"
)

// Структура сервиса, которая будет инкапсулировать бизнес-логику
type Service struct {
	repo Repo
}

// Согласно идеологии Go:
// - "принимайте интерфейсы и возвращайте структуры",
// - "объявляйте интерфейсы там, где вы собираетесь их использовать"
type Repo interface {
	FindById(id int64) (Entity, error)
	Create(*Entity) error
	CreateNamed(*Entity) error
	FindAll() ([]Entity, error)
	FilterByIDs([]int64) ([]Entity, error)
	DeleteById(int64) (int64, error)
	DeleteByIds([]int64) (int64, error)
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (srv *Service) FindById(id int64) (Response, error) {
	var entity, err = srv.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", id, err)
	}

	return entity.toResponse(), nil
}

func (srv *Service) Create(e Entity) error {
	var err = srv.repo.Create(&e)
	if err != nil {
		return fmt.Errorf("employee not created: %w", err)
	}

	return nil
}

func (srv *Service) CreateNamed(e Entity) error {
	var err = srv.repo.CreateNamed(&e)
	if err != nil {
		return fmt.Errorf("employee not created: %w", err)
	}

	return nil
}

func (srv *Service) FindAll() ([]Response, error) {
	var entities, err = srv.repo.FindAll()
	if err != nil {
		return []Response{}, fmt.Errorf("error get all employees: %w", err)
	}

	var resp []Response
	for _, e := range entities {
		resp = append(resp, e.toResponse())
	}
	return resp, nil
}

func (srv *Service) FilterByIDs(ids []int64) ([]Response, error) {
	var entities, err = srv.repo.FilterByIDs(ids)
	if err != nil {
		return []Response{}, fmt.Errorf("error get employees by ids: %w", err)
	}

	var resp []Response
	for _, e := range entities {
		resp = append(resp, e.toResponse())
	}
	return resp, nil
}

func (srv *Service) DeleteById(id int64) (int64, error) {
	var count, err = srv.repo.DeleteById(id)
	if err != nil {
		return 0, fmt.Errorf("error delete employee by id: %w", err)
	}

	return count, nil
}

func (srv *Service) DeleteByIds(ids []int64) (int64, error) {
	var count, err = srv.repo.DeleteByIds(ids)
	if err != nil {
		return 0, fmt.Errorf("error delete employee by ids: %w", err)
	}

	return count, nil
}
