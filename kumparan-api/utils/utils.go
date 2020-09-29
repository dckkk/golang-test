package utils

import "os"

func init() {
	LoadEnvs("")
}

// GetEnv is function to get data config from env variable
func GetEnv(key string, defaultValue string) string {
	result := os.Getenv(key)
	if result == "" {
		result = defaultValue
	}
	return result
}
