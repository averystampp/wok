package main

import (
	"net/http"

	"github.com/averystampp/wok"
)

func main() {

	config := &wok.WokConfig{
		Addr:    ":8080",
		Handler: new(http.ServeMux),
		TLS:     false,
	}

	dbconf := wok.DbConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "docker",
		Dbname:          "postgres",
		MigrationFolder: "./migrations",
	}

	wok.StartServer(config, dbconf)
}
