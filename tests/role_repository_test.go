package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhedevops/idm/inner/role"
	"testing"
)

func TestRoleRepository(t *testing.T) {
	a := assert.New(t)
	fixtureDb, err := NewFixtureDb()
	a.Nil(err, "expected error to be nil")
	err = fixtureDb.CreateRoleTable()
	a.Nil(err, "expected error to be nil")
	db := fixtureDb.testDb

	var clearDatabase = func() {
		db.MustExec("DELETE FROM role")
	}
	defer func() {
		if r := recover(); r != nil {
			clearDatabase()
		}
	}()
	var Repository = role.NewRepository(db)
	var fixture = NewFixtureRole(Repository)

	var newRoleId int64
	t.Run("Create Role and FindById", func(t *testing.T) {
		newRoleId = fixture.Role("John Doe")
		var role, err = Repository.FindById(newRoleId)
		a.Nil(err, "expected error to be nil")
		a.Equal(role.Name, "John Doe")
		a.Equal(role.Id, newRoleId)
		a.NotEmpty(role.CreatedAt, "CreatedAt must be not empty")
		a.NotEmpty(role.UpdatedAt, "UpdatedAt must be not empty")
	})

	var ids []int64
	t.Run("Create roles, get All and FilterByIDs", func(t *testing.T) {
		var newRoleId2 = fixture.Role("John Deer")
		var newRoleId3 = fixture.Role("John Smith")
		var roles, err = Repository.FindAll()
		a.Nil(err, "expected error to be nil")
		for _, e := range roles {
			if e.Id == newRoleId2 || e.Id == newRoleId3 {
				ids = append(ids, e.Id)
			}
		}

		a.Contains(ids, newRoleId2, "expected roles to contain newRoleId")
		a.Contains(ids, newRoleId3, "expected roles to contain newRoleId2")
		roles, err = Repository.FilterByIDs(ids)
		a.Nil(err, "expected error to be nil")
		for _, e := range roles {
			if e.Id == newRoleId2 {
				a.Equal(e.Name, "John Deer")
			}
			if e.Id == newRoleId3 {
				a.Equal(e.Name, "John Smith")
			}
		}
	})

	t.Run("Delete roles", func(t *testing.T) {
		fmt.Println(ids)
		var count, err = Repository.DeleteByIds(ids)
		a.Nil(err, "expected error to be nil")
		a.Equal(count, int64(2), "expected count to be 2")
		count, err = Repository.DeleteById(newRoleId)
		a.Nil(err, "expected error to be nil")
		a.Equal(count, int64(1), "expected count to be 1")
	})

	clearDatabase()
}
