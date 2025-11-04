package tests

import (
	"github.com/zhedevops/idm/inner/role"
)

type FixtureRole struct {
	roles *role.Repository
}

func NewFixtureRole(roles *role.Repository) *FixtureRole {
	return &FixtureRole{roles}
}

func (f *FixtureRole) Role(name string) int64 {
	var entity = role.Entity{
		Name: name,
	}
	err := f.roles.CreateNamed(&entity)
	if err != nil {
		return 0
	}
	return entity.Id
}
