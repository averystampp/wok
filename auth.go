package wok

import (
	"fmt"
	"net/http"
)

func UserisUser(ctx Context) error {
	id, err := ctx.r.Cookie("session_id")
	if err != nil {
		return fmt.Errorf("%s", http.StatusText(http.StatusInternalServerError))
	}

	qs := "SELECT role FROM users WHERE session_id=$1"
	row := database.QueryRow(qs, id.Value)
	var role string
	row.Scan(&role)

	if role == "user" {
		return nil
	}

	if role == "admin" {
		return nil
	}

	return fmt.Errorf("user is not authorized")

}

func UserisAdmin(ctx Context) error {
	id, err := ctx.r.Cookie("session_id")
	if err != nil {
		return fmt.Errorf("%s", http.StatusText(http.StatusInternalServerError))
	}

	qs := "SELECT role FROM users WHERE session_id=$1"
	row := database.QueryRow(qs, id.Value)
	var role string
	row.Scan(&role)

	if role == "admin" {
		return nil
	}

	return fmt.Errorf("user is not authorized")

}
