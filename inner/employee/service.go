package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
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
	BeginTransaction() (*sqlx.Tx, error)
	FindByNameTx(*sqlx.Tx, string) (bool, error)
	CreateTx(*sqlx.Tx, Entity) error
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

func (svc *Service) CreateEmployee(e Entity) (err error) {
	tx, err := svc.repo.BeginTransaction()

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
		return fmt.Errorf("error creating transaction: %w", err)
	}

	isExists, err := svc.repo.FindByNameTx(tx, e.Name)
	if err != nil {
		return fmt.Errorf("error finding employee by name: %w", err)
	}
	if isExists == false {
		err = svc.repo.CreateTx(tx, e)
	}
	if err != nil {
		return fmt.Errorf("error create employee: %w", err)
	}
	return nil
}
