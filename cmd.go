package wok

import (
	"bufio"
	"fmt"
	"os"
)

// drops the users table, nice function to have when testing and you want to start fresh
func dropTable(conf *DbConfig) {
	db, err := directToDB(conf)
	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Type Table Name to drop: ")
		scanner.Scan()
		table := scanner.Text()
		if len(table) == 0 {
			break
		}
		qs := "DROP TABLE " + table
		_, err = db.Exec(qs)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

}
