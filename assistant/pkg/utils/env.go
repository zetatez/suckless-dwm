package utils

import "os"

func GetEnv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultVal
}

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}
