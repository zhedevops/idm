package common

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// Config общая конфигурация всего приложения
type Config struct {
	DbDriverName string `validate:"required"`
	Dsn          string `validate:"required"`
}

// GetConfig загружает конфигурацию из .env файла или переменных окружения.
// Параметр envFile — путь к файлу .env
// Параметр overwrite — если true, значения из файла перезапишут уже существующие переменные окружения
func GetConfig(envFile string, overwrite bool) (Config, string) {
	var emptyCfg = Config{}

	var driver, drv_ok = os.LookupEnv("DB_DRIVER_NAME")
	var dsn, dsn_ok = os.LookupEnv("DB_DSN")

	if dsn_ok && drv_ok {
		emptyCfg.DbDriverName = driver
		emptyCfg.Dsn = dsn
	}

	if envFile == "" {
		return emptyCfg, "envFile path cannot be empty"
	}

	// Проверяем, что файл существует
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return emptyCfg, "envFile does not exist"
	}

	// Загружаем .env
	var err error
	if overwrite {
		err = godotenv.Overload(envFile) // перезаписываем переменные
	} else {
		err = godotenv.Load(envFile) // не перезаписываем существующие
	}

	if err != nil {
		return Config{}, "failed to load env file"
	}

	var cfg = Config{
		DbDriverName: os.Getenv("DB_DRIVER_NAME"),
		Dsn:          os.Getenv("DB_DSN"),
	}
	fmt.Printf("DB_DRIVER_NAME=%s, DB_DSN=%s\n", cfg.DbDriverName, cfg.Dsn)
	// Проверяем, что переменные окружения заполнены
	if cfg.DbDriverName == "" || cfg.Dsn == "" {
		return Config{}, "required environment variables are missing"
	}

	return cfg, ""
}
