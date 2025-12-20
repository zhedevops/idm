package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zhedevops/idm/inner/common"
)

// Структура сервиса, которая будет инкапсулировать бизнес-логику
type Service struct {
	repo      Repo
	validator Validator
}

type CreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=155"`
}

type ParamIdRequest struct {
	Id int64 `validate:"required,gt=0"`
}

type ParamIdsRequest struct {
	Ids []int64 `validate:"required,min=1,dive,gt=0"`
}

type Validator interface {
	Validate(request any) error
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
	BeginTransaction() (*sqlx.Tx, error)
	FindByNameTx(*sqlx.Tx, string) (bool, error)
	CreateTx(*sqlx.Tx, CreateRequest) (int64, error)
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{
		repo:      repo,
		validator: validator,
	}
}

func (req *CreateRequest) ToEntity() Entity {
	return Entity{Name: req.Name}
}

func (srv *Service) FindById(request ParamIdRequest) (Response, error) {
	var err = srv.validator.Validate(request)
	if err != nil {
		return Response{}, common.RequestValidationError{Message: err.Error()}
	}
	entity, err := srv.repo.FindById(request.Id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", request.Id, err)
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

func (srv *Service) FilterByIDs(request ParamIdsRequest) ([]Response, error) {
	var err = srv.validator.Validate(request)
	if err != nil {
		return []Response{}, common.RequestValidationError{Message: err.Error()}
	}
	entities, err := srv.repo.FilterByIDs(request.Ids)
	if err != nil {
		return []Response{}, fmt.Errorf("error get employees by ids: %w", err)
	}

	var resp []Response
	for _, e := range entities {
		resp = append(resp, e.toResponse())
	}
	return resp, nil
}

func (srv *Service) DeleteById(request ParamIdRequest) (int64, error) {
	var err = srv.validator.Validate(request)
	if err != nil {
		return 0, common.RequestValidationError{Message: err.Error()}
	}
	count, err := srv.repo.DeleteById(request.Id)
	if err != nil {
		return 0, fmt.Errorf("error delete employee by id: %w", err)
	}

	return count, nil
}

func (srv *Service) DeleteByIds(request ParamIdsRequest) (int64, error) {
	var err = srv.validator.Validate(request)
	if err != nil {
		return 0, common.RequestValidationError{Message: err.Error()}
	}
	count, err := srv.repo.DeleteByIds(request.Ids)
	if err != nil {
		return 0, fmt.Errorf("error delete employee by ids: %w", err)
	}

	return count, nil
}

// Метод для создания нового сотрудника
// принимает на вход CreateRequest - структура запроса на создание сотрудника
func (srv *Service) CreateEmployee(request CreateRequest) (int64, error) {
	var err = srv.validator.Validate(request)
	if err != nil {
		// возвращаем кастомную ошибку в случае, если запрос не прошёл валидацию (про кастомные ошибки - дальше)
		return 0, common.RequestValidationError{Message: err.Error()}
	}

	tx, err := srv.repo.BeginTransaction()

	// отложенная функция завершения транзакции
	defer func() {
		// проверяем, не было ли паники
		if r := recover(); r != nil {
			err = fmt.Errorf("creating employee panic: %v", r)
			// если была паника, то откатываем транзакцию
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else if err != nil {
			// если произошла другая ошибка (не паника), то откатываем транзакцию
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else {
			// если ошибок нет, то коммитим транзакцию
			errTx := tx.Commit()
			if errTx != nil {
				err = fmt.Errorf("creating employee: commiting transaction error: %w", errTx)
			}
		}
	}()

	if err != nil {
		return 0, fmt.Errorf("error creating transaction: %w", err)
	}

	isExists, err := srv.repo.FindByNameTx(tx, request.Name)
	if err != nil {
		return 0, fmt.Errorf("error finding employee by name: %w", err)
	}
	if isExists {
		return 0, common.AlreadyExistsError{Message: fmt.Sprintf("employee with name %s already exists", request.Name)}
	}
	newEmployeeId, err := srv.repo.CreateTx(tx, request)
	if err != nil {
		return 0, fmt.Errorf("error create employee with name: %s %w", request.Name, err)
	}
	return newEmployeeId, nil
}
