package main

import (
	"log"

	"github.com/b1tvect0r/exchangerates/pkg/server"
)

func main() {
	s, err := server.New(server.WithAESAPIKeyProvider("key"))
	if err != nil {
		log.Fatalf("failed to create server: %s", err.Error())
	}

	if err = s.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}
