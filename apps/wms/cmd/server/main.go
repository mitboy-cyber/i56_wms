// Package main is the entry point for the I56 WMS backend server.
// It initializes the server from internal/server and starts listening.
package main

import (
	"log"

	"github.com/i56/i56-apps/i56-wms/internal/server"
)

func main() {
	srv, err := server.New(server.Config{
		Port:      8080,
		DBDSN:     "postgres://ubuntu@localhost:5432/i56_dev?sslmode=disable",
		StaticDir: "static",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
