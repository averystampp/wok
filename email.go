package wok

import (
	"encoding/json"
	"net/smtp"
	"os"
)

type Email struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
}

func SendCreateUserEmail() error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	body, err := os.ReadFile("../public/NewUser.html")

	if err != nil {
		return err
	}

	subject := "Subject: Does this work?\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	if err := smtp.SendMail("smtp.gmail.com:587", auth, "noreply@wok.app", []string{"amstampp18@gmail.com"}, []byte(subject+mime+string(body))); err != nil {
		return err
	}

	return nil
}

func EmailsinQueue() ([]byte, error) {
	rows, err := Database.Query("SELECT * FROM signups")
	if err != nil {
		return nil, err
	}
	var elist []Email
	var e Email
	for rows.Next() {
		rows.Scan(&e.Id, &e.Address)
		elist = append(elist, e)
	}

	resp, err := json.Marshal(elist)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func AddEmailtoQueue(e *Email) error {
	_, err := Database.Exec("INSERT INTO signups (email) VALUES ($1)", e.Address)
	if err != nil {
		return err
	}

	return nil
}

func RemoveEmailFromQueue(e *Email) error {
	return nil
}
