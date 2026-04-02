package config

import (
	"fmt"
	"os"
	"strings"
)

type EnvironmentConfig struct {
	RunningEnvironment string
}

func GetEnvironmentConfig() (*EnvironmentConfig, error) {
	runningEnv := os.Getenv("ENVIRONMENT")
	if strings.TrimSpace(runningEnv) == "" {
		return nil, fmt.Errorf(".env value ENVIRONMENT not found")
	}

	return &EnvironmentConfig{
		RunningEnvironment: runningEnv,
	}, nil
}
