package main

import (
	"log"

	"github.com/averystampp/wok"
)

func main() {
	app := wok.NewWok(false, ":8080", "", "")

	dbconf := wok.DbConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "docker",
		Dbname:          "postgres",
		MigrationFolder: "./migrations",
	}

	log.Fatal(app.StartServer(dbconf))

}
