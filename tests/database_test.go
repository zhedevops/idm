package database_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhedevops/idm/inner/common"
	"github.com/zhedevops/idm/inner/database"
	"os"
	"testing"
)

func TestGetConfigWithoutEnv(t *testing.T) {
	a := assert.New(t)
	_, err := common.GetConfig("", false)
	a.Equal(err, "envFile path cannot be empty")
}

func TestGetConfigWithEnvOs(t *testing.T) {
	os.Setenv("DB_DRIVER_NAME", "driver")
	os.Setenv("DB_DSN", "dsn")
	a := assert.New(t)
	got, err := common.GetConfig("", false)
	a.Equal(err, "envFile path cannot be empty")
	a.Equal(got.DbDriverName, "driver")
	a.Equal(got.Dsn, "dsn")
}

func TestGetConfigWithOverwrite(t *testing.T) {
	a := assert.New(t)
	got, err := common.GetConfig("../.env", false)
	a.Empty(err, "envFile does not exist")
	a.Equal(got.DbDriverName, "driver")
	a.Equal(got.Dsn, "dsn")
	got, err = common.GetConfig("../.env", true)
	a.Empty(err, "envFile does not exist")
	a.NotEqual(got.DbDriverName, "driver")
	a.NotEqual(got.Dsn, "dsn")
	os.Unsetenv("DB_DRIVER_NAME")
	os.Unsetenv("DB_DSN")
}

func TestGetConfig(t *testing.T) {
	a := assert.New(t)
	got, err := common.GetConfig(".env.local", false)
	a.Equal(err, "envFile does not exist")
	a.Empty(got.DbDriverName, "DbDriverName must be empty")
	a.Empty(got.Dsn, "Dsn must be empty")
}

func TestGetRootConfig(t *testing.T) {
	a := assert.New(t)
	got, err := common.GetConfig("../.env", true)
	a.Empty(err, "err must not be empty")
	a.NotEmpty(got.DbDriverName, "DbDriverName must be empty")
	a.NotEmpty(got.Dsn, "Dsn must be not empty")
}

func TestConnectDbWithWrongCfg(t *testing.T) {
	a := assert.New(t)
	cfg := common.Config{
		DbDriverName: "invalid_driver",
		Dsn:          "invalid_dsn",
	}
	// Ловим панику
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		_ = database.ConnectDbWithCfg(cfg)
	}()

	a.True(didPanic, "expected panic when connecting with wrong config")
}

func TestConnectDb(t *testing.T) {
	a := assert.New(t)
	db := database.ConnectDb()
	a.NotNil(db, "expected db to be not nil")
}
