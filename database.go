package wok

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type DbConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Dbname          string
	MigrationFolder string
}

var Database *sql.DB

// connects to database on server startup, will create the users table if it not already in the database
func DbStartup(c *DbConfig) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Connected to database")
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
	  		id         serial PRIMARY KEY,
	  		email      VARCHAR( 128 ) NOT NULL,
	  		password     VARCHAR( 255 ) NOT NULL,
	  		role      VARCHAR( 128 ) NOT NULL,
	  		session_id      VARCHAR( 128 ) NOT NULL,
			logged_in BOOLEAN NOT NULL
			);`)

	if err != nil {
		fmt.Println("users")
		fmt.Println(err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS signups (
		id         serial PRIMARY KEY,
		email      VARCHAR( 128 ) NOT NULL
	);`)

	if err != nil {
		fmt.Println("signups")
		fmt.Println(err)
	}

	dir, err := os.ReadDir(c.MigrationFolder)

	if err != nil {
		fmt.Println(err)
	}

	for _, migration := range dir {

		if strings.Split(migration.Name(), ".")[1] != "sql" {
			return nil, fmt.Errorf("file does not have .sql extension: %s", migration.Name())
		}

		migrate, err := os.ReadFile("./migrations/" + migration.Name())
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec(string(migrate))
		if err != nil {
			fmt.Println(err)
		}

	}

	Database = db
	return db, nil

}

// connect to database without checking if users table is in the db
func DbConnect(c *DbConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil

}
