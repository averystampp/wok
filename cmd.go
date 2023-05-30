package wok

import (
	"bufio"
	"fmt"
	"os"
)

// function to create an admin, currently creating users via the /user route creates a user
// with the role of "user" this is a way for developers to create admins without having to
// expose an endpoint
func NewAdmin(conf *DbConfig) {
	db, err := DbConnect(conf)
	if err != nil {
		panic(err)
	}
	Database = db
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter Username: ")
		scanner.Scan()
		username := scanner.Text()
		if len(username) == 0 {
			break
		}
		fmt.Print("Enter Lastname: ")
		scanner.Scan()
		lastname := scanner.Text()
		if len(lastname) == 0 {
			break
		}

		fmt.Print("Enter Password: ")
		scanner.Scan()
		password := scanner.Text()
		if len(password) < 5 {
			fmt.Println("Password must be more than 5 characters long")
		} else {

			cliUser := &User{
				FirstName: username,
				LastName:  lastname,
				Password:  password,
			}

			if err := CreateAdmin(cliUser); err != nil {
				panic(err)
			}

			break
		}
	}
}

// drops the users table, nice function to have when testing and you want to start fresh
func DropUsersTable(conf *DbConfig) {
	db, err := DbConnect(conf)
	if err != nil {
		fmt.Println(err)
	}
	qs := "DROP TABLE users"
	_, err = db.Exec(qs)
	if err != nil {
		fmt.Println(err)
	}

	qs2 := "DROP TABLE signups"
	_, err = db.Exec(qs2)
	if err != nil {
		fmt.Println(err)
	}
}
