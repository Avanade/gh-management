package config

import (
	"log"
	"os"
	"strconv" // Add this import

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

func (ecm *envConfigManager) GetApprovalSystemAppId() string {
	return os.Getenv("APPROVAL_SYSTEM_APP_ID")
}

func (ecm *envConfigManager) GetApprovalSystemAppModuleAdoOrganization() string {
	return os.Getenv("APPROVAL_SYSTEM_APP_MODULE_ADO_ORGANIZATION")
}

func (ecm *envConfigManager) GetApprovalSystemAppUrl() string {
	return os.Getenv("APPROVAL_SYSTEM_APP_URL")
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

func (ecm *envConfigManager) GetLegalApprovalTypeId() int {
	value, err := strconv.Atoi(os.Getenv("LEGAL_APPROVAL_TYPE_ID"))
	if err != nil {
		return 0
	}
	return value
}
