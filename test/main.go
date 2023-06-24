package main

import (
<<<<<<< HEAD
	"net/http"
=======
	"log"
>>>>>>> a3201194f11bacd06b9c5392c66466b61727dcad

	"github.com/averystampp/wok"
)

func main() {
<<<<<<< HEAD

	config := &wok.WokConfig{
		Addr:    ":8080",
		Handler: new(http.ServeMux),
		TLS:     false,
	}
=======
	app := wok.NewWok(false, ":8080", "", "")
>>>>>>> a3201194f11bacd06b9c5392c66466b61727dcad

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
