package wok

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

type DbConfig struct {
	Host            string
	Port            int
	SSL             bool
	User            string
	Password        string
	Dbname          string
	MigrationFolder string
}

var Database *sql.DB

// connects to database on server startup, will create the users table if it not already in the database
func connectToDB(c *DbConfig) error {

	sslString := "disable"
	if c.SSL {
		sslString = "enable"
	}

	psqlInfo := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Dbname,
		sslString,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS csrf (
			id serial PRIMARY KEY,
			token VARCHAR( 1004 ) NOT NULL,
			expires VARCHAR( 256 ) NOT NULL
	);`)

	if err != nil {
		return err
	}

	dir, err := os.ReadDir(c.MigrationFolder)

	if err != nil {
		return err
	}

	for _, migration := range dir {
		if len(dir) < 1 {
			break
		}
		fileExt := strings.Split(migration.Name(), ".")[1]
		if fileExt != "sql" && fileExt != "psql" {
			return fmt.Errorf("file does not have .sql extension: %s", migration.Name())
		}

		migrate, err := os.ReadFile(c.MigrationFolder + "/" + migration.Name())
		if err != nil {
			return err
		}
		_, err = db.Exec(string(migrate))
		if err != nil {
			return err
		}
	}

	Database = db
	return nil
}

// connect to database without checking if users table is in the db
func directToDB(c *DbConfig) (*sql.DB, error) {
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
