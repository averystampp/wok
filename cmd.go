package wok

import (
	"bufio"
	"fmt"
	"os"
)

func NewUser(conf *DbConfig) {
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
				Role:      "admin",
			}

			if err := CreateUser(cliUser); err != nil {
				panic(err)
			}

			break
		}
	}
}
