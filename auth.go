package wok

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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
		return err
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

func CsrfProtect(ctx Context) error {
	token := ctx.Req.FormValue("_csrf")

	if token == "" {
		return errors.New("no csrf token in request")
	}

	uuid, err := uuid.Parse(token)

	if err != nil {
		return errors.New("token is not in valid format")
	}

	qs := "SELECT * FROM csrf WHERE token=$1"

	row := Database.QueryRow(qs, uuid.String())

	tokenRow := struct {
		id      int
		token   string
		expires string
	}{}

	if err := row.Scan(&tokenRow.id, &tokenRow.token, &tokenRow.expires); err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339, tokenRow.expires)
	if err != nil {
		return err
	}

	if !time.Now().Local().Before(t) {
		return errors.New("token is expired")
	}

	return nil
}

// Insert this into your handlers (GET methods) to create a CSRF token for the client if they don't have one
// or if the one they have is expired, works in tandem with CsrfProtect for POST handlers
func CsrfCreate(ctx Context) error {
	csrf, err := ctx.Req.Cookie("_csrf")
	if err != nil {
		if err := createCsrfToken(ctx); err != nil {
			return err
		}

		return err
	}

	qs := "SELECT * FROM csrf WHERE token=$1"

	row := Database.QueryRow(qs, csrf.Value)

	tokenRow := struct {
		id      int
		token   string
		expires string
	}{}

	if err := row.Scan(&tokenRow.id, &tokenRow.token, &tokenRow.expires); err != nil {
		if err := createCsrfToken(ctx); err != nil {
			return err
		}
		return err
	}
	// creates token if the row return nothing for the token value
	if tokenRow.token == "" {
		if err := createCsrfToken(ctx); err != nil {
			return err
		}
	}

	t, err := time.Parse(time.RFC3339, tokenRow.expires)
	if err != nil {
		return err
	}
	// creates token if the current token given is expired
	if !time.Now().Local().Before(t) {
		if err := createCsrfToken(ctx); err != nil {
			return err
		}
	}

	return nil
}

func createCsrfToken(ctx Context) error {
	qs := "INSERT INTO csrf (token, expires) VALUES ($1, $2)"

	token := uuid.New()
	expires := time.Now().Add(time.Minute * 15).Local()

	_, err := Database.Exec(qs, token.String(), expires.Format(time.RFC3339))
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		HttpOnly: true,
		Expires:  expires,
		Value:    token.String(),
		Name:     "_csrf",
	}

	http.SetCookie(ctx.Resp, cookie)
	http.Redirect(ctx.Resp, ctx.Req, ctx.Req.URL.Path, http.StatusSeeOther)
	return nil
}
