package wok

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type Wok struct {
	Host     string
	Address  string
	prefix   string
	mux      *http.ServeMux
	CertFile string
	KeyFile  string
	Database bool
}

const (
	WOK_VERSION = "Wok v1.0.0"
)

var (
	pool   *sync.Pool
	logger *woklog
	err    error
)

// Return a new Wok server struct
func NewWok(addr string, db bool) *Wok {
	return &Wok{
		Address:  addr,
		mux:      &http.ServeMux{},
		Database: db,
	}
}

func initpool() {
	pool = &sync.Pool{
		New: func() any {
			return Context{
				Resp: nil,
				Req:  nil,
				Ctx:  nil,
			}
		},
	}
}

func (w *Wok) StartWok(db ...DbConfig) {
	initpool()
	logger, err = newLogger()
	if err != nil {
		panic(err)
	}
	if w.Database {
		if err := validatedbconfig(db[0]); err != nil {
			panic(err)
		}

		for _, arg := range os.Args {
			if arg == "droptable" {
				dropTable(&db[0])
				os.Exit(0)
			}

			if arg == "logJSON" {
				convertLogToJSON()
				os.Exit(0)
			}
		}

		if err := connectToDB(&db[0]); err != nil {
			panic(err)
		}

		defer Database.Close()
		w.startServer()
	}

	for _, arg := range os.Args {
		if arg == "logJSON" {
			convertLogToJSON()
			os.Exit(0)
		}
	}
	w.startServer()
}

var WokSession = StartSession()

func (w *Wok) startServer() {
	fmt.Println(WOK_VERSION)
	fmt.Printf("---------------------------------\n| Server starting on port %s |\n---------------------------------\n", w.Address)
	if w.CertFile != "" && w.KeyFile != "" {
		log.Fatal(http.ListenAndServeTLS(w.Address, w.CertFile, w.KeyFile, w.mux))
	} else {
		log.Fatal(http.ListenAndServe(w.Address, w.mux))
	}
}
