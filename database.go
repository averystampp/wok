package wok

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	Host            string
	Port            int
	SSL             bool
	User            string
	Password        string
	Dbname          string
	MigrationFolder string
}

type Database struct {
	Database *sql.DB
}

var Store Database

func (d *Database) Connect(config *Config) {
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS csrf (
		id serial PRIMARY KEY,
		token VARCHAR( 1004 ) NOT NULL,
		expires VARCHAR( 256 ) NOT NULL
	);`)

	if err != nil {
		log.Fatal(err)
	}

	Store.Database = db
}
