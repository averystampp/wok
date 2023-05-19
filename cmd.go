package wok

import (
	"bufio"
	"fmt"
	"os"
)

func NewAdmin(conf *DbConfig) {
	db, err := DbConnect(conf)
	if err != nil {
		panic(err)
	}
	database = db
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
}
