package wok

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/smtp"
	"os"
)

type Email struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

func SendCreateUserEmail(email, tmpl string) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	body, err := template.ParseFiles(tmpl)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := body.Execute(buf, nil); err != nil {
		return err
	}

	subject := "Subject: Does this work?\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	if err := smtp.SendMail("smtp.gmail.com:587", auth, "noreply@wok.app", []string{email}, []byte(subject+mime+buf.String())); err != nil {
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
		rows.Scan(&e.Id, &e.Address, &e.Name)
		elist = append(elist, e)
	}

	resp, err := json.Marshal(elist)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func (e *Email) AddEmailtoQueue() error {
	_, err := Database.Exec("INSERT INTO signups (email, name) VALUES ($1, $2)", e.Address, e.Name)
	if err != nil {
		return err
	}

	return nil
}

func (e *Email) RemoveEmailFromQueue() error {
	_, err := Database.Exec("DELETE FROM signups WHERE id=$1", e.Id)
	if err != nil {
		return err
	}
	return nil
}
