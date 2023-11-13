package wok

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Wok struct {
	Address  string
	mux      *http.ServeMux
	CertFile string
	KeyFile  string
}

const (
	WOK_VERSION = "Wok v1.0.0"
)

var (
	err        error
	wokSession *Session
	wokLogger  *WokLog
)

// Return a new Wok server struct
func NewWok(addr string) *Wok {
	return &Wok{
		Address: addr,
		mux:     http.NewServeMux(),
	}
}

func (w *Wok) WithDatabase(config *Config) {
	Store.Connect(config)
}

func (w *Wok) WithSession() {
	wokSession = newSession()
}

func (w *Wok) WithLogger() {
	wokLogger, err = newLogger()
	if err != nil {
		log.Fatal(err)
	}
}

func (w *Wok) StartWok() {
	for _, arg := range os.Args {
		if arg == "drop table" {
			DropTable()
			os.Exit(0)
		}
	}
	w.startServer()
}

func (w *Wok) startServer() {
	if wokLogger != nil {
		wokLogger.General("Server starting on port" + w.Address)
	} else {
		fmt.Println("Server starting on port" + w.Address)
	}
	if w.CertFile != "" && w.KeyFile != "" {
		log.Fatal(http.ListenAndServeTLS(w.Address, w.CertFile, w.KeyFile, w.mux))
	} else {
		log.Fatal(http.ListenAndServe(w.Address, w.mux))
	}
}
