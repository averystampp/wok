package wok

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DatabaseOpts struct {
	Host            string
	Port            int
	SSL             bool
	User            string
	Password        string
	Dbname          string
	MigrationFolder string
}

var DB *sql.DB

func Connect(config *DatabaseOpts) {
	sslString := "disable"
	if config.SSL {
		sslString = "enable"
	}
	psqlInfo := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Dbname,
		sslString,
	)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}
	if config.MigrationFolder != "" {
		folder, err := os.ReadDir(config.MigrationFolder)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range folder {
			data, err := os.ReadFile(config.MigrationFolder + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec(string(data))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
