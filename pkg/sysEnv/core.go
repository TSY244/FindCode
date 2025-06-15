package sysEnv

import "os"

func SetEnv(key, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}
	return nil
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
