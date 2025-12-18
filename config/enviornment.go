package config

import (
	"fmt"
	"log"
	"os"

	"sitecrawler/newgo/utils"

	"github.com/joho/godotenv"
)

type DBConnectionParams struct {
	Host     string
	User     string
	Password string
	DBName   string
	DBPort   string
	DBSchema string
}

type ServiceParams struct {
	Port string
	Env  string
}

type EnvVariables struct {
	ServiceParams      *ServiceParams
	DBConnectionParams *DBConnectionParams
}

func LoadEnvVariables() {
	// to run in dev mode supply the param dev to runner as
	// $ go run main.go dev
	args := os.Args[1:]

	if args != nil && len(args) > 0 && args[0] == "dev" {
		if err := godotenv.Load("dev.env"); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Default to loading `dev.env` for local runs if present.
	// If the file doesn't exist, we fall back to the process env vars.
	_ = godotenv.Load("dev.env")
}

func LoadEnvConfiguration() *EnvVariables {
	serviceParams := &ServiceParams{
		Port: LoadEnvStr(PORT),
		Env:  LoadEnvStr(Env),
	}
	dbConnectionParams := &DBConnectionParams{
		Host:     LoadEnvStr(DB_HOST),
		User:     LoadEnvStr(DB_USER),
		Password: LoadEnvStr(DB_PASSWORD),
		DBName:   LoadEnvStr(DB_NAME),
		DBPort:   LoadEnvStr(DB_PORT),
		DBSchema: LoadEnvStr(DB_SCHEMA),
	}

	res := EnvVariables{
		ServiceParams:      serviceParams,
		DBConnectionParams: dbConnectionParams,
	}

	return &res
}

func LoadEnvStr(key string) string {
	v := os.Getenv(key)
	if utils.IsEmpty(v) {
		panic(fmt.Sprintf("env var '%s' is either missing or empty", key))
	}
	return v
}
