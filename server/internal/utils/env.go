package utils

import (
	"os"

	"github.com/mpcarolin/cinematch-server/internal/constants/env"
)

/*
 * Get the environment variable from the system
 * @return {string} The environment variable
 */
func GetEnv() string {
	envVar := os.Getenv("ENV")
	if envVar == "" {
		return env.Development // default to development
	}
	return envVar
}