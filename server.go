package wok

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
)

// takes in a server struct and runs the server on the speified struct with the specified handler
// also takes a db config, validates it and then run dbstartup
// TODO: remove error return and use panics instead beause each of these processes only run once on startup
// maybe implement a sudo recovery to ensure server is more resiliant
func (config *Wok) StartServer(dbconfig DbConfig) error {

	if err := validatedbconfig(dbconfig); err != nil {
		panic(err)
	}

	for _, arg := range os.Args {
		if arg == "createuser" {
			NewAdmin(&dbconfig)
			os.Exit(0)
		}

		if arg == "droptable" {
			DropUsersTable(&dbconfig)
			os.Exit(0)
		}
	}

	// Startup db, this will create a users table if it doesnt already exist.
	// Also prints to the console on successful connection
	_, err := DbStartup(&dbconfig)
	if err != nil {
		return err
	}
	fmt.Println("WOK version 0.0.0")
	fmt.Println("-------------------------------------")
	fmt.Printf("| Server starting on localhost%s |\n", config.address)
	fmt.Println("-------------------------------------")

	// call the default router
	DefaultRouter(config)
	if config.tls {
		if err := http.ListenAndServeTLS(config.address, config.certFile, config.keyFile, config.mux); err != nil {
			return err
		}
	} else {
		if err := http.ListenAndServe(config.address, config.mux); err != nil {
			return err
		}
	}

	defer Database.Close()

	return nil

}

// validate the dbonfig strut passed into startup server, this ensures that all fields have a value
// but is agnostic to the specific data passed into each field
func validatedbconfig(dbconfig DbConfig) error {
	v := reflect.ValueOf(dbconfig)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			return fmt.Errorf("must fill all config fields out")
		}

		if v.Field(i).Kind() == reflect.Int && v.Field(i).Int() == 0 {
			return fmt.Errorf("must include a database port")
		}
	}

	return nil
}
