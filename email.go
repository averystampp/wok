package wok

import (
	"net/smtp"
	"os"
)

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
