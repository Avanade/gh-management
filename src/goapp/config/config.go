package config

type Key string

type Config struct {
	DatabaseConnectionString string
}

type ConfigManager interface {
	GetDatabaseConnectionString() string
}
