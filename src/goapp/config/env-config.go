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

func (ecm *envConfigManager) GetEmailTenantID() string {
	return os.Getenv("EMAIL_TENANT_ID")
}

func (ecm *envConfigManager) GetEmailClientID() string {
	return os.Getenv("EMAIL_CLIENT_ID")
}

func (ecm *envConfigManager) GetEmailClientSecret() string {
	return os.Getenv("EMAIL_CLIENT_SECRET")
}

func (ecm *envConfigManager) GetEmailUserID() string {
	return os.Getenv("EMAIL_USER_ID")
}

func (ecm *envConfigManager) GetIsEmailEnabled() bool {
	if os.Getenv("EMAIL_ENABLED") != "true" {
		return false
	}
	return true
}
