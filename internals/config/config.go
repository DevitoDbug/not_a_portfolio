// Package config - Holds all configurations for the application e.g environment variables
package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func init() {
	locations := []string{
		".env",                      // Current directory
		filepath.Join("..", ".env"), // Parent directory (cross-platform)
	}

	fmt.Printf("Tried fetching from the following location %+v\n", locations)

	loaded := false
	for _, loc := range locations {
		if err := godotenv.Load(loc); err == nil {
			loaded = true
			break
		}
	}

	if !loaded {
		log.Printf("Warning: Could not load .env file from any location")
	}
}
