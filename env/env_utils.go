package env

import (
	"os"
	"strconv"
)

func GetAppName() string {
	return GetEnv(LoggerAppName, "my-app")
}

func GetEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key) // Gets the environment variable as a string
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value // Returns the environment variable as an integer
	}

	return defaultValue // Returns a default value if the environment variable is not an integer
}

func GetEnv(key string, defaultValue string) string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return valueStr
}

func GetEnvBool(envName string, defaultValue bool) bool {
	envVal, exists := os.LookupEnv(envName)
	if !exists {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(envVal)
	if err != nil {
		return defaultValue
	}
	return boolVal
}
