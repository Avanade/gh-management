package config

type Key string

type Config struct {
	DatabaseConnectionString string
}

type ConfigManager interface {
	GetApprovalSystemAppId() string
	GetApprovalSystemAppModuleAdoOrganization() string
	GetApprovalSystemAppUrl() string
	GetDatabaseConnectionString() string
	GetEmailTenantID() string
	GetEmailClientID() string
	GetEmailClientSecret() string
	GetEmailUserID() string
	GetIsEmailEnabled() bool
	GetLegalApprovalTypeId() int
}
