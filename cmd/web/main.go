package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/DevitoDbug/portfolio/internals/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	slog.Info("staring server...", "port", port)
	server := server.NewServer(port)

	err := server.StartServer()
	if err != nil {
		fmt.Printf("failed to start server. Error: %v\n", err)
	}
}
