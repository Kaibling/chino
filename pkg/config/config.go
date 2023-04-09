package config

import (
	"os"
)

type Config struct {
	BindingPort      string
	BindingIP        string
	DatabaseFilePath string
}

var osprefix = "CHINO"

func Init() *Config {
	return &Config{
		BindingIP:        getEnv("BINDING_IP", "0.0.0.0"),
		DatabaseFilePath: getEnv("DATABASE_PATH", "./chino.db"),
		BindingPort:      getEnv("BINDING_PORT", "3000"),
	}
}

func getEnv(key string, defaultValue string) string {
	fullKey := osprefix + "_" + key
	val := os.Getenv(osprefix + "_" + key)
	if val == "" {
		if defaultValue != "" {
			return defaultValue
		}
		panic(fullKey + " is not set")
	}
	return val
}
