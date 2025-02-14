package config

type Key string

type Config struct {
	DatabaseConnectionString string
}

type ConfigManager interface {
	GetDatabaseConnectionString() string
	GetEmailTenantID() string
	GetEmailClientID() string
	GetEmailClientSecret() string
	GetEmailUserID() string
	GetIsEmailEnabled() bool
	GetLegalApprovalTypeId() int
}
