package main

import (
	"github.com/averystampp/wok"
)

func main() {
	config := &wok.WokConfig{
		Addr:    ":8080",
		Handler: nil,
		TLS:     false,
	}

	dbconf := wok.DbConfig{
		Host:            "db",
		Port:            5432,
		User:            "postgres",
		Password:        "docker",
		Dbname:          "postgres",
		MigrationFolder: "./test/migrations",
	}

	wok.StartServer(config, dbconf)
}
