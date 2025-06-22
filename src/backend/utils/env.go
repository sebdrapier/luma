package utils

import "os"

type Config struct {
	ServerPort   string
	DMXPort      string
	DataFilePath string
	EnableDMX    bool
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:   GetEnv("SERVER_PORT", ":3000"),
		DataFilePath: GetEnv("DATA_FILE", ".data/project.yaml"),
		EnableDMX:    GetEnvBool("ENABLE_DMX", true),
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}
