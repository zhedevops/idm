package employee

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type StubRepoInterface interface {
	FilterByIds([]int64) ([]Entity, error)
}

type ServiceStub struct {
	repo StubRepo
}

type StubRepo struct {
	Employees []Entity // данные, которые будет возвращать stub
	Err       error    // опциональная ошибка, чтобы имитировать failure
}

func NewServiceStub(repo StubRepo) *ServiceStub {
	return &ServiceStub{
		repo: repo,
	}
}

func (s *ServiceStub) FilterByIds(ids []int64) ([]Entity, error) {
	if s.repo.Err != nil {
		return nil, s.repo.Err
	}

	var result []Entity
	for _, e := range s.repo.Employees {
		for _, id := range ids {
			if e.Id == id {
				result = append(result, e)
			}
		}
	}
	return result, nil
}

func TestFilterByIdsWithStub(t *testing.T) {
	var a = assert.New(t)
	t.Run("found employees in stub", func(t *testing.T) {
		stub := StubRepo{
			Employees: []Entity{
				{Id: 1, Name: "Alice Mitchel", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{Id: 2, Name: "Bob Wanger", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{Id: 3, Name: "Charlie Sheen", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			},
			Err: nil,
		}
		var srv = NewServiceStub(stub)
		var ids = []int64{1, 3}
		var want []Entity
		for _, e := range stub.Employees {
			for _, id := range ids {
				if e.Id == id {
					want = append(want, e)
				}
			}
		}
		var response, err = srv.FilterByIds(ids)
		a.Nil(err)
		a.Equal(want, response)
	})
	t.Run("not found employees in stub", func(t *testing.T) {
		var err = errors.New("not found employees")
		var want = fmt.Errorf("error get employees by ids: %w", err)
		var entity []Entity
		stub := StubRepo{
			Employees: entity,
			Err:       want,
		}
		var srv = NewServiceStub(stub)
		var ids = []int64{4, 5}
		var response, got = srv.FilterByIds(ids)
		a.NotNil(err)
		a.Equal(want, got)
		a.Equal(entity, response)
	})
}
