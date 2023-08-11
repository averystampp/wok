package wok

import (
	"bufio"
	"bytes"
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

func convertLogToJSON() {
	b, err := os.ReadFile("log.log")
	if err != nil {
		fmt.Println(err)
	}
	b = bytes.ReplaceAll(b, []byte("}"), []byte("},"))
	b = bytes.Replace(b, []byte("{"), []byte("[{"), 1)
	b[len(b)-2] = ']'

	out, err := os.Create("output.json")
	if err != nil {
		fmt.Println(err)
	}

	out.Write(b)

}
