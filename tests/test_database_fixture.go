package tests

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/zhedevops/idm/inner/common"
	"github.com/zhedevops/idm/inner/database"
)

type FixtureDb struct {
	testDb *sqlx.DB
}

func NewFixtureDb() (*FixtureDb, error) {
	db := database.ConnectDb()
	// Создать базу для тестов
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'test_db')").Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		_, err = db.Exec("CREATE DATABASE test_db")
		if err != nil {
			return nil, err
		}
	}
	cfg, errStr := common.GetConfig(".env.test", true)
	if errStr != "" {
		return nil, errors.New(errStr)
	}
	testDb := database.ConnectDbWithCfg(cfg)
	return &FixtureDb{testDb}, nil
}

func (f *FixtureDb) CreateEmployeeTable() error {
	query := `CREATE TABLE IF NOT EXISTS employee (
              id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
              name TEXT NOT NULL,
              created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
              updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
          );`
	_, err := f.testDb.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (f *FixtureDb) CreateRoleTable() error {
	query := `CREATE TABLE IF NOT EXISTS role (
              id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
              name TEXT NOT NULL,
              created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
              updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
          );`
	_, err := f.testDb.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
