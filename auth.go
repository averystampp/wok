package wok

import (
	"fmt"
	"net/http"
)

// takes in context and parses the session_id cookie, if the cookie is not present
// function returns an empty string and an error, if cookie is in the request but does
// but the user does not have the role of "user" or "admin", return an error stating they are not authorized
// if the cookie is in the request and the users role is user or admin, returns the session_id as a string and nil
// for the error
func UserisUser(ctx Context) (string, error) {
	id, err := ctx.Req.Cookie("session_id")
	if err != nil {
		return "", fmt.Errorf("%s", http.StatusText(http.StatusUnauthorized))
	}

	qs := "SELECT role FROM users WHERE session_id=$1"
	row := Database.QueryRow(qs, id.Value)
	var role string
	row.Scan(&role)

	if role == "user" {
		return id.Value, nil
	}

	if role == "admin" {
		return id.Value, nil
	}

	return "", fmt.Errorf("%s", http.StatusText(http.StatusUnauthorized))

}

// takes in context, returns string of users id to be used by handler if necessary
// also returns an error, wrote http into logic because there is multiple errors from this function
// this function mirrors the UserisUser function but only for admins
func UserisAdmin(ctx Context) (string, error) {
	id := ctx.Req.Header.Get("session_id")
	if id == "" {
		return "", fmt.Errorf("must include a session id")
	}

	qs := "SELECT role FROM users WHERE session_id=$1"
	row := Database.QueryRow(qs, id)
	var role string
	row.Scan(&role)

	if role == "admin" {
		return id, nil
	}

	return "", fmt.Errorf("user is not authorized")

}

func QueryUserfromDb(ctx Context) (string, error) {
	id, err := ctx.Req.Cookie("session_id")
	if err != nil {
		return "", fmt.Errorf("must send a session_id")
	}

	qs := "SELECT role FROM users WHERE session_id=$1"
	row := Database.QueryRow(qs, id.Value)
	var role string
	row.Scan(&role)

	if role != "" {
		return role, nil
	}

	return "", fmt.Errorf("error getting role from database")

}
