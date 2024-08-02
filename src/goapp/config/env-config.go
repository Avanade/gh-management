package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type envConfigManager struct {
	*Config
}

func NewEnvConfigManager() *envConfigManager {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}
	return &envConfigManager{}
}

func (ecm *envConfigManager) GetDatabaseConnectionString() string {
	return os.Getenv("GHMGMTDB_CONNECTION_STRING")
}
