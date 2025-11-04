package tests

import (
	"github.com/zhedevops/idm/inner/employee"
)

type Fixture struct {
	employees *employee.Repository
}

func NewFixtureEmployee(employees *employee.Repository) *Fixture {
	return &Fixture{employees}
}

func (f *Fixture) Employee(name string) int64 {
	var entity = employee.Entity{
		Name: name,
	}
	err := f.employees.CreateNamed(&entity)
	if err != nil {
		return 0
	}
	return entity.Id
}
