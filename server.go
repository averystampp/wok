package wok

import (
	"fmt"
	"net/http"
	"os"
)

const (
	WOK_VERSION = "Wok v1.0.0"
)

func (w *Wok) StartWok(db DbConfig) {
	if err := validatedbconfig(db); err != nil {
		panic(err)
	}

	for _, arg := range os.Args {
		if arg == "createuser" {
			newAdmin(&db)
			os.Exit(0)
		}

		if arg == "droptable" {
			dropUsersTable(&db)
			os.Exit(0)
		}
	}

	if err := connectToDB(&db); err != nil {
		panic(err)
	}

	fmt.Println(WOK_VERSION)
	fmt.Println("-------------------------------------")
	fmt.Printf("| Server starting on port %s |\n", w.address)
	fmt.Println("-------------------------------------")

	// call the default router
	DefaultRouter(w)

	http.ListenAndServe(w.address, w.mux)
	defer Database.Close()
}

func (w *Wok) StartWokTLS(db DbConfig) {
	if err := validatedbconfig(db); err != nil {
		panic(err)
	}

	for _, arg := range os.Args {
		if arg == "createuser" {
			newAdmin(&db)
			os.Exit(0)
		}

		if arg == "droptable" {
			dropUsersTable(&db)
			os.Exit(0)
		}
	}

	if err := connectToDB(&db); err != nil {
		panic(err)
	}
	fmt.Println(WOK_VERSION)
	fmt.Println("-------------------------------------")
	fmt.Printf("| Server starting on port %s |\n", w.address)
	fmt.Println("-------------------------------------")

	// call the default router
	DefaultRouter(w)

	http.ListenAndServeTLS(w.address, w.certFile, w.keyFile, w.mux)
	defer Database.Close()
}
