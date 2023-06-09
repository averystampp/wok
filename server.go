package wok

import (
	"fmt"
	"net/http"
	"os"
)

type Wok struct {
	address  string
	prefix   string
	mux      *http.ServeMux
	certFile string
	keyFile  string
}

const (
	WOK_VERSION = "Wok v1.0.0"
)

// Return a new Wok server struct
func NewWok(addr string) *Wok {
	return &Wok{
		address: addr,
		mux:     &http.ServeMux{},
	}
}

func NewWokTLS(addr, certfile, keyfile string) *Wok {
	return &Wok{
		address:  addr,
		certFile: certfile,
		keyFile:  keyfile,
		mux:      &http.ServeMux{},
	}
}

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
			dropTable(&db)
			os.Exit(0)
		}
	}

	if err := connectToDB(&db); err != nil {
		panic(err)
	}
	defer Database.Close()

	fmt.Println(WOK_VERSION)
	fmt.Printf("---------------------------------\n| Server starting on port %s |\n---------------------------------\n", w.address)

	DefaultRouter(w)

	if w.certFile != "" && w.keyFile != "" {
		http.ListenAndServeTLS(w.address, w.certFile, w.keyFile, w.mux)
	} else {
		http.ListenAndServe(w.address, w.mux)
	}

}
