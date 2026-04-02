package main

import (
	"fmt"

	"github.com/DevitoDbug/portfolio/internals/server"
)

func main() {
	port := ":8081"
	server := server.NewServer(port)

	err := server.StartServer()
	if err != nil {
		fmt.Printf("failed to start server. Error: %v\n", err)
	}
}
