package role

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

// реализуем интерфейс репозитория у мока
func (m *MockRepo) FindById(id int64) (Entity, error) {
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

func TestFindById(t *testing.T) {
	var a = assert.New(t)

	t.Run("should return found role", func(t *testing.T) {
		// создаём экземпляр мок-объекта
		var repo = new(MockRepo)
		// создаём экземпляр сервиса, который собираемся тестировать. Передаём в его конструктор мок вместо реального репозитория
		var svc = NewService(repo)
		// создаём Entity, которую должен вернуть репозиторий
		var entity = Entity{
			Id:        1,
			Name:      "John Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// создаём Response, который ожидаем получить от сервиса
		var want = entity.toResponse()
		// конфигурируем поведение мок-репозитория (при вызове метода FindById с аргументом 1 вернуть Entity, созданную нами выше)
		repo.On("FindById", int64(1)).Return(entity, nil)
		// вызываем сервис с аргументом id = 1
		var got, err = svc.FindById(1)
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
		var svc = NewService(repo)
		// создаём пустую структуру role.Entity, которую сервис вернёт вместе с ошибкой
		var entity = Entity{}
		// ошибка, которую вернёт репозиторий
		var err = errors.New("database error")
		// ошибка, которую должен будет вернуть сервис
		var want = fmt.Errorf("error finding role with id 1: %w", err)
		repo.On("FindById", int64(1)).Return(entity, err)
		var response, got = svc.FindById(1)
		// проверяем результаты теста
		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})
}

func TestCreateNamed(t *testing.T) {
	var a = assert.New(t)
	var repo = new(MockRepo)
	var svc = NewService(repo)
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
		var want = fmt.Errorf("role not created: %w", err)
		repo.On("CreateNamed", &entity).Return(err)
		var got = svc.CreateNamed(entity)
		a.NotNil(err)
		a.Equal(want, got)
	})
}

func TestFindAll(t *testing.T) {
	var a = assert.New(t)
	t.Run("found roles", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
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
	t.Run("not found roles", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
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
	var svc = NewService(repo)
	t.Run("found roles", func(t *testing.T) {
		var ids = []int64{1, 2}
		var want []Response
		for _, e := range entities {
			want = append(want, e.toResponse())
		}
		repo.On("FilterByIDs", ids).Return(entities, nil)
		var response, err = svc.FilterByIDs(ids)
		a.Nil(err)
		a.Equal(want, response)
	})
	t.Run("not found roles", func(t *testing.T) {
		var ids = []int64{3, 4}
		var err = errors.New("not found roles")
		var want = fmt.Errorf("error get roles by ids: %w", err)
		repo.On("FilterByIDs", ids).Return([]Entity{}, err)
		var response, got = svc.FilterByIDs(ids)
		a.NotNil(err)
		a.Equal(want, got)
		a.Equal(response, []Response{})
	})
}

func TestDeleteById(t *testing.T) {
	var a = assert.New(t)
	var repo = new(MockRepo)
	var svc = NewService(repo)
	t.Run("delete role", func(t *testing.T) {
		repo.On("DeleteById", int64(1)).Return(int64(1), nil)
		var response, err = svc.DeleteById(1)
		a.Nil(err)
		a.Equal(int64(1), response)
	})
	t.Run("error on delete role", func(t *testing.T) {
		var err = errors.New("not found role")
		var want = fmt.Errorf("error delete role by id: %w", err)
		repo.On("DeleteById", int64(3)).Return(int64(0), want)
		var response, got = svc.DeleteById(3)
		a.NotNil(got)
		a.Equal(int64(0), response)
	})
}

func TestDeleteByIds(t *testing.T) {
	var a = assert.New(t)
	var repo = new(MockRepo)
	var svc = NewService(repo)
	t.Run("delete roles", func(t *testing.T) {
		var ids = []int64{1, 2}
		repo.On("DeleteByIds", ids).Return(int64(2), nil)
		var response, err = svc.DeleteByIds(ids)
		a.Nil(err)
		a.Equal(int64(2), response)
	})
	t.Run("error on delete roles", func(t *testing.T) {
		var ids []int64
		var err = errors.New("not found roles")
		var want = fmt.Errorf("error delete roles by ids: %w", err)
		repo.On("DeleteByIds", ids).Return(int64(0), want)
		var response, got = svc.DeleteByIds(ids)
		a.NotNil(got)
		a.Equal(int64(0), response)
	})
}
