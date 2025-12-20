package employee

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zhedevops/idm/inner/validator"
)

type MockRepo struct {
	mock.Mock
}

// реализуем интерфейс репозитория у мока
func (m *MockRepo) FindById(id int64) (employee Entity, err error) {
	// Общая конфигурация поведения мок-объекта
	args := m.Called(id)
	return args.Get(0).(Entity), args.Error(1)
}

func (m *MockRepo) Create(e *Entity) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockRepo) CreateNamed(e *Entity) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *MockRepo) FindAll() ([]Entity, error) {
	args := m.Called()
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) FilterByIDs(ids []int64) ([]Entity, error) {
	args := m.Called(ids)
	return args.Get(0).([]Entity), args.Error(1)
}

func (m *MockRepo) DeleteById(id int64) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) DeleteByIds(ids []int64) (int64, error) {
	args := m.Called(ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) BeginTransaction() (*sqlx.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (m *MockRepo) FindByNameTx(tx *sqlx.Tx, name string) (bool, error) {
	args := m.Called(tx, name)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockRepo) CreateTx(tx *sqlx.Tx, request CreateRequest) (int64, error) {
	args := m.Called(tx, request)
	return args.Get(0).(int64), args.Error(1)
}

func TestFindById(t *testing.T) {
	var a = assert.New(t)
	var validator = validator.New()

	t.Run("should return found employee", func(t *testing.T) {
		// создаём экземпляр мок-объекта
		var repo = new(MockRepo)
		// создаём экземпляр сервиса, который собираемся тестировать. Передаём в его конструктор мок вместо реального репозитория
		var svc = NewService(repo, validator)
		// создаём Entity, которую должен вернуть репозиторий
		var entity = Entity{
			Id:        1,
			Name:      "John Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		req := ParamIdRequest{Id: 1}
		// создаём Response, который ожидаем получить от сервиса
		var want = entity.toResponse()
		// конфигурируем поведение мок-репозитория (при вызове метода FindById с аргументом 1 вернуть Entity, созданную нами выше)
		repo.On("FindById", req.Id).Return(entity, nil)
		// вызываем сервис с аргументом id = 1
		var got, err = svc.FindById(req)
		// проверяем, что сервис не вернул ошибку
		a.Nil(err)
		// проверяем, что сервис вернул нам тот employee.Response, который мы ожилали получить
		a.Equal(want, got)
		// проверяем, что сервис вызвал репозиторий ровно 1 раз
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		// Создаём для теста новый экземпляр мока репозитория.
		// Мы собираемся проверить счётчик вызовов, поэтому хотим, чтобы счётчик содержал количество вызовов к репозиторию,
		// выполненных в рамках одного нашего теста.
		// Ели сделать мок общим для нескольких тестов, то он посчитает вызовы, которые сделали все тесты
		var repo = new(MockRepo)
		var svc = NewService(repo, validator)
		// создаём пустую структуру employee.Entity, которую сервис вернёт вместе с ошибкой
		var entity = Entity{}
		req := ParamIdRequest{Id: 1}
		// ошибка, которую вернёт репозиторий
		var err = errors.New("database error")
		// ошибка, которую должен будет вернуть сервис
		var wantErr = fmt.Errorf("error finding employee with id 1: %w", err)
		repo.On("FindById", req.Id).Return(entity, err)
		var response, gotErr = svc.FindById(req)
		// проверяем результаты теста
		a.Empty(response)
		a.NotNil(gotErr)
		a.Equal(wantErr, gotErr)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})
}

func TestCreateNamed(t *testing.T) {
	var a = assert.New(t)
	var validator = validator.New()
	var repo = new(MockRepo)
	var svc = NewService(repo, validator)
	t.Run("error is nil", func(t *testing.T) {
		var entity = Entity{
			Name: "Grigory Leps",
		}
		repo.On("CreateNamed", &entity).Return(nil)
		var err = svc.CreateNamed(entity)
		a.Nil(err)
	})
	t.Run("error on creating", func(t *testing.T) {
		var entity = Entity{}
		var err = errors.New("database error")
		var want = fmt.Errorf("employee not created: %w", err)
		repo.On("CreateNamed", &entity).Return(err)
		var got = svc.CreateNamed(entity)
		a.NotNil(err)
		a.Equal(want, got)
	})
}

func TestFindAll(t *testing.T) {
	var a = assert.New(t)
	var validator = validator.New()
	t.Run("found employees", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator)
		var entity1 = Entity{
			Id:        1,
			Name:      "Grigory Leps",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		var entity2 = Entity{
			Id:        2,
			Name:      "Semen Altow",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		var entities = []Entity{entity1, entity2}
		var want []Response
		for _, e := range entities {
			want = append(want, e.toResponse())
		}
		repo.On("FindAll").Return(entities, nil)
		var response, err = svc.FindAll()
		a.Nil(err)
		a.Equal(want, response)
	})
	t.Run("not found employees", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator)
		var entities = []Entity{}
		var want []Response
		repo.On("FindAll").Return(entities, nil)
		var response, got = svc.FindAll()
		a.Nil(got)
		a.Equal(response, want)
	})
}

func TestFilterByIDs(t *testing.T) {
	var a = assert.New(t)
	var validator = validator.New()
	var entity1 = Entity{
		Id:        1,
		Name:      "Grigory Leps",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	var entity2 = Entity{
		Id:        2,
		Name:      "Semen Altow",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	var entities = []Entity{entity1, entity2}
	var repo = new(MockRepo)
	var svc = NewService(repo, validator)
	t.Run("found employees", func(t *testing.T) {
		var req = ParamIdsRequest{Ids: []int64{1, 2}}
		var want []Response
		for _, e := range entities {
			want = append(want, e.toResponse())
		}
		repo.On("FilterByIDs", req.Ids).Return(entities, nil)
		var response, err = svc.FilterByIDs(req)
		a.Nil(err)
		a.Equal(want, response)
	})
	t.Run("not found employees", func(t *testing.T) {
		var req = ParamIdsRequest{Ids: []int64{3, 4}}
		var err = errors.New("not found employees")
		var want = fmt.Errorf("error get employees by ids: %w", err)
		repo.On("FilterByIDs", req.Ids).Return([]Entity{}, err)
		var response, got = svc.FilterByIDs(req)
		a.NotNil(err)
		a.Equal(want, got)
		a.Equal(response, []Response{})
	})
}

func TestDeleteById(t *testing.T) {
	var a = assert.New(t)
	var validator = validator.New()
	var repo = new(MockRepo)
	var svc = NewService(repo, validator)
	t.Run("delete employee", func(t *testing.T) {
		var req = ParamIdRequest{Id: 1}
		repo.On("DeleteById", req.Id).Return(int64(1), nil)
		var response, err = svc.DeleteById(req)
		a.Nil(err)
		a.Equal(int64(1), response)
	})
	t.Run("error on delete employee", func(t *testing.T) {
		var req = ParamIdRequest{Id: 3}
		var err = errors.New("not found employee")
		var want = fmt.Errorf("error delete employee by id: %w", err)
		repo.On("DeleteById", req.Id).Return(int64(0), want)
		var response, got = svc.DeleteById(req)
		a.NotNil(got)
		a.Equal(int64(0), response)
	})
}

func TestDeleteByIds(t *testing.T) {
	var a = assert.New(t)
	var validator = validator.New()
	var repo = new(MockRepo)
	var svc = NewService(repo, validator)
	t.Run("delete employees", func(t *testing.T) {
		var req = ParamIdsRequest{Ids: []int64{1, 2}}
		repo.On("DeleteByIds", req.Ids).Return(int64(2), nil)
		var response, err = svc.DeleteByIds(req)
		a.Nil(err)
		a.Equal(int64(2), response)
	})
	t.Run("error on delete employees", func(t *testing.T) {
		var req = ParamIdsRequest{Ids: []int64{}}
		var err = errors.New("not found employees")
		var want = fmt.Errorf("error delete employee by ids: %w", err)
		repo.On("DeleteByIds", req.Ids).Return(int64(0), want)
		var response, got = svc.DeleteByIds(req)
		a.NotNil(got)
		a.Equal(int64(0), response)
	})
}

func TestCreateEmployee(t *testing.T) {
	var a = assert.New(t)
	var validator = validator.New()
	var request = CreateRequest{
		Name: "Uncle Bob",
	}

	t.Run("success begin transaction and create employee", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator)
		db, mock, err := sqlmock.New()
		a.Nil(err)
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		mock.ExpectBegin()  // ожидаем транзакцию
		mock.ExpectCommit() // ожидаем коммит

		tx, err := sqlxDB.Beginx()
		a.Nil(err)
		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, request.Name).Return(false, nil)
		repo.On("CreateTx", tx, request).Return(int64(1), nil)
		_, err = svc.CreateEmployee(request)
		a.Nil(err)
	})

	t.Run("failure begin transaction", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator)
		db, mock, err := sqlmock.New()
		a.Nil(err)
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		mock.ExpectBegin()    // ожидаем транзакцию
		mock.ExpectRollback() // ожидаем откат

		tx, err := sqlxDB.Beginx()
		a.Nil(err)
		err = errors.New("transaction not begin")
		var want = fmt.Errorf("error creating transaction: %w", err)
		repo.On("BeginTransaction").Return(tx, want)
		_, err = svc.CreateEmployee(request)
		a.NotNil(err)
	})

	t.Run("failure on FindByNameTx", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator)
		db, mock, err := sqlmock.New()
		a.Nil(err)
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		mock.ExpectBegin()    // ожидаем транзакцию
		mock.ExpectRollback() // ожидаем откат
		tx, err := sqlxDB.Beginx()
		a.Nil(err)
		var requestNone = CreateRequest{
			Name: "None",
		}
		err = errors.New("finding error")
		var want = fmt.Errorf("error finding employee by name: %w", err)
		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, requestNone.Name).Return(false, err)
		_, err = svc.CreateEmployee(requestNone)
		a.NotNil(err)
		a.Equal(want, err)
	})

	t.Run("entity already exists", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator)
		db, mock, err := sqlmock.New()
		a.Nil(err)
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		mock.ExpectBegin()    // ожидаем транзакцию
		mock.ExpectRollback() // ожидаем откат
		tx, err := sqlxDB.Beginx()
		a.Nil(err)
		err = errors.New("already exists")
		var want = fmt.Errorf("error finding employee by name: %w", err)
		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, request.Name).Return(true, err)
		_, err = svc.CreateEmployee(request)
		a.NotNil(err)
		a.Equal(want, err)
	})

	t.Run("error create employee", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator)
		db, mock, err := sqlmock.New()
		a.Nil(err)
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		mock.ExpectBegin()    // ожидаем транзакцию
		mock.ExpectRollback() // ожидаем откат
		tx, err := sqlxDB.Beginx()
		a.Nil(err)
		err = errors.New("something wrong")
		var want = fmt.Errorf("error create employee with name: %s %w", request.Name, err)
		repo.On("BeginTransaction").Return(tx, nil)
		repo.On("FindByNameTx", tx, request.Name).Return(false, nil)
		repo.On("CreateTx", tx, request).Return(int64(0), err)
		id, err := svc.CreateEmployee(request)
		a.NotNil(err)
		a.Equal(want, err)
		a.Equal(id, int64(0))
	})
}

