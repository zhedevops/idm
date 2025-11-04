package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhedevops/idm/inner/employee"
	"testing"
)

func TestEmployeeRepository(t *testing.T) {
	a := assert.New(t)
	fixtureDb, err := NewFixtureDb()
	a.Nil(err, "expected error to be nil")
	err = fixtureDb.CreateEmployeeTable()
	a.Nil(err, "expected error to be nil")
	db := fixtureDb.testDb

	var clearDatabase = func() {
		db.MustExec("DELETE FROM employee")
	}
	defer func() {
		if r := recover(); r != nil {
			clearDatabase()
		}
	}()
	var Repository = employee.NewRepository(db)
	var fixture = NewFixtureEmployee(Repository)

	var newEmployeeId int64
	t.Run("Create Employee and FindById", func(t *testing.T) {
		newEmployeeId = fixture.Employee("John Doe")
		var employee, err = Repository.FindById(newEmployeeId)
		a.Nil(err, "expected error to be nil")
		a.Equal(employee.Name, "John Doe")
		a.Equal(employee.Id, newEmployeeId)
		a.NotEmpty(employee.CreatedAt, "CreatedAt must be not empty")
		a.NotEmpty(employee.UpdatedAt, "UpdatedAt must be not empty")
	})

	var ids []int64
	t.Run("Create Employees, get All and FilterByIDs", func(t *testing.T) {
		var newEmployeeId2 = fixture.Employee("John Deer")
		var newEmployeeId3 = fixture.Employee("John Smith")
		var employees, err = Repository.FindAll()
		a.Nil(err, "expected error to be nil")
		for _, e := range employees {
			if e.Id == newEmployeeId2 || e.Id == newEmployeeId3 {
				ids = append(ids, e.Id)
			}
		}

		a.Contains(ids, newEmployeeId2, "expected employees to contain newEmployeeId")
		a.Contains(ids, newEmployeeId3, "expected employees to contain newEmployeeId2")
		employees, err = Repository.FilterByIDs(ids)
		a.Nil(err, "expected error to be nil")
		for _, e := range employees {
			if e.Id == newEmployeeId2 {
				a.Equal(e.Name, "John Deer")
			}
			if e.Id == newEmployeeId3 {
				a.Equal(e.Name, "John Smith")
			}
		}
	})

	t.Run("Delete Employees", func(t *testing.T) {
		fmt.Println(ids)
		var count, err = Repository.DeleteByIds(ids)
		a.Nil(err, "expected error to be nil")
		a.Equal(count, int64(2), "expected count to be 2")
		count, err = Repository.DeleteById(newEmployeeId)
		a.Nil(err, "expected error to be nil")
		a.Equal(count, int64(1), "expected count to be 1")
	})

	clearDatabase()
}
