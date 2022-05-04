package main

import (
	"os"

	"github.com/cin-lawrence/gosandbox/pkg/api"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	server := api.NewAPIServer()
	server.Run("0.0.0.0", 8080)
}