// Тест сформированный через gotests и доработанный (с моками неудобно)
// А вот если есть много вариантов для проверки, то для тестирования лучше использовать таблицы (table-driven tests)
// с входящими данными и ожидаемыми результатами
func TestService_FindById(t *testing.T) {
	type fields struct {
		repo      *MockRepo
		validator Validator
	}
	type args struct {
		request ParamIdRequest
	}
	var repo = new(MockRepo)
	var validator = validator.New()
	req := ParamIdRequest{Id: 1}
	tm := time.Now()
	var err = errors.New("database error")
	tests := []struct {
		name    string
		fields  fields
		args    args
		entity  Entity
		want    Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "should return found employee",
			fields: fields{
				repo:      repo,
				validator: validator,
			},
			args: args{
				request: req,
			},
			entity: Entity{
				Id:        1,
				Name:      "John Doe",
				CreatedAt: tm,
				UpdatedAt: tm,
			},
			want: (&Entity{
				Id:        1,
				Name:      "John Doe",
				CreatedAt: tm,
				UpdatedAt: tm,
			}).toResponse(),
			wantErr: false,
		},
		{
			name: "should return wrapped error",
			fields: fields{
				repo:      repo,
				validator: validator,
			},
			args: args{
				request: ParamIdRequest{Id: 2},
			},
			entity:  Entity{},
			want:    Response{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &Service{
				repo:      tt.fields.repo,
				validator: tt.fields.validator,
			}
			if tt.wantErr {
				tt.fields.repo.On("FindById", tt.args.request.Id).Return(tt.entity, err)
			} else {
				tt.fields.repo.On("FindById", tt.args.request.Id).Return(tt.entity, nil)
			}
			got, err := srv.FindById(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.FindById() = %v, want %v", got, tt.want)
			}
		})
	}
}
