package wok

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// takes in context and parses the session_id cookie, if the cookie is not present
// function returns an empty string and an error, if cookie is in the request but does
// but the user does not have the role of "user" or "admin", return an error stating they are not authorized
// if the cookie is in the request and the users role is user or admin, returns the session_id as a string and nil
// for the error
func UserisValid(ctx Context) error {
	id, err := ctx.Req.Cookie("session_id")
	if err != nil {
		return fmt.Errorf("%s", http.StatusText(http.StatusUnauthorized))
	}

	uuid, err := uuid.Parse(id.Value)
	if err != nil {
		return err
	}

	qs := "SELECT session_id FROM users WHERE session_id=$1"
	row := Database.QueryRow(qs, uuid.String())
	var session string
	if err := row.Scan(&session); err != nil {
		return err
	}

	return nil

}

func UserisAdmin(ctx Context) error {
	id, err := ctx.Req.Cookie("session_id")
	if err != nil {
		return err
	}

	uuid, err := uuid.Parse(id.Value)
	if err != nil {
		return err
	}

	qs := "SELECT role FROM users WHERE session_id=$1"
	row := Database.QueryRow(qs, uuid.String())
	var role string

	if err := row.Scan(&role); err != nil {
		return err
	}

	if role != "admin" {
		return fmt.Errorf(http.StatusText(http.StatusUnauthorized))
	}

	return nil

}
