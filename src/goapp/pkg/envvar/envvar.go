package envvar

import (
	"os"
)

func GetEnvVar(envVarName string, defaultValue string) string {

	if os.Getenv(envVarName) != "" {
		return os.Getenv(envVarName)
	}
	return defaultValue

}
