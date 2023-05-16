package main

import (
	"net/http"

	"github.com/averystampp/wok"
)

func main() {
	config := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	dbconf := wok.DbConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "docker",
		Dbname:   "postgres",
	}

	wok.Router()
	wok.StartServer(config, dbconf)
}
