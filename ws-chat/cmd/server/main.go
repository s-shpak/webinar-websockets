package main

import (
	"fmt"
	"log"

	config "ws-chat/internal/config/server"
	"ws-chat/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfig()
	srv, err := server.NewServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to create a new server: %w", err)
	}
	if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		return err
	}
	return nil
}
